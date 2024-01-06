package model

import (
	"time"

	"gorm.io/datatypes"
)

type RdtrGeojson struct {
	Uuid          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rdtr_id       NullInt64      `gorm:"column:rdtr_id" json:"rdtr_id,omitempty"`
	Rdtr_group_id NullInt64      `gorm:"column:rdtr_group_id" json:"rdtr_group_id,omitempty"`
	Rdtr_file_id  NullInt64      `gorm:"column:rdtr_file_id" json:"rdtr_file_id,omitempty"`
	Geojson       string         `gorm:"type:geometry(GEOMETRY,4326)" json:"geojson"`
	Properties    datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (c *RdtrGeojson) TableName() string {
	return "rdtr_geojson" // Replace with your existing table name
}
