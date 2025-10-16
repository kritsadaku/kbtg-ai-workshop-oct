package usecase

import (
	"errors"
	"temp-backend/internal/domain"
	"testing"
)

// Mock UserRepository for testing
type mockUserRepository struct {
	users     map[int]*domain.User
	points    map[int]float64
	nextID    int
	shouldErr bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[int]*domain.User),
		points: make(map[int]float64),
		nextID: 1,
	}
}

func (m *mockUserRepository) GetAll() ([]*domain.User, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	var result []*domain.User
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *mockUserRepository) GetByID(id int) (*domain.User, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) Create(user *domain.User) (*domain.User, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.points[user.ID] = user.Points
	return user, nil
}

func (m *mockUserRepository) Update(id int, updates map[string]interface{}) (*domain.User, error) {
	if m.shouldErr {
		return nil, errors.New("mock error")
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Apply updates
	if firstName, ok := updates["first_name"]; ok {
		user.FirstName = firstName.(string)
	}
	if lastName, ok := updates["last_name"]; ok {
		user.LastName = lastName.(string)
	}
	if points, ok := updates["points"]; ok {
		user.Points = points.(float64)
		m.points[id] = points.(float64)
	}

	m.users[id] = user
	return user, nil
}

func (m *mockUserRepository) Delete(id int) error {
	if m.shouldErr {
		return errors.New("mock error")
	}
	delete(m.users, id)
	delete(m.points, id)
	return nil
}

func (m *mockUserRepository) Exists(id int) (bool, error) {
	if m.shouldErr {
		return false, errors.New("mock error")
	}
	_, exists := m.users[id]
	return exists, nil
}

func (m *mockUserRepository) UpdatePoints(userID int, newPoints float64) error {
	if m.shouldErr {
		return errors.New("mock error")
	}
	if user, exists := m.users[userID]; exists {
		user.Points = newPoints
		m.points[userID] = newPoints
		return nil
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) GetUserPoints(userID int) (float64, error) {
	if m.shouldErr {
		return 0, errors.New("mock error")
	}
	points, exists := m.points[userID]
	if !exists {
		return 0, errors.New("user not found")
	}
	return points, nil
}

// Test cases for user name validation (Rule 1)
func TestCreateUser_NameLengthValidation(t *testing.T) {
	mockRepo := newMockUserRepository()
	usecase := NewUserUsecase(mockRepo)

	tests := []struct {
		name          string
		request       *domain.CreateUserRequest
		shouldError   bool
		expectedError string
	}{
		{
			name: "Valid names - exactly 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "จิน", // 3 Thai characters without tone marks
				LastName:  "ดอน", // 3 Thai characters
				Phone:     "0123456789",
				Email:     "john@example.com",
			},
			shouldError: false,
		},
		{
			name: "Valid names - less than 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "จน", // 2 Thai characters
				LastName:  "ด",  // 1 Thai character
				Phone:     "0123456789",
				Email:     "john@example.com",
			},
			shouldError: false,
		},
		{
			name: "Invalid first name - more than 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "จิรายุ", // 6 Thai characters
				LastName:  "ดอน",    // 3 Thai characters
				Phone:     "0123456789",
				Email:     "john@example.com",
			},
			shouldError:   true,
			expectedError: "first name must not exceed 3 characters",
		},
		{
			name: "Invalid last name - more than 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "จิน",        // 3 Thai characters
				LastName:  "ดอนเมเยอร์", // 7 Thai characters
				Phone:     "0123456789",
				Email:     "john@example.com",
			},
			shouldError:   true,
			expectedError: "last name must not exceed 3 characters",
		},
		{
			name: "Valid English names - exactly 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "Tom",
				LastName:  "Lee",
				Phone:     "0123456789",
				Email:     "tom@example.com",
			},
			shouldError: false,
		},
		{
			name: "Invalid English first name - more than 3 characters",
			request: &domain.CreateUserRequest{
				FirstName: "Thomas",
				LastName:  "Lee",
				Phone:     "0123456789",
				Email:     "tom@example.com",
			},
			shouldError:   true,
			expectedError: "first name must not exceed 3 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usecase.CreateUser(tt.request)

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

func TestUpdateUser_NameLengthValidation(t *testing.T) {
	mockRepo := newMockUserRepository()
	usecase := NewUserUsecase(mockRepo)

	// Create a user first
	user, _ := usecase.CreateUser(&domain.CreateUserRequest{
		FirstName: "จอน",
		LastName:  "ดู",
		Phone:     "0123456789",
		Email:     "john@example.com",
	})

	tests := []struct {
		name          string
		request       *domain.UpdateUserRequest
		shouldError   bool
		expectedError string
	}{
		{
			name: "Valid update - names within 3 characters",
			request: &domain.UpdateUserRequest{
				FirstName: stringPtr("จิน"), // 3 Thai characters
				LastName:  stringPtr("ดอน"), // 3 Thai characters
			},
			shouldError: false,
		},
		{
			name: "Invalid update - first name more than 3 characters",
			request: &domain.UpdateUserRequest{
				FirstName: stringPtr("จิรายุ"), // 6 Thai characters
			},
			shouldError:   true,
			expectedError: "first name must not exceed 3 characters",
		},
		{
			name: "Invalid update - last name more than 3 characters",
			request: &domain.UpdateUserRequest{
				LastName: stringPtr("ดอนเมเยอร์"), // 7 Thai characters
			},
			shouldError:   true,
			expectedError: "last name must not exceed 3 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usecase.UpdateUser(user.ID, tt.request)

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

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
