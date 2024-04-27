package service

import (
	"encoding/json"
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"
	"gis_map_info/support/gorm_support"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RtrwService struct {
	DB *gorm.DB

	// Embed the other struct
	RtrwServiceAddType    rtrwServiceAddType
	RtrwServiceUpdateType rtrwServiceUpdateType
	RtrwGroupAddType      rtrwGroupAddType
	RtrwMbtilePayload     struct {
		Id         *int64
		File_name  string
		Uuid       string
		Rtrw_id    int64
		Created_at string
		Updated_at string
		Checked_at string
	}
}

func (c *RtrwService) Gets() *gorm.DB {
	RtrwModel := Model.RtrwType{}
	rtrwDB := gorm_support.DB.Model(&RtrwModel)
	return rtrwDB
}

func (c *RtrwService) GetById(id int) *gorm.DB {
	rtrwDB := c.Gets()
	rtrwDB.Where("id = ?", id)
	return rtrwDB
}

type rtrwServiceAddType struct {
	Name           string
	RegProvince_id int64
	RegRegency_id  int64
	RegDistrict_id int64
	RegVillage_id  int64
	Place_string   string
	Status         string
}

func (c *RtrwService) Add(props rtrwServiceAddType) (Model.RtrwType, error) {
	rtrwModel := Model.RtrwType{}
	rtrwModel.Name = props.Name
	rtrwModel.RegProvince_id = props.RegProvince_id
	rtrwModel.RegRegency_id = props.RegRegency_id
	rtrwModel.RegDistrict_id = props.RegDistrict_id
	rtrwModel.RegVillage_id = props.RegVillage_id
	rtrwModel.Place_string = props.Place_string
	rtrwModel.Status = props.Status
	err := c.DB.Create(&rtrwModel).Error
	if err != nil {
		return Model.RtrwType{}, err
	}
	return rtrwModel, nil
}

type rtrwServiceUpdateType struct {
	rtrwServiceAddType
	Id int64
}

func (c *RtrwService) Update(props rtrwServiceUpdateType) (Model.RtrwType, error) {

	rtrwModel := Model.RtrwType{}
	rtrwModel.Id = props.Id
	rtrwModel.Name = props.Name
	rtrwModel.RegProvince_id = props.RegProvince_id
	rtrwModel.RegRegency_id = props.RegRegency_id
	rtrwModel.RegDistrict_id = props.RegDistrict_id
	rtrwModel.RegVillage_id = props.RegVillage_id
	rtrwModel.Place_string = props.Place_string
	rtrwModel.Status = props.Status

	err := c.DB.Save(&rtrwModel).Error
	if err != nil {
		return Model.RtrwType{}, err
	}
	return rtrwModel, nil
}

func (c *RtrwService) DeleteByIds(arr []int) error {
	err := c.DB.Where("id IN ?", arr).Delete(&Model.RtrwType{}).Error
	return err
}

func (c *RtrwService) GetGroupsByRtrwId(rtrw_id int) ([]Model.RtrwGroupView, error) {
	rtrwGroups := []Model.RtrwGroupView{}
	err := c.DB.Where("rtrw_id = ?", rtrw_id).Find(&rtrwGroups).Error
	if err != nil {
		return []Model.RtrwGroupView{}, err
	}
	return rtrwGroups, nil
}

type rtrwGroupAddType struct {
	Id         int64
	Uuid       string
	Rtrw_id    int64  `validate:"required"`
	Asset_key  string `validate:"required"`
	Status     string
	Name       string
	Properties map[string]interface{}
}

func (c *RtrwService) AddGroup(props rtrwGroupAddType) (Model.RtrwGroup, error) {
	props.Name = strings.ReplaceAll(props.Asset_key, "_", " ")
	rtrwGroup := Model.RtrwGroup{}
	if props.Id != 0 {
		rtrwGroup.Id = props.Id
	}
	rtrwGroup.Uuid = props.Uuid
	if props.Uuid == "" {
		uuid := uuid.New()
		rtrwGroup.Uuid = uuid.String()
	}
	rtrwGroup.Rtrw_id = props.Rtrw_id
	rtrwGroup.Asset_key = props.Asset_key
	rtrwGroup.Status = props.Status
	rtrwGroup.Name = props.Name
	propertiesByte, _ := json.Marshal(props.Properties)
	rtrwGroup.Properties = propertiesByte
	err := c.DB.Create(&rtrwGroup).Error
	if err != nil {
		return Model.RtrwGroup{}, err
	}
	return rtrwGroup, nil
}

func (c *RtrwService) DeleteGroupByRtrwId(rtrw_id int) error {
	err := c.DB.Where("rtrw_id = ?", rtrw_id).Delete(Model.RtrwGroup{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *RtrwService) GetRtrwGroups() (tx *gorm.DB) {
	return c.DB.Model(&model.RtrwGroup{})
}

func (c *RtrwService) DeleteMbtileByRtrwId(rtrw_id int) error {
	err := c.DB.Model(Model.RtrwMbtile{}).Where("rtrw_id = ?", rtrw_id).Delete(Model.RtrwMbtile{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *RtrwService) DeleteMbtileExceptRtrwMbtileIds_withRtrw_id(rtrw_mbtile_ids []int, rtrw_id int) error {
	err := c.DB.Model(Model.RtrwMbtile{}).Where("id NOT IN ?", rtrw_mbtile_ids).Where("rtrw_id = ?", rtrw_id).Delete(Model.RtrwMbtile{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *RtrwService) AddMbtile() (model.RtrwMbtile, error) {
	rtrwMbtilePayload := c.RtrwMbtilePayload
	rtrwMbtile := model.RtrwMbtile{}
	if rtrwMbtilePayload.Id != nil {
		rtrwMbtile.Id = *rtrwMbtilePayload.Id
		err := c.DB.Model(model.RtrwMbtile{}).Where("id = ?", rtrwMbtile.Id).First(&rtrwMbtile).Error
		if err != nil {
			return model.RtrwMbtile{}, err
		}
	}
	rtrwMbtile.UUID = rtrwMbtilePayload.Uuid
	rtrwMbtile.File_name = rtrwMbtilePayload.File_name
	rtrwMbtile.Rtrw_id = Model.NullInt64{
		Valid: true,
		Int64: rtrwMbtilePayload.Rtrw_id,
	}
	var err error = nil
	err = c.DB.Model(model.RtrwMbtile{}).Where("id = ?", rtrwMbtile.Id).Save(&rtrwMbtile).Error
	if err != nil {
		return model.RtrwMbtile{}, err
	}
	return rtrwMbtile, nil
}
