package model

import (
	"time"
)

type RegDistrict struct {
	Id           int64     `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"varchar(255)" json:"name"`
	RegRegencyID int64     `gorm:"column:reg_regency_id" json:"reg_regency_id"` // Assuming reg_regency_id is the column name in the database and it's an int64
	Latitude     float64   `gorm:"column:latitude" json:"latitude"`             // Assuming latitude is a float64; adjust if needed
	Longitude    float64   `gorm:"column:longitude" json:"longitude"`           // Assuming longitude is a float64; adjust if needed
	AltName      string    `gorm:"column:alt_name" json:"alt_name"`             // Assuming alt_name is a string; adjust if needed
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RegDistrict) TableName() string {
	return "reg_district" // Replace with your existing table name
}
