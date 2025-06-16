package repository

import "blog-fanchiikawa-service/db"

// UserRepository defines the interface for user data access
type UserRepository interface {
	// GetByEmail retrieves a user by email address
	GetByEmail(email string) (*db.User, error)
	
	// GetByID retrieves a user by ID
	GetByID(id int64) (*db.User, error)
	
	// Create creates a new user
	Create(user *db.User) error
	
	// List retrieves users with pagination
	List(limit int, offset int) ([]*db.User, error)
}

// UserDeviceRepository defines the interface for user device data access
type UserDeviceRepository interface {
	// Create creates a new user device
	Create(device *db.UserDevice) error
	
	// GetByUserID retrieves devices for a user
	GetByUserID(userID int64) ([]*db.UserDevice, error)
}

// TransactionManager defines the interface for transaction management
type TransactionManager interface {
	// WithTransaction executes a function within a database transaction
	WithTransaction(fn func() error) error
}