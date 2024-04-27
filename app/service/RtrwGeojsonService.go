package service

import (
	"encoding/json"
	model "gis_map_info/app/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RtrwGeojsonService struct {
	DB *gorm.DB
	// Custom own payload
	RtrwGeojsonAdd    rtrwGeojsonAddPayload
	RtrwGeojsonUpdate rtrwGeojsonUpdatePayload
}

type rtrwGeojsonAddPayload struct {
	Order_number  int64
	Uuid          string
	Geojson       string
	Properties    datatypes.JSON
	Rtrw_id       int64
	Rtrw_group_id int64
	Rtrw_file_id  int64
}

type rtrwGeojsonUpdatePayload struct {
}

func (c *RtrwGeojsonService) Add(props rtrwGeojsonAddPayload) (model.RtrwGeojson, error) {
	rtrwGeojsonModel := model.RtrwGeojson{}
	rtrwGeojsonModel.Order_number = props.Order_number
	rtrwGeojsonModel.Uuid = props.Uuid
	rtrwGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	rtrwGeojsonModel.Properties = propertiesByte
	rtrwGeojsonModel.Rtrw_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_file_id,
	}
	rtrwGeojsonModel.Rtrw_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_id,
	}
	rtrwGeojsonModel.Rtrw_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_group_id,
	}
	err := c.DB.Create(&rtrwGeojsonModel).Error
	if err != nil {
		return model.RtrwGeojson{}, err
	}
	return rtrwGeojsonModel, nil
}

func (c *RtrwGeojsonService) Update(props rtrwGeojsonAddPayload) (model.RtrwGeojson, error) {
	rtrwGeojsonModel := model.RtrwGeojson{}
	rtrwGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	rtrwGeojsonModel.Properties = propertiesByte
	rtrwGeojsonModel.Rtrw_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_file_id,
	}
	rtrwGeojsonModel.Rtrw_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_id,
	}
	rtrwGeojsonModel.Rtrw_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_group_id,
	}
	err := c.DB.Where("uuid = ?", props.Uuid).Save(&rtrwGeojsonModel).Error
	if err != nil {
		return model.RtrwGeojson{}, err
	}
	return rtrwGeojsonModel, nil
}

func (c *RtrwGeojsonService) DeleteByRtrwGroupId(rtrw_group_id int64) error {
	err := c.DB.Where("rtrw_group_id = ?", rtrw_group_id).Delete(model.RtrwGeojson{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *RtrwGeojsonService) Gets() (tx *gorm.DB) {
	return c.DB.Model(model.RtrwGeojson{})
}

func (c *RtrwGeojsonService) GetByUUID(uuid string) (model.RtrwGeojson, error) {
	gets := c.Gets()
	rdD := model.RtrwGeojson{}
	if err := gets.Where("uuid = ?", uuid).First(&rdD).Error; err != nil {
		return model.RtrwGeojson{}, err
	}
	return rdD, nil
}
