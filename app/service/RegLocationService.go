package service

import (
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RegLocationService struct {
	DB *gorm.DB
}

func (c *RegLocationService) GetProvinces() ([]Model.RegProvince, error) {
	regProvinceDatas := []Model.RegProvince{}
	err := c.DB.Find(&regProvinceDatas).Error
	if err != nil {
		return []Model.RegProvince{}, err
	}
	return regProvinceDatas, nil
}

func (c *RegLocationService) GetRegenciesByProvinceId(province_id int) ([]Model.RegRegency, error) {
	regRegencyDatas := []Model.RegRegency{}
	err := c.DB.Where("reg_province_id = ?", province_id).Find(&regRegencyDatas).Error
	if err != nil {
		return []Model.RegRegency{}, err
	}
	return regRegencyDatas, nil
}

func (c *RegLocationService) GetDistrictsByRegencyId(regency_id int) ([]Model.RegDistrict, error) {
	regDistrictDatas := []Model.RegDistrict{}
	err := c.DB.Where("reg_regency_id = ?", regency_id).Find(&regDistrictDatas).Error
	if err != nil {
		return []Model.RegDistrict{}, err
	}
	return regDistrictDatas, nil
}

func (c *RegLocationService) GetVillagesByDistrictId(district_id int) ([]Model.RegVillage, error) {
	regVillageDatas := []Model.RegVillage{}
	err := c.DB.Where("reg_district_id = ?", district_id).Find(&regVillageDatas).Error
	if err != nil {
		return []Model.RegVillage{}, err
	}
	return regVillageDatas, nil
}
