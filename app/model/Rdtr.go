package model

import "time"

type RdtrType struct {
	Id             int64       `gorm:"primaryKey" json:"id,omitempty"`
	Name           string      `gorm:"type:varchar(255);column:name" json:"name"`
	RegProvince_id int64       `gorm:"type:uint;column:reg_province_id" json:"reg_province_id"`
	RegRegency_id  int64       `gorm:"type:uint;column:reg_regency_id" json:"reg_regency_id"`
	RegDistrict_id int64       `gorm:"type:uint;column:reg_district_id" json:"reg_district_id,omitempty"`
	RegVillage_id  int64       `gorm:"type:uint;column:reg_village_id" json:"reg_village_id,omitempty"`
	Status         string      `gorm:"type:uint;column:status" json:"status"`
	Place_string   string      `gorm:"type:uint;column:place_string" json:"place_string"`
	Rdtr_groups    []RdtrGroup `gorm:"foreignKey:rdtr_id" json:"rdtr_groups"`
	CreatedAt      time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Set the table name for the User model
func (RdtrType) TableName() string {
	return "rdtr" // Replace with your existing table name
}
