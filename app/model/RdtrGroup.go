package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RdtrGroup struct {
	Id         int64            `gorm:"primaryKey" json:"id,omitempty"`
	Rdtr_id    int64            `gorm:"column:rdtr_id" json:"rdtr_id"`
	Properties pgtype.JSONCodec `gorm:"column:properties;type:json" json:"properties"`
	Status     string           `gorm:"column:status;type:varchar" json:"status"`
	Name       string           `gorm:"column:name;type:varchar(255)" json:"name"`
	Cat_key    string           `gorm:"column:cat_key;type:varchar(255)" json:"cat_key"`
	CreatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RdtrGroup) TableName() string {
	return "rdtr_group" // Replace with your existing table name
}
