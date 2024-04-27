package model

import (
	"time"
)

type RtrwFile struct {
	Id            int64      `gorm:"primaryKey" json:"id,omitempty"`
	UUID          string     `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rtrw_group_id NullInt64  `gorm:"column:rtrw_group_id" json:"rtrw_group_id,omitempty"`
	Rtrw_id       NullInt64  `gorm:"column:rtrw_id" json:"rtrw_id,omitempty"`
	Rtrw_group    *RtrwGroup `gorm:"foreignKey:rtrw_group_id" json:"rtrw_group,omitempty"`
	Rtrw          *RtrwType  `gorm:"foreignKey:rtrw_id" json:"rtrw,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ValidatedAt   *time.Time `gorm:"column:validated_at;type:timestamp;" json:"validated_at"`
}

// Set the table name for the User model
func (RtrwFile) TableName() string {
	return "rtrw_file" // Replace with your existing table name
}
