package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Photo contain data regarding photo
type Photo struct {
	ID       int
	Location string
	Exif     map[string]interface{}
}

// Exif Exif data for a photo
type Exif map[string]interface{}

// Value a function that satisfies the driver.Value interface in database/sql
func (exif Exif) Value() (driver.Value, error) {
	j, err := json.Marshal(exif)
	return j, err
}

// Scan a function that satisfies the driver.Scan interface in database/sql
func (exif *Exif) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion []byte failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*exif, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion map[string]interface{} failed")
	}

	return nil
}
