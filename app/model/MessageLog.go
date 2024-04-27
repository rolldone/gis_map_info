package model

import (
	"time"
)

type MessageLog struct {
	Asynq_uuid string    `gorm:"column:asynq_uuid" json:"asynq_uuid"`
	Data_log   string    `gorm:"column:data_log" json:"data_log"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type MessageLogView struct {
	MessageLog
}

// Set the table name for the User model
func (c *MessageLog) TableName() string {
	return "message_log" // Replace with your existing table name
}
