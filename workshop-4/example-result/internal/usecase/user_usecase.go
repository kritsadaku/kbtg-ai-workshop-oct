package usecase

import (
	"errors"
	"strings"
	"temp-backend/internal/domain"
	"unicode/utf8"
)

// userUsecase implements domain.UserUsecase
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users
func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.userRepo.GetAll()
}

// GetUserByID retrieves a user by ID
func (u *userUsecase) GetUserByID(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CreateUser creates a new user
func (u *userUsecase) CreateUser(req *domain.CreateUserRequest) (*domain.User, error) {
	// Validation
	if err := u.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Set default membership level if not provided
	if req.MembershipLevel == "" {
		req.MembershipLevel = "Bronze"
	}

	// Create user entity
	user := &domain.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Email:           req.Email,
		MembershipLevel: req.MembershipLevel,
		Points:          req.Points,
	}

	return u.userRepo.Create(user)
}

// UpdateUser updates user information
func (u *userUsecase) UpdateUser(id int, req *domain.UpdateUserRequest) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Check if user exists
	exists, err := u.userRepo.Exists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Validate update request
	if err := u.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.MembershipLevel != nil {
		updates["membership_level"] = *req.MembershipLevel
	}
	if req.Points != nil {
		updates["points"] = *req.Points
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	return u.userRepo.Update(id, updates)
}

// DeleteUser removes a user
func (u *userUsecase) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	exists, err := u.userRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	return u.userRepo.Delete(id)
}

// validateCreateRequest validates the create user request
func (u *userUsecase) validateCreateRequest(req *domain.CreateUserRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}

	// Validate name length (should not exceed 3 characters)
	if utf8.RuneCountInString(strings.TrimSpace(req.FirstName)) > 3 {
		return errors.New("first name must not exceed 3 characters")
	}
	if utf8.RuneCountInString(strings.TrimSpace(req.LastName)) > 3 {
		return errors.New("last name must not exceed 3 characters")
	}

	// Validate membership level
	validLevels := map[string]bool{
		"Bronze":   true,
		"Silver":   true,
		"Gold":     true,
		"Platinum": true,
	}
	if req.MembershipLevel != "" && !validLevels[req.MembershipLevel] {
		return errors.New("invalid membership level")
	}

	return nil
}

// validateUpdateRequest validates the update user request
func (u *userUsecase) validateUpdateRequest(req *domain.UpdateUserRequest) error {
	// Validate name length (should not exceed 3 characters)
	if req.FirstName != nil && utf8.RuneCountInString(strings.TrimSpace(*req.FirstName)) > 3 {
		return errors.New("first name must not exceed 3 characters")
	}
	if req.LastName != nil && utf8.RuneCountInString(strings.TrimSpace(*req.LastName)) > 3 {
		return errors.New("last name must not exceed 3 characters")
	}

	// Validate membership level
	if req.MembershipLevel != nil {
		validLevels := map[string]bool{
			"Bronze":   true,
			"Silver":   true,
			"Gold":     true,
			"Platinum": true,
		}
		if !validLevels[*req.MembershipLevel] {
			return errors.New("invalid membership level")
		}
	}

	return nil
}
