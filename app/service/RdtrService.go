package service

import (
	Helper "gis_map_info/app/helper"
	Model "gis_map_info/app/model"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type RdtrService struct {
	DB *gorm.DB
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

func (c *RdtrService) Add(props interface{}) (Model.RdtrType, error) {

	var result struct {
		Name        string `json:"name"`
		Province_id int64  `json:"province_id"`
		Regency_id  int64  `json:"regency_id"`
		District_id int64  `json:"district_id"`
		Village_id  int64  `json:"village_id"`
	}

	err := Helper.ToStructFromMap(props, &result)
	if err != nil {
		return Model.RdtrType{}, err
	}

	rdtrModel := Model.RdtrType{}
	rdtrModel.Name = result.Name
	rdtrModel.RegProvince_id = result.Province_id
	rdtrModel.RegRegency_id = result.Regency_id
	rdtrModel.RegDistrict_id = result.District_id
	rdtrModel.RegVillage_id = result.Village_id

	err = c.DB.Create(&rdtrModel).Error
	if err != nil {
		return Model.RdtrType{}, err
	}
	return rdtrModel, nil
}

func (c *RdtrService) Update(props interface{}) (Model.RdtrType, error) {

	var result struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Province_id int64  `json:"province_id"`
		Regency_id  int64  `json:"regency_id"`
		District_id int64  `json:"district_id"`
		Village_id  int64  `json:"village_id"`
	}

	err := Helper.ToStructFromMap(props, &result)
	if err != nil {
		return Model.RdtrType{}, err
	}

	rdtrModel := Model.RdtrType{}
	rdtrModel.Id = result.Id
	rdtrModel.Name = result.Name
	rdtrModel.RegProvince_id = result.Province_id
	rdtrModel.RegRegency_id = result.Regency_id
	rdtrModel.RegDistrict_id = result.District_id
	rdtrModel.RegVillage_id = result.Village_id

	err = c.DB.Save(&rdtrModel).Error
	if err != nil {
		return Model.RdtrType{}, err
	}
	return rdtrModel, nil
}

func (c *RdtrService) DeleteByIds(arr []int) error {
	err := Model.DB.Where("id IN ?", arr).Delete(&Model.RdtrType{}).Error
	return err
}

func (c *RdtrService) GetGroupsByRdtrId(rdtr_id int) ([]Model.RdtrGroup, error) {
	rdtrGroups := []Model.RdtrGroup{}
	err := Model.DB.Find(&rdtrGroups).Where("rdtr_id = ?", rdtr_id).Error
	if err != nil {
		return []Model.RdtrGroup{}, err
	}
	return rdtrGroups, nil
}

func (c *RdtrService) AddGroup(props interface{}) (Model.RdtrGroup, error) {
	var propsT struct {
		Id         *big.Int         `json:"id"`
		Rdtr_id    int64            `json:"rdtr_id"`
		Properties pgtype.JSONCodec `json:"properties"`
		Status     string           `json:"status"`
		Name       string           `json:"name"`
		Cat_key    string           `json:"cat_key"`
	}
	Helper.ToStructFromMap(props, &propsT)
	rdtrGroup := Model.RdtrGroup{}
	rdtrGroup.Rdtr_id = propsT.Rdtr_id
	rdtrGroup.Cat_key = propsT.Cat_key
	rdtrGroup.Status = propsT.Status
	rdtrGroup.Name = propsT.Name
	rdtrGroup.Properties = propsT.Properties
	err := c.DB.Create(&rdtrGroup).Error
	if err != nil {
		return Model.RdtrGroup{}, err
	}
	return rdtrGroup, nil
}

func (c *RdtrService) DeleteGroupByRdtrId(rdtr_id int) error {
	err := Model.DB.Delete(Model.RdtrGroup{}).Where("rdtr_id = ?", rdtr_id).Error
	if err != nil {
		return err
	}
	return nil
}
