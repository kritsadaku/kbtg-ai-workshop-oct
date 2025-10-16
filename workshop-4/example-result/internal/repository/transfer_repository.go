package repository

import (
	"database/sql"
	"temp-backend/internal/domain"
	"time"
)

// transferRepository implements domain.TransferRepository
type transferRepository struct {
	db *sql.DB
}

// NewTransferRepository creates a new transfer repository
func NewTransferRepository(db *sql.DB) domain.TransferRepository {
	return &transferRepository{
		db: db,
	}
}

// Create creates a new transfer
func (r *transferRepository) Create(transfer *domain.Transfer) (*domain.Transfer, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	result, err := r.db.Exec(`
		INSERT INTO transfers (from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		transfer.FromUserID, transfer.ToUserID, transfer.Amount, transfer.Status,
		transfer.Note, transfer.IdemKey, now, now)

	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByIdemKey(transfer.IdemKey)
}

// GetByIdemKey retrieves a transfer by idempotency key
func (r *transferRepository) GetByIdemKey(idemKey string) (*domain.Transfer, error) {
	transfer := &domain.Transfer{}
	var createdAt, updatedAt string
	var completedAt, failReason sql.NullString

	err := r.db.QueryRow(`
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key,
			   created_at, updated_at, completed_at, fail_reason
		FROM transfers WHERE idempotency_key = ?`, idemKey).Scan(
		&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount,
		&transfer.Status, &transfer.Note, &transfer.IdemKey,
		&createdAt, &updatedAt, &completedAt, &failReason,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	if transfer.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		return nil, err
	}
	if transfer.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		return nil, err
	}

	if completedAt.Valid {
		if parsed, err := time.Parse("2006-01-02 15:04:05", completedAt.String); err == nil {
			transfer.CompletedAt = &parsed
		}
	}

	if failReason.Valid {
		transfer.FailReason = &failReason.String
	}

	return transfer, nil
}

// GetByUserID retrieves transfers by user ID with pagination
func (r *transferRepository) GetByUserID(userID int, page, pageSize int) ([]*domain.Transfer, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?`, userID, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := r.db.Query(`
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key,
			   created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, userID, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []*domain.Transfer
	for rows.Next() {
		transfer := &domain.Transfer{}
		var createdAt, updatedAt string
		var completedAt, failReason sql.NullString

		err := rows.Scan(
			&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount,
			&transfer.Status, &transfer.Note, &transfer.IdemKey,
			&createdAt, &updatedAt, &completedAt, &failReason,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse timestamps
		if transfer.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			return nil, 0, err
		}
		if transfer.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			return nil, 0, err
		}

		if completedAt.Valid {
			if parsed, err := time.Parse("2006-01-02 15:04:05", completedAt.String); err == nil {
				transfer.CompletedAt = &parsed
			}
		}

		if failReason.Valid {
			transfer.FailReason = &failReason.String
		}

		transfers = append(transfers, transfer)
	}

	return transfers, total, nil
}

// GetLastTransferFromUser retrieves the last completed transfer from a specific user
func (r *transferRepository) GetLastTransferFromUser(fromUserID int) (*domain.Transfer, error) {
	transfer := &domain.Transfer{}
	var createdAt, updatedAt string
	var completedAt, failReason sql.NullString

	err := r.db.QueryRow(`
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key,
			   created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? AND status = 'completed'
		ORDER BY completed_at DESC
		LIMIT 1`, fromUserID).Scan(
		&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount,
		&transfer.Status, &transfer.Note, &transfer.IdemKey,
		&createdAt, &updatedAt, &completedAt, &failReason,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	if transfer.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		return nil, err
	}
	if transfer.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		return nil, err
	}

	if completedAt.Valid {
		if parsed, err := time.Parse("2006-01-02 15:04:05", completedAt.String); err == nil {
			transfer.CompletedAt = &parsed
		}
	}

	if failReason.Valid {
		transfer.FailReason = &failReason.String
	}

	return transfer, nil
}

// UpdateStatus updates transfer status
func (r *transferRepository) UpdateStatus(idemKey string, status string, completedAt *time.Time, failReason *string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	var completedAtStr sql.NullString
	if completedAt != nil {
		completedAtStr.String = completedAt.Format("2006-01-02 15:04:05")
		completedAtStr.Valid = true
	}

	var failReasonStr sql.NullString
	if failReason != nil {
		failReasonStr.String = *failReason
		failReasonStr.Valid = true
	}

	_, err := r.db.Exec(`
		UPDATE transfers 
		SET status = ?, updated_at = ?, completed_at = ?, fail_reason = ?
		WHERE idempotency_key = ?`,
		status, now, completedAtStr, failReasonStr, idemKey)

	return err
}
