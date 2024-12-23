package model

import "time"

type ZlpType struct {
	Id             int64          `gorm:"primaryKey" json:"id,omitempty"`
	Name           string         `gorm:"type:varchar(255);column:name" json:"name"`
	RegProvince_id int64          `gorm:"type:uint;column:reg_province_id" json:"reg_province_id"`
	RegRegency_id  int64          `gorm:"type:uint;column:reg_regency_id" json:"reg_regency_id"`
	RegDistrict_id int64          `gorm:"type:uint;column:reg_district_id" json:"reg_district_id,omitempty"`
	RegVillage_id  int64          `gorm:"type:uint;column:reg_village_id" json:"reg_village_id,omitempty"`
	Status         string         `gorm:"type:uint;column:status" json:"status"`
	Place_string   string         `gorm:"type:uint;column:place_string" json:"place_string"`
	Zlp_groups     []ZlpGroupView `gorm:"foreignKey:zlp_id" json:"zlp_groups"`
	Zlp_mbtiles    []ZlpMbtile    `gorm:"foreignKey:zlp_id" json:"zlp_mbtiles"`
	Reg_province   RegProvince    `gorm:"foreignKey:RegProvince_id" json:"reg_province,omitempty"`
	Reg_regency    RegRegency     `gorm:"foreignKey:RegRegency_id" json:"reg_regency,omitempty"`
	Reg_district   RegDistrict    `gorm:"foreignKey:RegDistrict_id" json:"reg_district,omitempty"`
	CreatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (ZlpType) TableName() string {
	return "zlp" // Replace with your existing table name
}
