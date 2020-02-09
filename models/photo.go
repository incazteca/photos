package models

// Photo contain data regarding photo
type Photo struct {
	ID       int
	Location string
	Exif     map[string]interface{}
}
