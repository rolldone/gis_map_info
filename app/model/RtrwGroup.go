package model

import (
	"time"

	"gorm.io/datatypes"
)

type RtrwGroup struct {
	Id         int64          `gorm:"primaryKey" json:"id,omitempty"`
	Uuid       string         `gorm:"column:uuid" json:"uuid,omitempty"`
	Rtrw_id    int64          `gorm:"column:rtrw_id" json:"rtrw_id"`
	Properties datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	Status     string         `gorm:"column:status;type:varchar" json:"status"`
	Name       string         `gorm:"column:name;type:varchar(255)" json:"name"`
	Asset_key  string         `gorm:"column:asset_key;type:varchar(255)" json:"asset_key"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Datas      []RtrwFile     `gorm:"foreignKey:rtrw_group_id" json:"datas"`
}

type RtrwGroupView struct {
	RtrwGroup
	Unvalidated NullInt64 `gorm:"column:unvalidated" json:"unvalidated,omitempty"`
	Validated   NullInt64 `gorm:"column:validated" json:"validated,omitempty"`
}

// Set the table name for the User model
func (RtrwGroup) TableName() string {
	return "rtrw_group" // Replace with your existing table name
}
