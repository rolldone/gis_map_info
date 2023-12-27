package model

type RtrwType struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	UUID        string `gorm:"unique;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Province_id int64  `gorm:"type:uint" json:"province_id"`
	Regency_id  int64  `gorm:"type:uint" json:"regency_id"`
	District_id int64  `gorm:"type:uint" json:"district_id"`
	Village_id  int64  `gorm:"type:uint" json:"village_id"`
}
