package model

type RdtrType struct {
	Id             int64  `gorm:"primaryKey" json:"id,omitempty"`
	Name           string `gorm:"type:varchar(255);column:name" json:"name"`
	RegProvince_id int64  `gorm:"type:uint;column:reg_province_id" json:"reg_province_id"`
	RegRegency_id  int64  `gorm:"type:uint;column:reg_regency_id" json:"reg_regency_id"`
	RegDistrict_id int64  `gorm:"type:uint;column:reg_district_id" json:"reg_district_id,omitempty"`
	RegVillage_id  int64  `gorm:"type:uint;column:reg_village_id" json:"reg_village_id,omitempty"`
	Status         string `gorm:"type:uint;column:status" json:"status"`
}

// Set the table name for the User model
func (RdtrType) TableName() string {
	return "rdtr" // Replace with your existing table name
}
