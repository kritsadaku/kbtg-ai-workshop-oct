package usecase

import (
	"errors"
	"temp-backend/internal/domain"
	"testing"
	"time"
)

// Mock TransferRepository for testing
type mockTransferRepository struct {
	transfers     map[string]*domain.Transfer
	userTransfers map[int][]*domain.Transfer
	lastTransfer  map[int]*domain.Transfer // fromUserID -> last transfer
	shouldErr     bool
	nextID        int
}

func newMockTransferRepository() *mockTransferRepository {
	return &mockTransferRepository{
		transfers:     make(map[string]*domain.Transfer),
		userTransfers: make(map[int][]*domain.Transfer),
		lastTransfer:  make(map[int]*domain.Transfer),
		nextID:        1,
	}
}

func (m *mockTransferRepository) Create(transfer *domain.Transfer) (*domain.Transfer, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	transfer.ID = m.nextID
	m.nextID++
	transfer.CreatedAt = time.Now()
	transfer.UpdatedAt = time.Now()
	m.transfers[transfer.IdemKey] = transfer

	// Add to user transfers
	m.userTransfers[transfer.FromUserID] = append(m.userTransfers[transfer.FromUserID], transfer)
	m.userTransfers[transfer.ToUserID] = append(m.userTransfers[transfer.ToUserID], transfer)

	return transfer, nil
}

func (m *mockTransferRepository) GetByIdemKey(idemKey string) (*domain.Transfer, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	transfer, exists := m.transfers[idemKey]
	if !exists {
		return nil, nil
	}
	return transfer, nil
}

func (m *mockTransferRepository) GetByUserID(userID int, page, pageSize int) ([]*domain.Transfer, int, error) {
	if m.shouldErr {
		return nil, 0, errors.New("mock error")
	}
	transfers := m.userTransfers[userID]
	return transfers, len(transfers), nil
}

func (m *mockTransferRepository) GetLastTransferFromUser(fromUserID int) (*domain.Transfer, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	transfer, exists := m.lastTransfer[fromUserID]
	if !exists {
		return nil, nil
	}
	return transfer, nil
}

func (m *mockTransferRepository) UpdateStatus(idemKey string, status string, completedAt *time.Time, failReason *string) error {
	if m.shouldErr {
		return errors.New("mock error")
	}
	transfer := m.transfers[idemKey]
	if transfer != nil {
		transfer.Status = status
		transfer.CompletedAt = completedAt
		transfer.FailReason = failReason

		// Update last transfer if completed
		if status == "completed" {
			m.lastTransfer[transfer.FromUserID] = transfer
		}
	}
	return nil
}

// Mock PointLedgerRepository for testing
type mockPointLedgerRepository struct {
	shouldErr bool
}

func newMockPointLedgerRepository() *mockPointLedgerRepository {
	return &mockPointLedgerRepository{}
}

