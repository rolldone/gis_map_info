package model

import (
	"time"
)

type RegVillage struct {
	Id            int64     `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"varchar(255)" json:"name"`
	RegDistrictID int64     `gorm:"column:reg_district_id" json:"reg_district_id"` // Assuming reg_district_id is the column name in the database and it's an int64
	Latitude      float64   `gorm:"column:latitude" json:"latitude"`               // Assuming latitude is a float64; adjust if needed
	Longitude     float64   `gorm:"column:longitude" json:"longitude"`             // Assuming longitude is a float64; adjust if needed
	AltName       string    `gorm:"column:alt_name" json:"alt_name"`               // Assuming alt_name is a string; adjust if needed
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RegVillage) TableName() string {
	return "reg_village" // Replace with your existing table name
}
