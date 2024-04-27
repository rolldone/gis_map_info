package model

import (
	"time"
)

type ZlpMbtile struct {
	Id              int64      `gorm:"primaryKey" json:"id,omitempty"`
	File_name       string     `gorm:"column:file_name" json:"file_name"`
	UUID            string     `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Zlp_id          NullInt64  `gorm:"column:zlp_id" json:"zlp_id,omitempty"`
	Zlp             *ZlpType   `gorm:"foreignKey:zlp_id" json:"zlp,omitempty"`
	Asset_key       string     `gorm:"column:asset_key" json:"asset_key"`
	Reg_province_id int64      `gorm:"column:reg_province_id" json:"reg_province_id,omitempty"`
	Zlp_group_id    int64      `gorm:"column:zlp_group_id" json:"zlp_group_id,omitempty"`
	CreatedAt       time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CheckedAt       *time.Time `gorm:"column:checked_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"checked_at,omitempty"`
}

// Set the table name for the User model
func (ZlpMbtile) TableName() string {
	return "zlp_mbtile" // Replace with your existing table name
}
