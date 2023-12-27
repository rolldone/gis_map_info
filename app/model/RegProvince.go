package model

import (
	"time"
)

type RegProvince struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"varchar(255)" json:"name"`
	Latitude  float64   `gorm:"column:latitude" json:"latitude"`   // Assuming latitude is a float64; adjust if needed
	Longitude float64   `gorm:"column:longitude" json:"longitude"` // Assuming longitude is a float64; adjust if needed
	AltName   string    `gorm:"column:alt_name" json:"alt_name"`   // Assuming alt_name is a string; adjust if needed
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	// DeletedAt gorm.DeletedAt
}

type Tabler interface {
	TableName() string
}

// Set the table name for the User model
func (c *RegProvince) TableName() string {
	return "reg_province" // Replace with your existing table name
}
