package service

import (
	"encoding/json"
	Model "gis_map_info/app/model"
	"strings"

	"gorm.io/gorm"
)

type RdtrService struct {
	DB *gorm.DB

	// Embed the other struct
	RdtrServiceAddType    rdtrServiceAddType
	RdtrServiceUpdateType rdtrServiceUpdateType
	RdtrGroupAddType      rdtrGroupAddType
}

func (c *RdtrService) Gets() *gorm.DB {
	RdtrModel := Model.RdtrType{}
	rdtrDB := Model.DB.Model(&RdtrModel)
	return rdtrDB
}

func (c *RdtrService) GetById(id int) *gorm.DB {
	rdtrDB := c.Gets()
	rdtrDB.Where("id = ?", id)
	return rdtrDB
}

type rdtrServiceAddType struct {
	Name           string
	RegProvince_id int64
	RegRegency_id  int64
	RegDistrict_id int64
	RegVillage_id  int64
	Place_string   string
	Status         string
}

func (c *RdtrService) Add(props rdtrServiceAddType) (Model.RdtrType, error) {
	rdtrModel := Model.RdtrType{}
	rdtrModel.Name = props.Name
	rdtrModel.RegProvince_id = props.RegProvince_id
	rdtrModel.RegRegency_id = props.RegRegency_id
	rdtrModel.RegDistrict_id = props.RegDistrict_id
	rdtrModel.RegVillage_id = props.RegVillage_id
	rdtrModel.Status = props.Status
	err := c.DB.Create(&rdtrModel).Error
	if err != nil {
		return Model.RdtrType{}, err
	}
	return rdtrModel, nil
}

type rdtrServiceUpdateType struct {
	rdtrServiceAddType
	Id int64
}

func (c *RdtrService) Update(props rdtrServiceUpdateType) (Model.RdtrType, error) {

	rdtrModel := Model.RdtrType{}
	rdtrModel.Id = props.Id
	rdtrModel.Name = props.Name
	rdtrModel.RegProvince_id = props.RegProvince_id
	rdtrModel.RegRegency_id = props.RegRegency_id
	rdtrModel.RegDistrict_id = props.RegDistrict_id
	rdtrModel.RegVillage_id = props.RegVillage_id
	rdtrModel.Place_string = props.Place_string
	rdtrModel.Status = props.Status

	err := c.DB.Save(&rdtrModel).Error
	if err != nil {
		return Model.RdtrType{}, err
	}
	return rdtrModel, nil
}

func (c *RdtrService) DeleteByIds(arr []int) error {
	err := c.DB.Where("id IN ?", arr).Delete(&Model.RdtrType{}).Error
	return err
}

func (c *RdtrService) GetGroupsByRdtrId(rdtr_id int) ([]Model.RdtrGroup, error) {
	rdtrGroups := []Model.RdtrGroup{}
	err := c.DB.Where("rdtr_id = ?", rdtr_id).Find(&rdtrGroups).Error
	if err != nil {
		return []Model.RdtrGroup{}, err
	}
	return rdtrGroups, nil
}

type rdtrGroupAddType struct {
	Id         int64
	Rdtr_id    int64  `validate:"required"`
	Asset_key  string `validate:"required"`
	Status     string
	Name       string
	Properties map[string]interface{}
}

func (c *RdtrService) AddGroup(props rdtrGroupAddType) (Model.RdtrGroup, error) {
	props.Name = strings.ReplaceAll(props.Asset_key, "_", " ")

	rdtrGroup := Model.RdtrGroup{}
	rdtrGroup.Rdtr_id = props.Rdtr_id
	rdtrGroup.Asset_key = props.Asset_key
	rdtrGroup.Status = props.Status
	rdtrGroup.Name = props.Name
	propertiesByte, _ := json.Marshal(props.Properties)
	rdtrGroup.Properties = propertiesByte
	err := c.DB.Create(&rdtrGroup).Error
	if err != nil {
		return Model.RdtrGroup{}, err
	}
	return rdtrGroup, nil
}

func (c *RdtrService) DeleteGroupByRdtrId(rdtr_id int) error {
	err := c.DB.Where("rdtr_id = ?", rdtr_id).Delete(Model.RdtrGroup{}).Error
	if err != nil {
		return err
	}
	return nil
}
