package db

import (
	"time"
)

// User represents the user table
type User struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	Nickname  string    `xorm:"varchar(100) notnull 'nickname'" json:"nickname"`
	Email     string    `xorm:"varchar(255) notnull unique 'email'" json:"email"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'" json:"updatedAt"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "user"
}

// UserDevice represents the user_device table
type UserDevice struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	UserID    int64     `xorm:"notnull 'user_id'" json:"userId"`
	DeviceID  string    `xorm:"varchar(255) notnull 'device_id'" json:"deviceId"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'" json:"updatedAt"`
}

// TableName returns the table name for UserDevice
func (UserDevice) TableName() string {
	return "user_device"
}