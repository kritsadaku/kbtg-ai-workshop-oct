package domain

import "time"

// User represents the user entity
type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	MemberSince     string    `json:"member_since"`
	MembershipLevel string    `json:"membership_level"`
	Points          float64   `json:"points"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	FirstName       string  `json:"first_name" validate:"required"`
	LastName        string  `json:"last_name" validate:"required"`
	Phone           string  `json:"phone" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	MembershipLevel string  `json:"membership_level"`
	Points          float64 `json:"points"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName       *string  `json:"first_name,omitempty"`
	LastName        *string  `json:"last_name,omitempty"`
	Phone           *string  `json:"phone,omitempty"`
	Email           *string  `json:"email,omitempty"`
	MembershipLevel *string  `json:"membership_level,omitempty"`
	Points          *float64 `json:"points,omitempty"`
}

// Transfer represents the transfer entity
type Transfer struct {
	ID          int        `json:"transferId,omitempty"`
	IdemKey     string     `json:"idemKey"`
	FromUserID  int        `json:"fromUserId"`
	ToUserID    int        `json:"toUserId"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	Note        *string    `json:"note,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	FailReason  *string    `json:"failReason,omitempty"`
}

// TransferCreateRequest represents the request to create a transfer
type TransferCreateRequest struct {
	FromUserID int     `json:"fromUserId" validate:"required,min=1"`
	ToUserID   int     `json:"toUserId" validate:"required,min=1"`
	Amount     float64 `json:"amount" validate:"required,min=0.01"`
	Note       *string `json:"note,omitempty"`
}

// TransferCreateResponse represents the response after creating a transfer
type TransferCreateResponse struct {
	Transfer *Transfer `json:"transfer"`
}

// TransferGetResponse represents the response for getting a transfer
type TransferGetResponse struct {
	Transfer *Transfer `json:"transfer"`
}

// TransferListResponse represents the response for listing transfers
type TransferListResponse struct {
	Data     []*Transfer `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
}

// PointLedger represents a point ledger entry
type PointLedger struct {
	ID           int     `json:"id"`
	UserID       int     `json:"userId"`
	Change       float64 `json:"change"`
	BalanceAfter float64 `json:"balanceAfter"`
	EventType    string  `json:"eventType"`
	TransferID   *int    `json:"transferId,omitempty"`
	Reference    *string `json:"reference,omitempty"`
	Metadata     *string `json:"metadata,omitempty"`
	CreatedAt    string  `json:"createdAt"`
}

// UserRepository interface defines the contract for user data operations
type UserRepository interface {
	GetAll() ([]*User, error)
	GetByID(id int) (*User, error)
	Create(user *User) (*User, error)
	Update(id int, updates map[string]interface{}) (*User, error)
	Delete(id int) error
	Exists(id int) (bool, error)
	UpdatePoints(userID int, newPoints float64) error
	GetUserPoints(userID int) (float64, error)
}

// TransferRepository interface defines the contract for transfer data operations
type TransferRepository interface {
	Create(transfer *Transfer) (*Transfer, error)
	GetByIdemKey(idemKey string) (*Transfer, error)
	GetByUserID(userID int, page, pageSize int) ([]*Transfer, int, error)
	GetLastTransferFromUser(fromUserID int) (*Transfer, error)
	UpdateStatus(idemKey string, status string, completedAt *time.Time, failReason *string) error
}

// PointLedgerRepository interface defines the contract for point ledger operations
type PointLedgerRepository interface {
	Create(entry *PointLedger) error
	GetByUserID(userID int, limit int) ([]*PointLedger, error)
}

// UserUsecase interface defines the business logic contract
type UserUsecase interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(req *CreateUserRequest) (*User, error)
	UpdateUser(id int, req *UpdateUserRequest) (*User, error)
	DeleteUser(id int) error
}

// TransferUsecase interface defines the transfer business logic contract
type TransferUsecase interface {
	CreateTransfer(req *TransferCreateRequest) (*Transfer, error)
	GetTransferByIdemKey(idemKey string) (*Transfer, error)
	GetTransfersByUserID(userID int, page, pageSize int) (*TransferListResponse, error)
}
