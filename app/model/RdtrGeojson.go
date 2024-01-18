package model

import (
	"time"

	"gorm.io/datatypes"
)

type RdtrGeojson struct {
	Uuid          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Order_number  int64          `gorm:"column:order_number" json:"order_number,omitempty"`
	Rdtr_id       NullInt64      `gorm:"column:rdtr_id" json:"rdtr_id,omitempty"`
	Rdtr_group_id NullInt64      `gorm:"column:rdtr_group_id" json:"rdtr_group_id,omitempty"`
	Rdtr_file_id  NullInt64      `gorm:"column:rdtr_file_id" json:"rdtr_file_id,omitempty"`
	Geojson       string         `gorm:"type:geometry(GEOMETRY,4326)" json:"geojson"`
	Properties    datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type RdtrGeojsonView struct {
	Uuid          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Order_number  int64          `gorm:"column:order_number" json:"order_number,omitempty"`
	Rdtr_id       NullInt64      `gorm:"column:rdtr_id" json:"rdtr_id,omitempty"`
	Rdtr_group_id NullInt64      `gorm:"column:rdtr_group_id" json:"rdtr_group_id,omitempty"`
	Rdtr_file_id  NullInt64      `gorm:"column:rdtr_file_id" json:"rdtr_file_id,omitempty"`
	Geojson       datatypes.JSON `gorm:"type:geometry(GEOMETRY,4326)" json:"geojson"`
	Properties    datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// type GeoJSONType string

// // Implement the Scanner interface for GeoJSONType
// func (g *GeoJSONType) Scan(value interface{}) error {
// 	// Check if the value is a []byte
// 	*g = value.(string)
// 	return nil
// 	// if data, ok := value.([]byte); ok {
// 	// 	// Set GeoJSONType to the string representation
// 	// 	*g = GeoJSONType(data)
// 	// 	return nil
// 	// }
// 	// return errors.New("failed to scan GeoJSONType")
// }

// Set the table name for the User model
func (c *RdtrGeojson) TableName() string {
	return "rdtr_geojson" // Replace with your existing table name
}
