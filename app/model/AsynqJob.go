package model

import (
	"time"

	"gorm.io/datatypes"
)

type AsynqJob struct {
	App_uuid     string         `gorm:"column:app_uuid"`
	Asynq_uuid   string         `gorm:"column:asynq_uuid"`
	Table_name   string         `gorm:"column:table_name"`
	Status       string         `gorm:"column:status"`
	Message_text string         `gorm:"column:message_text"`
	Order_number NullInt64      `gorm:"column:order_number"`
	Payload      datatypes.JSON `gorm:"column:payload"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type AsyncJobView struct {
	AsynqJob
}

// Set the table name for the User model
func (c *AsynqJob) TableName() string {
	return "asynq_job" // Replace with your existing table name
}
