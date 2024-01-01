package service

import (
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RdtrFileService struct {
	DB *gorm.DB
	// Custom own payload
	RdtrFileAdd random_92923m23unjnvjfv
}

type random_92923m23unjnvjfv struct {
	Uuid string
}

func (c *RdtrFileService) Add(props random_92923m23unjnvjfv) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	rdtrFile.UUID = props.Uuid
	if err := c.DB.Create(&rdtrFile).Error; err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrFileService) DeleteById(id int) bool {
	return true
}

func (c *RdtrFileService) Gets(props interface{}) []Model.RdtrFile {
	return []Model.RdtrFile{}
}

func (c *RdtrFileService) GetByUUID(uuid string) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	err := c.DB.Where("uuid = ?", uuid).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrFileService) GetById(id int) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	err := c.DB.Where("id = ?", id).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}