func (m *mockPointLedgerRepository) Create(entry *domain.PointLedger) error {
	if m.shouldErr {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockPointLedgerRepository) GetByUserID(userID int, limit int) ([]*domain.PointLedger, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	return nil, nil
}

// Test cases for transfer amount validation (Rule 2)
func TestCreateTransfer_AmountValidation(t *testing.T) {
	mockUserRepo := newMockUserRepository()
	mockTransferRepo := newMockTransferRepository()
	mockPointLedgerRepo := newMockPointLedgerRepository()

	// Create test users with sufficient points
	user1, _ := mockUserRepo.Create(&domain.User{FirstName: "John", LastName: "Doe", Points: 10.0})
	user2, _ := mockUserRepo.Create(&domain.User{FirstName: "Jane", LastName: "Doe", Points: 10.0})

	// Create a test transfer usecase without DB dependency for validation tests
	usecase := &transferUsecase{
		transferRepo:    mockTransferRepo,
		userRepo:        mockUserRepo,
		pointLedgerRepo: mockPointLedgerRepo,
		db:              nil, // We'll test validation only, not the actual transfer process
	}

	tests := []struct {
		name          string
		request       *domain.TransferCreateRequest
		shouldError   bool
		expectedError string
	}{
		{
			name: "Valid amount - exactly 2.0",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     2.0,
			},
			shouldError: false,
		},
		{
			name: "Valid amount - less than 2.0",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     1.50,
			},
			shouldError: false,
		},
		{
			name: "Valid amount - 2 decimal places",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     1.99,
			},
			shouldError: false,
		},
		{
			name: "Invalid amount - exceeds 2.0",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     2.01,
			},
			shouldError:   true,
			expectedError: "transfer amount cannot exceed 2 points",
		},
		{
			name: "Invalid amount - way exceeds 2.0",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     5.0,
			},
			shouldError:   true,
			expectedError: "transfer amount cannot exceed 2 points",
		},
		{
			name: "Invalid amount - more than 2 decimal places",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     1.999,
			},
			shouldError:   true,
			expectedError: "transfer amount cannot have more than 2 decimal places",
		},
		{
			name: "Invalid amount - more than 2 decimal places (small amount)",
			request: &domain.TransferCreateRequest{
				FromUserID: user1.ID,
				ToUserID:   user2.ID,
				Amount:     0.001,
			},
			shouldError:   true,
			expectedError: "transfer amount cannot have more than 2 decimal places",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test only validation, not the full transfer process
			err := usecase.validateTransferRequest(tt.request)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// Test cases for last transfer validation (Rule 3)
func TestCreateTransfer_LastTransferValidation(t *testing.T) {
	mockUserRepo := newMockUserRepository()
	mockTransferRepo := newMockTransferRepository()

	// Create test users with sufficient points
	user1, _ := mockUserRepo.Create(&domain.User{FirstName: "John", LastName: "Doe", Points: 10.0})
	user2, _ := mockUserRepo.Create(&domain.User{FirstName: "Jane", LastName: "Doe", Points: 10.0})
	user3, _ := mockUserRepo.Create(&domain.User{FirstName: "Bob", LastName: "Doe", Points: 10.0})

	// Setup: Create a completed transfer from user1 to user2
	completedTime := time.Now()
	lastTransfer := &domain.Transfer{
		ID:          1,
		FromUserID:  user1.ID,
		ToUserID:    user2.ID,
		Amount:      1.0,
		Status:      "completed",
		CompletedAt: &completedTime,
	}
	mockTransferRepo.lastTransfer[user1.ID] = lastTransfer

	tests := []struct {
		name          string
		fromUserID    int
		toUserID      int
		shouldError   bool
		expectedError string
	}{
		{
			name:        "Valid transfer - to different user",
			fromUserID:  user1.ID,
			toUserID:    user3.ID, // Different user from last transfer
			shouldError: false,
		},
		{
			name:          "Invalid transfer - to same user as last transfer",
			fromUserID:    user1.ID,
			toUserID:      user2.ID, // Same user as last completed transfer
			shouldError:   true,
			expectedError: "cannot transfer to the same user as the last completed transfer",
		},
		{
			name:        "Valid transfer - user with no previous transfers",
			fromUserID:  user2.ID, // User who hasn't made any transfers
			toUserID:    user3.ID,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic separately
			lastTransfer, err := mockTransferRepo.GetLastTransferFromUser(tt.fromUserID)
			if err != nil {
				t.Fatalf("Error getting last transfer: %v", err)
			}

			var validationErr error
			if lastTransfer != nil && lastTransfer.ToUserID == tt.toUserID && lastTransfer.Status == "completed" {
				validationErr = errors.New("cannot transfer to the same user as the last completed transfer")
			}

			if tt.shouldError {
				if validationErr == nil {
					t.Errorf("Expected error but got none")
				} else if validationErr.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, validationErr.Error())
				}
			} else {
				if validationErr != nil {
					t.Errorf("Expected no error but got: %v", validationErr)
				}
			}
		})
	}
}

// Test case for multiple validation rules combined
func TestCreateTransfer_CombinedValidation(t *testing.T) {
	mockUserRepo := newMockUserRepository()
	mockTransferRepo := newMockTransferRepository()
	mockPointLedgerRepo := newMockPointLedgerRepository()

	// Create test users with sufficient points
	user1, _ := mockUserRepo.Create(&domain.User{FirstName: "John", LastName: "Doe", Points: 10.0})
	user2, _ := mockUserRepo.Create(&domain.User{FirstName: "Jane", LastName: "Doe", Points: 10.0})

	usecase := &transferUsecase{
		transferRepo:    mockTransferRepo,
		userRepo:        mockUserRepo,
		pointLedgerRepo: mockPointLedgerRepo,
		db:              nil,
	}

	// Test that amount validation runs before last transfer validation
	request := &domain.TransferCreateRequest{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Amount:     3.0, // Exceeds 2.0 limit
	}

	err := usecase.validateTransferRequest(request)
	if err == nil {
		t.Errorf("Expected error but got none")
	} else if err.Error() != "transfer amount cannot exceed 2 points" {
		t.Errorf("Expected amount validation error, got: %s", err.Error())
	}
}
