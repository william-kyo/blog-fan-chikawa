package repository

import (
	"blog-fanchiikawa-service/db"
)

// userDeviceRepository implements UserDeviceRepository interface
type userDeviceRepository struct{}

// NewUserDeviceRepository creates a new UserDeviceRepository instance
func NewUserDeviceRepository() UserDeviceRepository {
	return &userDeviceRepository{}
}

// Create creates a new user device
func (r *userDeviceRepository) Create(device *db.UserDevice) error {
	_, err := db.Engine.Insert(device)
	return err
}

// GetByUserID retrieves devices for a user
func (r *userDeviceRepository) GetByUserID(userID int64) ([]*db.UserDevice, error) {
	var devices []*db.UserDevice
	err := db.Engine.Where("user_id = ?", userID).Find(&devices)
	return devices, err
}