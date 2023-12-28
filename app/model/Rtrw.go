package model

type RtrwType struct {
	Id          int64 `gorm:"primaryKey" json:"id"`
	Province_id int64 `gorm:"type:uint" json:"province_id"`
	Regency_id  int64 `gorm:"type:uint" json:"regency_id"`
	District_id int64 `gorm:"type:uint" json:"district_id"`
	Village_id  int64 `gorm:"type:uint" json:"village_id"`
}

// Set the table name for the User model
func (RtrwType) TableName() string {
	return "rtrw" // Replace with your existing table name
}
