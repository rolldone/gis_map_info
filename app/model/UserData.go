package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         int64           `gorm:"column:id" json:"id"`
	Uuid       string          `gorm:"column:uuid" json:"uuid"`
	Name       string          `gorm:"column:name" json:"name"`
	Username   string          `gorm:"column:username" json:"username"`
	Email      string          `gorm:"column:email" json:"email"`
	Passkey    string          `gorm:"column:passkey" json:"passkey"`
	Salt       string          `gorm:"column:salt" json:"salt"`
	Status     string          `gorm:"column:status" json:"status"`
	Created_at time.Time       `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	Updated_at time.Time       `gorm:"column:updated_at;type:timestamp;autoUpdateTime:true" json:"updated_at"`
	Deleted_at *gorm.DeletedAt `gorm:"type:timestamp" json:"deleted_at,omitempty"`
}

type UserView struct {
	Id         int64           `gorm:"column:id" json:"id"`
	Uuid       string          `gorm:"column:uuid" json:"uuid"`
	Name       string          `gorm:"column:name" json:"name"`
	Username   string          `gorm:"column:username" json:"username"`
	Email      string          `gorm:"column:email" json:"email"`
	Status     string          `gorm:"column:status" json:"status"`
	Created_at time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	Updated_at time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Deleted_at *gorm.DeletedAt `gorm:"type:timestamp" json:"deleted_at,omitempty"`
	Password   *string         `gorm:"-" json:"password,omitempty"`
}

// Set the table name for the User model
func (c *User) TableName() string {
	return "user_data" // Replace with your existing table name
}
