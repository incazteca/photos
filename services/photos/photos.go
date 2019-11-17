package photos

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"io"
	"os"
)

// Exif Exif data for a photo
type Exif map[string]interface{}

// Photo contain data regarding photo
type Photo struct {
	ID       int
	Location string
	Exif     Exif
}

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

// FetchAll get all photos
func FetchAll(db *sql.DB) ([]*Photo, error) {
	rows, err := db.Query("SELECT id, location, exif FROM photos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	photos := make([]*Photo, 0)
	for rows.Next() {
		photo := new(Photo)
		err := rows.Scan(&photo.ID, &photo.Location, &photo.Exif)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// Fetch retrieve a photo
func Fetch(id int, db *sql.DB) (*Photo, error) {
	row := db.QueryRow("SELECT id, location, exif FROM photos WHERE id = ?", id)

	photo := new(Photo)
	err := row.Scan(&photo.ID, &photo.Location, &photo.Exif)
	if err != nil {
		return nil, err
	}

	return photo, nil
}

// Create Places photo in file system and stores info in database
func Create(db *sql.DB, photo io.ReadCloser) (int, error) {

	fileName := createFileName()

	err := createFile(fileName, photo)
	if err != nil {
		return 0, err
	}

	id := 0
	err = db.QueryRow(
		"INSERT INTO photos (location) VALUES ($1) RETURNING id",
		fileName,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func createFileName() string {
	fileUUID := uuid.Must(uuid.NewV4()).String()

	// Get actual filename
	return "/storage/" + fileUUID + ".jpeg"
}

func createFile(fileName string, photo io.ReadCloser) error {
	homeDir, _ := os.UserHomeDir()
	outFile, err := os.Create(homeDir + fileName)

	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, photo)

	if err != nil {
		return err
	}

	return nil
}
