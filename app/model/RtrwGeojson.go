package model

import (
	"math/big"
	"time"
)

type RtrwGeojson struct {
	Id           *big.Int  `gorm:"primaryKey" json:"id"`
	Rtrw_id      *big.Int  `gorm:"column:rtrw_id" json:"rtrw_id"`
	Rtrw_file_id *big.Int  `gorm:"column:rtrw_file_id" json:"rtrw_file_id"`
	Geometry     string    `gorm:"type:geometry(GEOMETRY,4326)" json:"geometry"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (c *RtrwGeojson) TableName() string {
	return "rtrw_geojson" // Replace with your existing table name
}
