package repository

import (
	"blog-fanchiikawa-service/db"
)

// userRepository implements UserRepository interface
type userRepository struct{}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// GetByEmail retrieves a user by email address
func (r *userRepository) GetByEmail(email string) (*db.User, error) {
	var user db.User
	has, err := db.Engine.Where("email = ?", email).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil // User not found
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*db.User, error) {
	var user db.User
	has, err := db.Engine.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil // User not found
	}
	return &user, nil
}

// Create creates a new user
func (r *userRepository) Create(user *db.User) error {
	_, err := db.Engine.Insert(user)
	return err
}

// List retrieves users with pagination
func (r *userRepository) List(limit int, offset int) ([]*db.User, error) {
	var users []*db.User
	err := db.Engine.Limit(limit, offset).Find(&users)
	return users, err
}