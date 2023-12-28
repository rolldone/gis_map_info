package model

import (
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RdtrGroup struct {
	Id         *big.Int         `gorm:"primaryKey" json:"id"`
	Rdtr_id    *big.Int         `gorm:"column:rdtr_id" json:"rdtr_id"`
	Properties pgtype.JSONCodec `gorm:"column:properties;type:json" json:"properties"`
	Status     string           `gorm:"column:status;type:varchar" json:"status"`
	Name       string           `gorm:"column:name;type:varchar" json:"name"`
	Cat_key    string           `gorm:"column:type:varchar" json:"cat_key"`
	CreatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RdtrGroup) TableName() string {
	return "rdtr_group" // Replace with your existing table name
}
