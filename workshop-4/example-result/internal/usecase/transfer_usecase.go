package usecase

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"temp-backend/internal/domain"
	"time"
)

// transferUsecase implements domain.TransferUsecase
type transferUsecase struct {
	transferRepo    domain.TransferRepository
	userRepo        domain.UserRepository
	pointLedgerRepo domain.PointLedgerRepository
	db              *sql.DB
}

// NewTransferUsecase creates a new transfer usecase
func NewTransferUsecase(
	transferRepo domain.TransferRepository,
	userRepo domain.UserRepository,
	pointLedgerRepo domain.PointLedgerRepository,
	db *sql.DB,
) domain.TransferUsecase {
	return &transferUsecase{
		transferRepo:    transferRepo,
		userRepo:        userRepo,
		pointLedgerRepo: pointLedgerRepo,
		db:              db,
	}
}

// CreateTransfer creates a new transfer with atomic transaction
func (u *transferUsecase) CreateTransfer(req *domain.TransferCreateRequest) (*domain.Transfer, error) {
	// Validation
	if err := u.validateTransferRequest(req); err != nil {
		return nil, err
	}

	// Check if users exist
	fromExists, err := u.userRepo.Exists(req.FromUserID)
	if err != nil {
		return nil, err
	}
	if !fromExists {
		return nil, errors.New("sender user not found")
	}

	toExists, err := u.userRepo.Exists(req.ToUserID)
	if err != nil {
		return nil, err
	}
	if !toExists {
		return nil, errors.New("receiver user not found")
	}

	// Check if sender can transfer to this user (rule 3: cannot transfer to last recipient)
	lastTransfer, err := u.transferRepo.GetLastTransferFromUser(req.FromUserID)
	if err != nil {
		return nil, err
	}
	if lastTransfer != nil && lastTransfer.ToUserID == req.ToUserID && lastTransfer.Status == "completed" {
		return nil, errors.New("cannot transfer to the same user as the last completed transfer")
	}

	// Check if sender has enough points
	senderPoints, err := u.userRepo.GetUserPoints(req.FromUserID)
	if err != nil {
		return nil, err
	}
	if senderPoints < req.Amount {
		return nil, errors.New("insufficient points")
	}

	// Generate idempotency key
	idemKey, err := u.generateIdemKey()
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create transfer record
	transfer := &domain.Transfer{
		IdemKey:    idemKey,
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Amount:     req.Amount,
		Status:     "pending",
		Note:       req.Note,
	}

	createdTransfer, err := u.transferRepo.Create(transfer)
	if err != nil {
		return nil, err
	}

	// Process the transfer immediately (for simplicity)
	err = u.processTransfer(createdTransfer, tx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the updated transfer
	return u.transferRepo.GetByIdemKey(idemKey)
}

// GetTransferByIdemKey retrieves a transfer by idempotency key
func (u *transferUsecase) GetTransferByIdemKey(idemKey string) (*domain.Transfer, error) {
	if idemKey == "" {
		return nil, errors.New("idempotency key is required")
	}

	transfer, err := u.transferRepo.GetByIdemKey(idemKey)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, errors.New("transfer not found")
	}

	return transfer, nil
}

// GetTransfersByUserID retrieves transfers by user ID with pagination
func (u *transferUsecase) GetTransfersByUserID(userID int, page, pageSize int) (*domain.TransferListResponse, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Validate pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}

	// Check if user exists
	exists, err := u.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	transfers, total, err := u.transferRepo.GetByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &domain.TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// processTransfer processes the transfer within a transaction
func (u *transferUsecase) processTransfer(transfer *domain.Transfer, tx *sql.Tx) error {
	// Get current points
	senderPoints, err := u.userRepo.GetUserPoints(transfer.FromUserID)
	if err != nil {
		return err
	}

	receiverPoints, err := u.userRepo.GetUserPoints(transfer.ToUserID)
	if err != nil {
		return err
	}

	// Check points again (double-check)
	if senderPoints < transfer.Amount {
		now := time.Now()
		failReason := "insufficient points at processing time"
		u.transferRepo.UpdateStatus(transfer.IdemKey, "failed", &now, &failReason)
		return errors.New("insufficient points")
	}

	// Update points
	newSenderPoints := senderPoints - transfer.Amount
	newReceiverPoints := receiverPoints + transfer.Amount

	err = u.userRepo.UpdatePoints(transfer.FromUserID, newSenderPoints)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdatePoints(transfer.ToUserID, newReceiverPoints)
	if err != nil {
		return err
	}

	// Create ledger entries
	senderLedger := &domain.PointLedger{
		UserID:       transfer.FromUserID,
		Change:       -transfer.Amount,
		BalanceAfter: newSenderPoints,
		EventType:    "transfer_out",
		TransferID:   &transfer.ID,
	}

	receiverLedger := &domain.PointLedger{
		UserID:       transfer.ToUserID,
		Change:       transfer.Amount,
		BalanceAfter: newReceiverPoints,
		EventType:    "transfer_in",
		TransferID:   &transfer.ID,
	}

	err = u.pointLedgerRepo.Create(senderLedger)
	if err != nil {
		return err
	}

	err = u.pointLedgerRepo.Create(receiverLedger)
	if err != nil {
		return err
	}

	// Update transfer status to completed
	now := time.Now()
	return u.transferRepo.UpdateStatus(transfer.IdemKey, "completed", &now, nil)
}

// validateTransferRequest validates the transfer request
func (u *transferUsecase) validateTransferRequest(req *domain.TransferCreateRequest) error {
	if req.FromUserID <= 0 {
		return errors.New("invalid sender user ID")
	}
	if req.ToUserID <= 0 {
		return errors.New("invalid receiver user ID")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if req.FromUserID == req.ToUserID {
		return errors.New("cannot transfer to yourself")
	}
	if req.Note != nil && len(*req.Note) > 512 {
		return errors.New("note too long (max 512 characters)")
	}

	// Rule 2: Transfer amount should not exceed 2 points and should have max 2 decimal places
	if req.Amount > 2.0 {
		return errors.New("transfer amount cannot exceed 2 points")
	}

	// Check decimal places (max 2)
	rounded := math.Round(req.Amount*100) / 100
	if req.Amount != rounded {
		return errors.New("transfer amount cannot have more than 2 decimal places")
	}

	return nil
}

// generateIdemKey generates a unique idempotency key
func (u *transferUsecase) generateIdemKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}
