package model

import (
	"time"

	"gorm.io/datatypes"
)

type ZlpGeojson struct {
	Uuid         string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Order_number int64          `gorm:"column:order_number" json:"order_number,omitempty"`
	Zlp_id       NullInt64      `gorm:"column:zlp_id" json:"zlp_id,omitempty"`
	Zlp_group_id NullInt64      `gorm:"column:zlp_group_id" json:"zlp_group_id,omitempty"`
	Zlp_file_id  NullInt64      `gorm:"column:zlp_file_id" json:"zlp_file_id,omitempty"`
	Geojson      string         `gorm:"type:geometry(GEOMETRY,4326)" json:"geojson"`
	Properties   datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ZlpGeojsonView struct {
	Uuid         string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Order_number int64          `gorm:"column:order_number" json:"order_number,omitempty"`
	Zlp_id       NullInt64      `gorm:"column:zlp_id" json:"zlp_id,omitempty"`
	Zlp_group_id NullInt64      `gorm:"column:zlp_group_id" json:"zlp_group_id,omitempty"`
	Zlp_file_id  NullInt64      `gorm:"column:zlp_file_id" json:"zlp_file_id,omitempty"`
	Geojson      datatypes.JSON `gorm:"type:geometry(GEOMETRY,4326)" json:"geojson"`
	Properties   datatypes.JSON `gorm:"column:properties;type:json" json:"properties"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
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
func (c *ZlpGeojson) TableName() string {
	return "zlp_geojson" // Replace with your existing table name
}
