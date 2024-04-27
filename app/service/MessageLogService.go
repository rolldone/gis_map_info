package service

import (
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type MessageLogService struct {
	DB *gorm.DB
	// Custom own payload
	Asynq_uuid string
}

type MessageLogAdd struct {
	Asynq_uuid string
}

func (c *MessageLogService) Add(message string) (Model.MessageLog, error) {
	messageLog := Model.MessageLog{}
	messageLog.Asynq_uuid = c.Asynq_uuid
	if err := c.DB.Create(&messageLog).Error; err != nil {
		return Model.MessageLog{}, err
	}
	return messageLog, nil
}

func (c *MessageLogService) Gets(take int64) *gorm.DB {
	return c.DB.Model(&Model.MessageLog{})
}
