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

type ZlpMbtilePayload struct {
	Id              *int64
	File_name       string
	Uuid            string
	Zlp_id          int64
	Asset_key       string
	Zlp_group_id    int64
	Reg_province_id int64
	Created_at      string
	Updated_at      string
	Checked_at      string
}

type ZlpService struct {
	DB *gorm.DB

	// Embed the other struct
	ZlpServiceAddType    zlpServiceAddType
	ZlpServiceUpdateType zlpServiceUpdateType
	ZlpGroupAddType      zlpGroupAddType
	ZlpMbtilePayload
}

func (c *ZlpService) Gets() *gorm.DB {
	ZlpModel := Model.ZlpType{}
	zlpDB := gorm_support.DB.Model(&ZlpModel)
	return zlpDB
}

func (c *ZlpService) GetById(id int) *gorm.DB {
	zlpDB := c.Gets()
	zlpDB.Where("id = ?", id)
	return zlpDB
}

type zlpServiceAddType struct {
	Name           string
	RegProvince_id int64
	RegRegency_id  int64
	RegDistrict_id int64
	RegVillage_id  int64
	Place_string   string
	Status         string
}

func (c *ZlpService) Add(props zlpServiceAddType) (Model.ZlpType, error) {
	zlpModel := Model.ZlpType{}
	zlpModel.Name = props.Name
	zlpModel.RegProvince_id = props.RegProvince_id
	zlpModel.RegRegency_id = props.RegRegency_id
	zlpModel.RegDistrict_id = props.RegDistrict_id
	zlpModel.RegVillage_id = props.RegVillage_id
	zlpModel.Place_string = props.Place_string
	zlpModel.Status = props.Status
	err := c.DB.Create(&zlpModel).Error
	if err != nil {
		return Model.ZlpType{}, err
	}
	return zlpModel, nil
}

type zlpServiceUpdateType struct {
	zlpServiceAddType
	Id int64
}

func (c *ZlpService) Update(props zlpServiceUpdateType) (Model.ZlpType, error) {

	zlpModel := Model.ZlpType{}
	zlpModel.Id = props.Id
	zlpModel.Name = props.Name
	zlpModel.RegProvince_id = props.RegProvince_id
	zlpModel.RegRegency_id = props.RegRegency_id
	zlpModel.RegDistrict_id = props.RegDistrict_id
	zlpModel.RegVillage_id = props.RegVillage_id
	zlpModel.Place_string = props.Place_string
	zlpModel.Status = props.Status

	err := c.DB.Save(&zlpModel).Error
	if err != nil {
		return Model.ZlpType{}, err
	}
	return zlpModel, nil
}

func (c *ZlpService) DeleteByIds(arr []int) error {
	err := c.DB.Where("id IN ?", arr).Delete(&Model.ZlpType{}).Error
	return err
}

func (c *ZlpService) GetGroupsByZlpId(zlp_id int) ([]Model.ZlpGroupView, error) {
	zlpGroups := []Model.ZlpGroupView{}
	err := c.DB.Where("zlp_id = ?", zlp_id).Find(&zlpGroups).Error
	if err != nil {
		return []Model.ZlpGroupView{}, err
	}
	return zlpGroups, nil
}

type zlpGroupAddType struct {
	Id         int64
	Uuid       string
	Zlp_id     int64  `validate:"required"`
	Asset_key  string `validate:"required"`
	Status     string
	Name       string
	Properties map[string]interface{}
}

func (c *ZlpService) AddGroup(props zlpGroupAddType) (Model.ZlpGroup, error) {
	props.Name = strings.ReplaceAll(props.Asset_key, "_", " ")
	zlpGroup := Model.ZlpGroup{}
	if props.Id != 0 {
		zlpGroup.Id = props.Id
	}
	zlpGroup.Uuid = props.Uuid
	if props.Uuid == "" {
		uuid := uuid.New()
		zlpGroup.Uuid = uuid.String()
	}
	zlpGroup.Zlp_id = props.Zlp_id
	zlpGroup.Asset_key = props.Asset_key
	zlpGroup.Status = props.Status
	zlpGroup.Name = props.Name
	propertiesByte, _ := json.Marshal(props.Properties)
	zlpGroup.Properties = propertiesByte
	err := c.DB.Create(&zlpGroup).Error
	if err != nil {
		return Model.ZlpGroup{}, err
	}
	return zlpGroup, nil
}

func (c *ZlpService) DeleteGroupByZlpId(zlp_id int) error {
	err := c.DB.Where("zlp_id = ?", zlp_id).Delete(Model.ZlpGroup{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *ZlpService) GetZlpGroups() (tx *gorm.DB) {
	return c.DB.Model(&model.ZlpGroup{})
}

func (c *ZlpService) DeleteMbtileByZlpId(zlp_id int) error {
	err := c.DB.Model(Model.ZlpMbtile{}).Where("zlp_id = ?", zlp_id).Delete(Model.ZlpMbtile{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *ZlpService) DeleteMbtileExceptZlpMbtileIds_withZlp_id(zlp_mbtile_ids []int, zlp_id int) error {
	err := c.DB.Model(Model.ZlpMbtile{}).Where("id NOT IN ?", zlp_mbtile_ids).Where("zlp_id = ?", zlp_id).Delete(Model.ZlpMbtile{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *ZlpService) AddMbtile() (model.ZlpMbtile, error) {
	zlpMbtilePayload := c.ZlpMbtilePayload
	zlpMbtile := model.ZlpMbtile{}
	if zlpMbtilePayload.Id != nil {
		zlpMbtile.Id = *zlpMbtilePayload.Id
		err := c.DB.Model(model.ZlpMbtile{}).Where("id = ?", zlpMbtile.Id).First(&zlpMbtile).Error
		if err != nil {
			return model.ZlpMbtile{}, err
		}
	}
	zlpMbtile.UUID = zlpMbtilePayload.Uuid
	zlpMbtile.File_name = zlpMbtilePayload.File_name
	zlpMbtile.Asset_key = zlpMbtilePayload.Asset_key
	zlpMbtile.Zlp_group_id = zlpMbtilePayload.Zlp_group_id       // default 0 by golang
	zlpMbtile.Reg_province_id = zlpMbtilePayload.Reg_province_id // default 0 by golang
	zlpMbtile.Zlp_id = Model.NullInt64{
		Valid: true,
		Int64: zlpMbtilePayload.Zlp_id,
	}
	var err error = nil
	err = c.DB.Model(model.ZlpMbtile{}).Where("id = ?", zlpMbtile.Id).Save(&zlpMbtile).Error
	if err != nil {
		return model.ZlpMbtile{}, err
	}
	return zlpMbtile, nil
}
