package service

import (
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RtrwMbtileService struct {
	DB         *gorm.DB
	AddPayload *struct {
		Uuid      string
		File_name string
	}
	UpdatePayload *struct {
		Id            int64
		Rtrw_group_id int64
		Rtrw_id       int64
	}
}

func (c *RtrwMbtileService) Add() (model.RtrwMbtile, error) {
	props := c.AddPayload
	rtrwFile := Model.RtrwMbtile{}
	rtrwFile.UUID = props.Uuid
	rtrwFile.File_name = props.File_name

	if err := c.DB.Create(&rtrwFile).Error; err != nil {
		return Model.RtrwMbtile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwMbtileService) Update() (Model.RtrwMbtile, error) {
	props := c.UpdatePayload
	rtrwFile := Model.RtrwMbtile{}
	c.DB.Where("id = ?", props.Id).First(&rtrwFile)
	rtrwFile.Rtrw_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_id,
	}
	err := c.DB.Save(&rtrwFile).Error
	if err != nil {
		return Model.RtrwMbtile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwMbtileService) DeleteById(id int) bool {
	return true
}

func (c *RtrwMbtileService) Gets() *gorm.DB {
	rtrwFiles := Model.RtrwMbtile{}
	return c.DB.Model(&rtrwFiles)
}

func (c *RtrwMbtileService) GetByUUID(uuid string) (Model.RtrwMbtile, error) {
	rtrwFile := Model.RtrwMbtile{}
	err := c.DB.Where("uuid = ?", uuid).First(&rtrwFile).Error
	if err != nil {
		return Model.RtrwMbtile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwMbtileService) GetById(id int) (Model.RtrwMbtile, error) {
	rtrwFile := Model.RtrwMbtile{}
	err := c.DB.Where("id = ?", id).First(&rtrwFile).Error
	if err != nil {
		return Model.RtrwMbtile{}, err
	}
	return rtrwFile, nil
}
