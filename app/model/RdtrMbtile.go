package model

import (
	"time"
)

type RdtrMbtile struct {
	Id        int64      `gorm:"primaryKey" json:"id,omitempty"`
	File_name string     `gorm:"column:file_name" json:"file_name"`
	UUID      string     `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rdtr_id   NullInt64  `gorm:"column:rdtr_id" json:"rdtr_id,omitempty"`
	Rdtr      *RdtrType  `gorm:"foreignKey:rdtr_id" json:"rdtr,omitempty"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CheckedAt *time.Time `gorm:"column:checked_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"checked_at,omitempty"`
}

// Set the table name for the User model
func (RdtrMbtile) TableName() string {
	return "rdtr_mbtile" // Replace with your existing table name
}
