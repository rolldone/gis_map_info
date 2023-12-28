package model

import (
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RtrwGroup struct {
	Id         *big.Int         `gorm:"primaryKey" json:"id"`
	Rtrw_id    *big.Int         `gorm:"column:rtrw_id" json:"rtrw_id"`
	Properties pgtype.JSONCodec `gorm:"type:json" json:"properties"`
	Status     string           `gorm:"type:varchar" json:"status"`
	Name       string           `gorm:"type:varchar" json:"name"`
	Cat_key    string           `gorm:"type:varchar" json:"cat_key"`
	CreatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RtrwGroup) TableName() string {
	return "rtrw_group" // Replace with your existing table name
}
