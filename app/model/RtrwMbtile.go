package model

import (
	"time"
)

type RtrwMbtile struct {
	Id        int64      `gorm:"primaryKey" json:"id,omitempty"`
	File_name string     `gorm:"column:file_name" json:"file_name"`
	UUID      string     `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rtrw_id   NullInt64  `gorm:"column:rtrw_id" json:"rtrw_id,omitempty"`
	Rtrw      *RtrwType  `gorm:"foreignKey:rtrw_id" json:"rtrw,omitempty"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CheckedAt *time.Time `gorm:"column:checked_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"checked_at,omitempty"`
}

// Set the table name for the User model
func (RtrwMbtile) TableName() string {
	return "rtrw_mbtile" // Replace with your existing table name
}
