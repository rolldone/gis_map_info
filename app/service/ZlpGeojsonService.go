package service

import (
	"encoding/json"
	model "gis_map_info/app/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ZlpGeojsonService struct {
	DB *gorm.DB
	// Custom own payload
	ZlpGeojsonAdd    zlpGeojsonAddPayload
	ZlpGeojsonUpdate zlpGeojsonUpdatePayload
}

type zlpGeojsonAddPayload struct {
	Order_number int64
	Uuid         string
	Geojson      string
	Properties   datatypes.JSON
	Zlp_id       int64
	Zlp_group_id int64
	Zlp_file_id  int64
}

type zlpGeojsonUpdatePayload struct {
}

func (c *ZlpGeojsonService) Add(props zlpGeojsonAddPayload) (model.ZlpGeojson, error) {
	zlpGeojsonModel := model.ZlpGeojson{}
	zlpGeojsonModel.Order_number = props.Order_number
	zlpGeojsonModel.Uuid = props.Uuid
	zlpGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	zlpGeojsonModel.Properties = propertiesByte
	zlpGeojsonModel.Zlp_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_file_id,
	}
	zlpGeojsonModel.Zlp_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_id,
	}
	zlpGeojsonModel.Zlp_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_group_id,
	}
	err := c.DB.Create(&zlpGeojsonModel).Error
	if err != nil {
		return model.ZlpGeojson{}, err
	}
	return zlpGeojsonModel, nil
}

func (c *ZlpGeojsonService) Update(props zlpGeojsonAddPayload) (model.ZlpGeojson, error) {
	zlpGeojsonModel := model.ZlpGeojson{}
	zlpGeojsonModel.Geojson = props.Geojson
	propertiesByte, _ := json.Marshal(props.Properties)
	zlpGeojsonModel.Properties = propertiesByte
	zlpGeojsonModel.Zlp_file_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_file_id,
	}
	zlpGeojsonModel.Zlp_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_id,
	}
	zlpGeojsonModel.Zlp_group_id = model.NullInt64{
		Valid: true,
		Int64: props.Zlp_group_id,
	}
	err := c.DB.Where("uuid = ?", props.Uuid).Save(&zlpGeojsonModel).Error
	if err != nil {
		return model.ZlpGeojson{}, err
	}
	return zlpGeojsonModel, nil
}

func (c *ZlpGeojsonService) DeleteByZlpGroupId(zlp_group_id int64) error {
	err := c.DB.Where("zlp_group_id = ?", zlp_group_id).Delete(model.ZlpGeojson{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *ZlpGeojsonService) Gets() (tx *gorm.DB) {
	return c.DB.Model(model.ZlpGeojson{})
}

func (c *ZlpGeojsonService) GetByUUID(uuid string) (model.ZlpGeojson, error) {
	gets := c.Gets()
	rdD := model.ZlpGeojson{}
	if err := gets.Where("uuid = ?", uuid).First(&rdD).Error; err != nil {
		return model.ZlpGeojson{}, err
	}
	return rdD, nil
}
