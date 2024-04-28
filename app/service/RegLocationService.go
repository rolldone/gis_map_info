package service

import (
	"fmt"
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"
	"strconv"

	"gorm.io/gorm"
)

type RegLocationService struct {
	DB *gorm.DB
}

func (c *RegLocationService) GetProvinces() ([]Model.RegProvince, error) {
	regProvinceDatas := []Model.RegProvince{}
	err := c.DB.Order("name ASC").Find(&regProvinceDatas).Error
	if err != nil {
		return []Model.RegProvince{}, err
	}
	return regProvinceDatas, nil
}

func (c *RegLocationService) GetRegenciesByProvinceId(province_id int) ([]Model.RegRegency, error) {
	regRegencyDatas := []Model.RegRegency{}
	err := c.DB.Where("reg_province_id = ?", province_id).Order("name ASC").Find(&regRegencyDatas).Error
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

func (c *RegLocationService) GetNearProvinceByLocation(lat float64, lng float64, limit int, km int) ([]Model.RegProvince, error) {
	regProvinceDatas := []Model.RegProvince{}
	lat_string := fmt.Sprintf("%f", lat)
	lng_string := fmt.Sprintf("%f", lng)
	err := c.DB.Model(&model.RegProvince{}).Where("ST_DWithin(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(" + lng_string + ", " + lat_string + ")::geography, " + strconv.Itoa(km) + ")").Order("ST_Distance(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(" + lng_string + "," + lat_string + ")::geography)").Limit(limit).Find(&regProvinceDatas).Error
	if err != nil {
		return []Model.RegProvince{}, err
	}
	return regProvinceDatas, nil
}
