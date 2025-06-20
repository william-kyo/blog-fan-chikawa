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

type ImageRepository interface {
	Create(image *db.Image) error

	GetByID(id int64) (*db.Image, error)

	UpdateLabelDetected(id int64, labelDetected bool) (int64, error)

	GetByLabelDetected(labelDetected bool) ([]*db.Image, error)
}

type LabelRepository interface {
	Create(label *db.Label) error

	GetByName(labelName string) (*db.Label, error)
}

type ImageLabelRepository interface {
	Create(imageLabel *db.ImageLabel) error
	GetByImageAndLabel(imageID, labelID int64) (*db.ImageLabel, error)
}

// TransactionManager defines the interface for transaction management
type TransactionManager interface {
	// WithTransaction executes a function within a database transaction
	WithTransaction(fn func() error) error
}
