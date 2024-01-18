package service

import (
	"context"
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/support/redis_support"
	"log"
	"os"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AsynqJobService struct {
	DB         *gorm.DB
	OS         *os.File
	Async_uuid *string
}

type AsynqJobAddPayload struct {
	App_uuid     string
	Asynq_uuid   string
	Payload      string
	Status       string
	Message_text string
	Table_name   string
}

type asyncJobStatus struct {
	STATUS_PENDING   string
	STATUS_PROCESS   string
	STATUS_FAILED    string
	STATUS_COMPLETED string
	STATUS_STOPPED   string
}

func (c *AsynqJobService) GetStatus() asyncJobStatus {
	return asyncJobStatus{
		STATUS_PENDING:   "PENDING",
		STATUS_PROCESS:   "PROCESS",
		STATUS_FAILED:    "FAILED",
		STATUS_COMPLETED: "COMPLETED",
		STATUS_STOPPED:   "STOPPED",
	}
}

func (c *AsynqJobService) Construct(asyncUUid string) {
	path := "./storage/log/job"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Println("Create folder ", path, " :: ", err)
	}
	fileName := fmt.Sprint(path, "/", asyncUUid, ".log")
	cc, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Create ", fileName, " :: ", err)
	}
	c.OS = cc
}

func (c *AsynqJobService) IsPossibleContinue() bool {
	gg := model.AsyncJobView{}
	err := c.DB.Where("asynq_uuid = ?", c.Async_uuid).First(&gg).Error
	if err != nil {
		log.Println("Check asynq_uuid ", c.Async_uuid, " ", err)
		return false
	}
	hh := model.AsyncJobView{}
	err = c.DB.Where("app_uuid = ?", gg.App_uuid).Order("created_at DESC").First(&hh).Error
	if err != nil {
		log.Println("Check app_uuid with same asynq_uuid ", c.Async_uuid, " not found")
		return false
	}
	if hh.Asynq_uuid == gg.Asynq_uuid {
		return true
	}
	c.DB.Model(&model.AsynqJob{}).Where("asynq_uuid = ?", c.Async_uuid).Update("message_text", "Job is stopped").Update("status", c.GetStatus().STATUS_STOPPED)
	return false
}

func (c *AsynqJobService) Add(props AsynqJobAddPayload) (model.AsynqJob, error) {
	asynqjobModel := model.AsynqJob{}
	asynqjobModel.App_uuid = props.App_uuid
	asynqjobModel.Asynq_uuid = props.Asynq_uuid
	asynqjobModel.Status = c.GetStatus().STATUS_PENDING
	asynqjobModel.Payload = datatypes.JSON([]byte(props.Payload))
	asynqjobModel.Message_text = props.Message_text
	asynqjobModel.Table_name = props.Table_name
	time.Sleep(1 * time.Second)
	currentTime := time.Now()
	UnixNano := currentTime.UnixNano()
	asynqjobModel.Order_number = model.NullInt64{
		Valid: true,
		Int64: UnixNano,
	}
	if err := c.DB.Create(&asynqjobModel); err != nil {
		return model.AsynqJob{}, err.Error
	}
	return model.AsynqJob{}, nil
}

func (c *AsynqJobService) UpdateByAsynqUUID(status string, message string) (model.AsyncJobView, error) {
	// Write log
	c.WriteToLog(message)
	// Update the status
	mm := model.AsyncJobView{}
	err := c.DB.Model(model.AsynqJob{}).Where("asynq_uuid = ?", c.Async_uuid).Update("status", status).Error
	if err != nil {
		return model.AsyncJobView{}, err
	}
	return mm, nil
}

func (c *AsynqJobService) WriteToLog(message string) {
	if c.OS == nil {
		log.Println("You have no define yet construct -> call .construct on this AsynqJobService")
		panic(1)
	}
	c.OS.Write([]byte(message))
}

func (c *AsynqJobService) DeleteByUuids(uuids []string) (bool, error) {
	asynqJobDatas := []model.AsyncJobView{}
	err := c.DB.Model(&model.AsynqJob{}).Where("uuid IN ?", uuids).Find(&asynqJobDatas).Error
	if err != nil {
		return false, err
	}
	for _, v := range asynqJobDatas {
		redis_support.RedisClient.Del(context.Background(), fmt.Sprintln("*:", v.Asynq_uuid))
	}
	if err := c.DB.Where("uuid IN ?", uuids).Delete(model.AsynqJob{}); err != nil {
		return false, err.Error
	}
	return true, nil
}

func (c *AsynqJobService) Gets() *gorm.DB {
	return c.DB.Model(model.AsyncJobView{})
}

func (c *AsynqJobService) StopLastProcess() {
	err := c.DB.Model(model.AsynqJob{}).Where("status IN ?", []string{c.GetStatus().STATUS_PROCESS, c.GetStatus().STATUS_PENDING}).Update("status", c.GetStatus().STATUS_STOPPED).Error
	if err != nil {
		log.Fatalln(err.Error())
		panic(1)
	}
}
