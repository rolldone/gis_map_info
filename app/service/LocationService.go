package service

import (
	"gis_map_info/app/model"

	"gorm.io/gorm"
)

type LocationService struct {
	DB *gorm.DB
}

func (c *LocationService) GetsProvince() *gorm.DB {
	return c.DB.Model(model.RegProvince{})
}
