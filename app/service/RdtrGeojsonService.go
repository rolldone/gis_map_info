package service

import (
	"encoding/json"
	model "gis_map_info/app/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RdtrGeojsonService struct {
	DB *gorm.DB
	// Custom own payload
	RdtrGeojsonAdd    rdtrGeojsonAddPayload
	RdtrGeojsonUpdate rdtrGeojsonUpdatePayload
}

type rdtrGeojsonAddPayload struct {
	Uuid          string
	Geojson       string
	Properties    datatypes.JSON
	Rdtr_id       int64
	Rdtr_group_id int64
	Rdtr_file_id  int64
}

type rdtrGeojsonUpdatePayload struct {
}

func (c *RdtrGeojsonService) Add(props rdtrGeojsonAddPayload) (model.RdtrGeojson, error) {
	rdtrGeojsonModel := model.RdtrGeojson{}
	rdtrGeojsonModel.Uuid = props.Uuid
	rdtrGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	rdtrGeojsonModel.Properties = propertiesByte
	rdtrGeojsonModel.Rdtr_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_file_id,
	}
	rdtrGeojsonModel.Rdtr_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_id,
	}
	rdtrGeojsonModel.Rdtr_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_group_id,
	}
	err := c.DB.Create(&rdtrGeojsonModel).Error
	if err != nil {
		return model.RdtrGeojson{}, err
	}
	return rdtrGeojsonModel, nil
}

func (c *RdtrGeojsonService) Update(props rdtrGeojsonAddPayload) (model.RdtrGeojson, error) {
	rdtrGeojsonModel := model.RdtrGeojson{}
	rdtrGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	rdtrGeojsonModel.Properties = propertiesByte
	rdtrGeojsonModel.Rdtr_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_file_id,
	}
	rdtrGeojsonModel.Rdtr_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_id,
	}
	rdtrGeojsonModel.Rdtr_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_group_id,
	}
	err := c.DB.Where("uuid = ?", props.Uuid).Save(&rdtrGeojsonModel).Error
	if err != nil {
		return model.RdtrGeojson{}, err
	}
	return rdtrGeojsonModel, nil
}

func (c *RdtrGeojsonService) DeleteByRdtrGroupId(rdtr_group_id int64) error {
	err := c.DB.Where("rdtr_group_id = ?", rdtr_group_id).Delete(model.RdtrGeojson{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *RdtrGeojsonService) Gets() (tx *gorm.DB) {
	return c.DB.Model(model.RdtrGeojson{})
}

func (c *RdtrGeojsonService) GetByUUID(uuid string) (model.RdtrGeojson, error) {
	gets := c.Gets()
	rdD := model.RdtrGeojson{}
	if err := gets.Where("uuid = ?", uuid).First(&rdD).Error; err != nil {
		return model.RdtrGeojson{}, err
	}
	return rdD, nil
}
