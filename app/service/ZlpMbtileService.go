package service

import (
	"fmt"
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type ZlpMbtileService struct {
	DB         *gorm.DB
	AddPayload *struct {
		Uuid            string
		File_name       string
		Asset_key       string
		Reg_province_id int64
		Zlp_group_id    int64
	}
	UpdatePayload *struct {
		Id              int64
		Zlp_id          int64
		Asset_key       string
		Reg_province_id int64
		Zlp_group_id    int64
	}
}

func (c *ZlpMbtileService) Add() (model.ZlpMbtile, error) {
	props := c.AddPayload
	zlpFile := Model.ZlpMbtile{}
	zlpFile.UUID = props.Uuid
	zlpFile.File_name = props.File_name
	zlpFile.Asset_key = props.Asset_key
	zlpFile.Reg_province_id = props.Reg_province_id
	zlpFile.Zlp_group_id = props.Zlp_group_id
	fmt.Println("asset_key", props.Asset_key)
	if err := c.DB.Create(&zlpFile).Error; err != nil {
		return Model.ZlpMbtile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpMbtileService) Update() (Model.ZlpMbtile, error) {
	props := c.UpdatePayload
	zlpFile := Model.ZlpMbtile{}
	c.DB.Where("id = ?", props.Id).First(&zlpFile)
	zlpFile.Zlp_id = Model.NullInt64{
		Valid: true,
		Int64: props.Zlp_id,
	}
	zlpFile.Asset_key = props.Asset_key
	zlpFile.Reg_province_id = props.Reg_province_id
	zlpFile.Zlp_group_id = props.Zlp_group_id
	err := c.DB.Save(&zlpFile).Error
	if err != nil {
		return Model.ZlpMbtile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpMbtileService) DeleteById(id int) bool {
	return true
}

func (c *ZlpMbtileService) Gets() *gorm.DB {
	zlpFiles := Model.ZlpMbtile{}
	return c.DB.Model(&zlpFiles)
}

func (c *ZlpMbtileService) GetByUUID(uuid string) (Model.ZlpMbtile, error) {
	zlpFile := Model.ZlpMbtile{}
	err := c.DB.Where("uuid = ?", uuid).First(&zlpFile).Error
	if err != nil {
		return Model.ZlpMbtile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpMbtileService) GetById(id int) (Model.ZlpMbtile, error) {
	zlpFile := Model.ZlpMbtile{}
	err := c.DB.Where("id = ?", id).First(&zlpFile).Error
	if err != nil {
		return Model.ZlpMbtile{}, err
	}
	return zlpFile, nil
}
