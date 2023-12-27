package model

import (
	"time"
)

type RegRegency struct {
	// gorm.Model
	Id            int64     `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"varchar(255)" json:"name"`
	FullName      string    `gorm:"column:full_name" json:"full_name"`             // Assuming full_name is the column name in the database
	RegProvinceID int64     `gorm:"column:reg_province_id" json:"reg_province_id"` // Assuming reg_province_id is the column name in the database
	Type          string    `gorm:"type:varchar(255)" json:"type"`                 // Example type; adjust as per your requirements
	Latitude      float64   `gorm:"column:latitude" json:"latitude"`               // Assuming latitude is a float64; adjust if needed
	Longitude     float64   `gorm:"column:longitude" json:"longitude"`             // Assuming longitude is a float64; adjust if needed
	AltName       string    `gorm:"column:alt_name" json:"alt_name"`               // Assuming alt_name is a string; adjust if needed
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RegRegency) TableName() string {
	return "reg_regency" // Replace with your existing table name
}
