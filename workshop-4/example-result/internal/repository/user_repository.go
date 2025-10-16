package repository

import (
	"database/sql"
	"temp-backend/internal/domain"
	"time"
)

// userRepository implements domain.UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetAll retrieves all users from database
func (r *userRepository) GetAll() ([]*domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, first_name, last_name, phone, email, member_since, 
			   membership_level, points, created_at, updated_at 
		FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Phone,
			&user.Email, &user.MemberSince, &user.MembershipLevel,
			&user.Points, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		if user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			return nil, err
		}
		if user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	var createdAt, updatedAt string

	err := r.db.QueryRow(`
		SELECT id, first_name, last_name, phone, email, member_since,
			   membership_level, points, created_at, updated_at
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Phone,
		&user.Email, &user.MemberSince, &user.MembershipLevel,
		&user.Points, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	if user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		return nil, err
	}
	if user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	memberSince := time.Now().Format("2/1/2006") // Thai date format

	result, err := r.db.Exec(`
		INSERT INTO users (first_name, last_name, phone, email, member_since, 
						  membership_level, points, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.FirstName, user.LastName, user.Phone, user.Email,
		memberSince, user.MembershipLevel, user.Points, now, now)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

// Update updates user fields
func (r *userRepository) Update(id int, updates map[string]interface{}) (*domain.User, error) {
	if len(updates) == 0 {
		return r.GetByID(id)
	}

	// Build dynamic query
	setParts := make([]string, 0, len(updates)+1)
	args := make([]interface{}, 0, len(updates)+2)

	for field, value := range updates {
		setParts = append(setParts, field+" = ?")
		args = append(args, value)
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	args = append(args, id)

	query := "UPDATE users SET "
	for i, part := range setParts {
		if i > 0 {
			query += ", "
		}
		query += part
	}
	query += " WHERE id = ?"

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

// Delete removes a user
func (r *userRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// Exists checks if user exists
func (r *userRepository) Exists(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", id).Scan(&exists)
	return exists, err
}

// UpdatePoints updates user points
func (r *userRepository) UpdatePoints(userID int, newPoints float64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := r.db.Exec("UPDATE users SET points = ?, updated_at = ? WHERE id = ?", newPoints, now, userID)
	return err
}

// GetUserPoints gets current user points
func (r *userRepository) GetUserPoints(userID int) (float64, error) {
	var points float64
	err := r.db.QueryRow("SELECT points FROM users WHERE id = ?", userID).Scan(&points)
	return points, err
}
