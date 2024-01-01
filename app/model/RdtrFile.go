package model

import (
	"time"
)

type RdtrFile struct {
	Id            int64     `gorm:"primaryKey" json:"id,omitempty"`
	UUID          string    `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rdtr_group_id NullInt64 `gorm:"column:rdtr_group_id" json:"rdtr_group_id,omitempty"`
	Rdtr_id       NullInt64 `gorm:"column:rdtr_id" json:"rdtr_id,omitempty"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RdtrFile) TableName() string {
	return "rdtr_file" // Replace with your existing table name
}
