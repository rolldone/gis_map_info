package service

import (
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RdtrMbtileService struct {
	DB         *gorm.DB
	AddPayload *struct {
		Uuid      string
		File_name string
	}
	UpdatePayload *struct {
		Id            int64
		Rdtr_group_id int64
		Rdtr_id       int64
	}
}

func (c *RdtrMbtileService) Add() (model.RdtrMbtile, error) {
	props := c.AddPayload
	rdtrFile := Model.RdtrMbtile{}
	rdtrFile.UUID = props.Uuid
	rdtrFile.File_name = props.File_name

	if err := c.DB.Create(&rdtrFile).Error; err != nil {
		return Model.RdtrMbtile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrMbtileService) Update() (Model.RdtrMbtile, error) {
	props := c.UpdatePayload
	rdtrFile := Model.RdtrMbtile{}
	c.DB.Where("id = ?", props.Id).First(&rdtrFile)
	rdtrFile.Rdtr_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_id,
	}
	err := c.DB.Save(&rdtrFile).Error
	if err != nil {
		return Model.RdtrMbtile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrMbtileService) DeleteById(id int) bool {
	return true
}

func (c *RdtrMbtileService) Gets() *gorm.DB {
	rdtrFiles := Model.RdtrMbtile{}
	return c.DB.Model(&rdtrFiles)
}

func (c *RdtrMbtileService) GetByUUID(uuid string) (Model.RdtrMbtile, error) {
	rdtrFile := Model.RdtrMbtile{}
	err := c.DB.Where("uuid = ?", uuid).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrMbtile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrMbtileService) GetById(id int) (Model.RdtrMbtile, error) {
	rdtrFile := Model.RdtrMbtile{}
	err := c.DB.Where("id = ?", id).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrMbtile{}, err
	}
	return rdtrFile, nil
}
