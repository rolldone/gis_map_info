package model

import (
	"time"
)

type ZlpFile struct {
	Id           int64      `gorm:"primaryKey" json:"id,omitempty"`
	UUID         string     `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Zlp_group_id NullInt64  `gorm:"column:zlp_group_id" json:"zlp_group_id,omitempty"`
	Zlp_id       NullInt64  `gorm:"column:zlp_id" json:"zlp_id,omitempty"`
	Zlp_group    *ZlpGroup  `gorm:"foreignKey:zlp_group_id" json:"zlp_group,omitempty"`
	Zlp          *ZlpType   `gorm:"foreignKey:zlp_id" json:"zlp,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ValidatedAt  *time.Time `gorm:"column:validated_at;type:timestamp;" json:"validated_at"`
}

// Set the table name for the User model
func (ZlpFile) TableName() string {
	return "zlp_file" // Replace with your existing table name
}
