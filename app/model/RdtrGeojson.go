package model

import (
	"math/big"
	"time"
)

type RdtrGeojson struct {
	Id           *big.Int  `gorm:"primaryKey" json:"id"`
	Rdtr_id      *big.Int  `gorm:"column:rdtr_id" json:"rdtr_id"`
	Rdtr_file_id *big.Int  `gorm:"column:rdtr_file_id" json:"rdtr_file_id"`
	Geometry     string    `gorm:"type:geometry(GEOMETRY,4326)" json:"geometry"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (c *RdtrGeojson) TableName() string {
	return "rdtr_geojson" // Replace with your existing table name
}
