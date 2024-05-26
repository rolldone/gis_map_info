package model

import (
	"time"

	"gorm.io/datatypes"
)

type ZlpGroup struct {
	Id         int64          `gorm:"primaryKey" json:"id,omitempty"`
	Uuid       string         `gorm:"column:uuid" json:"uuid,omitempty"`
	Zlp_id     int64          `gorm:"column:zlp_id" json:"zlp_id"`
	Properties datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	Status     string         `gorm:"column:status;type:varchar" json:"status"`
	Name       string         `gorm:"column:name;type:varchar(255)" json:"name"`
	Asset_key  string         `gorm:"column:asset_key;type:varchar(255)" json:"asset_key"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Datas      []ZlpFile      `gorm:"foreignKey:zlp_group_id" json:"datas"`
	Mbtiles    []ZlpMbtile    `gorm:"foreignKey:zlp_group_id;References:id" json:"mbtiles"`

	// Relation
	Zlp *ZlpType `gorm:"<-:false;foreignKey:zlp_id;References:id" json:"zlp,omitempty"`
}

type ZlpGroupView struct {
	ZlpGroup
	Unvalidated NullInt64 `gorm:"column:unvalidated" json:"unvalidated,omitempty"`
	Validated   NullInt64 `gorm:"column:validated" json:"validated,omitempty"`
}

type ZlpGroupDistinctAssetView struct {
	Name      string      `gorm:"column:name;type:varchar(255)" json:"name"`
	Asset_key string      `gorm:"column:asset_key;type:varchar(255)" json:"asset_key"`
	Mbtiles   []ZlpMbtile `gorm:"foreignKey:asset_key;References:asset_key" json:"mbtiles"`
}

// Set the table name for the User model
func (ZlpGroup) TableName() string {
	return "zlp_group" // Replace with your existing table name
}
