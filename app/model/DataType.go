package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type NullInt struct {
	Int   int
	Valid bool
}

func (ni *NullInt) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int:
		ni.Int = v
		ni.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (ni NullInt) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int, nil
}

func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int)
}

func (ni *NullInt) UnmarshalJSON(b []byte) error {
	var tempInt int
	if err := json.Unmarshal(b, &tempInt); err != nil {
		return err
	}
	ni.Int = tempInt
	ni.Valid = true
	return nil
}

type NullInt16 struct {
	Int   int16
	Valid bool
}

func (ni *NullInt16) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int16:
		ni.Int = v
		ni.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (ni NullInt16) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int, nil
}

func (ni NullInt16) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int)
}

func (ni *NullInt16) UnmarshalJSON(b []byte) error {
	var tempInt int16
	if err := json.Unmarshal(b, &tempInt); err != nil {
		return err
	}
	ni.Int = tempInt
	ni.Valid = true
	return nil
}

type NullInt32 struct {
	Int   int32
	Valid bool
}

func (ni *NullInt32) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int32:
		ni.Int = v
		ni.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (ni NullInt32) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int, nil
}

func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int)
}

func (ni *NullInt32) UnmarshalJSON(b []byte) error {
	var tempInt int32
	if err := json.Unmarshal(b, &tempInt); err != nil {
		return err
	}
	ni.Int = tempInt
	ni.Valid = true
	return nil
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	Int64 int64
	Valid bool
}

// Scan implements the sql.Scanner interface.
func (ni *NullInt64) Scan(value interface{}) error {
	if value == nil {
		ni.Int64, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int64:
		ni.Int64 = v
		ni.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

// Value implements the driver.Valuer interface.
func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	var tempInt64 int64
	if err := json.Unmarshal(b, &tempInt64); err != nil {
		return err
	}
	ni.Int64 = tempInt64
	ni.Valid = true
	return nil
}

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

func (nf *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		nf.Float32, nf.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case float32:
		nf.Float32 = v
		nf.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (nf NullFloat32) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float32, nil
}

func (nf NullFloat32) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float32)
}

func (nf *NullFloat32) UnmarshalJSON(b []byte) error {
	var tempFloat float32
	if err := json.Unmarshal(b, &tempFloat); err != nil {
		return err
	}
	nf.Float32 = tempFloat
	nf.Valid = true
	return nil
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

func (nf *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		nf.Float64, nf.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case float64:
		nf.Float64 = v
		nf.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	var tempFloat float64
	if err := json.Unmarshal(b, &tempFloat); err != nil {
		return err
	}
	nf.Float64 = tempFloat
	nf.Valid = true
	return nil
}

type NullBool struct {
	Bool  bool
	Valid bool
}

// Scan implements the sql.Scanner interface.
func (nb *NullBool) Scan(value interface{}) error {
	if value == nil {
		nb.Bool, nb.Valid = false, false
		return nil
	}

	switch v := value.(type) {
	case bool:
		nb.Bool = v
		nb.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

// Value implements the driver.Valuer interface.
func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	var tempBool bool
	if err := json.Unmarshal(b, &tempBool); err != nil {
		return err
	}
	nb.Bool = tempBool
	nb.Valid = true
	return nil
}

type NullString struct {
	String string
	Valid  bool
}

func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}

	switch v := value.(type) {
	case string:
		ns.String = v
		ns.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	var tempString string
	if err := json.Unmarshal(b, &tempString); err != nil {
		return err
	}
	ns.String = tempString
	ns.Valid = true
	return nil
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
		nt.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan type: %T", value)
	}
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	var tempTime time.Time
	if err := json.Unmarshal(b, &tempTime); err != nil {
		return err
	}
	nt.Time = tempTime
	nt.Valid = true
	return nil
}

func Float64(value float64) float64 {
	return (value)
}
