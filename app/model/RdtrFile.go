package model

import (
	"math/big"
	"time"
)

type RdtrFile struct {
	Id            *big.Int  `gorm:"primaryKey" json:"id"`
	UUID          string    `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Rdtr_group_id *big.Int  `gorm:"column:rdtr_group_id" json:"rdtr_group_id"`
	Rdtr_id       *big.Int  `gorm:"column:rdtr_id" json:"rdtr_id"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RdtrFile) TableName() string {
	return "rdtr_group" // Replace with your existing table name
}
