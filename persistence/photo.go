package persistence

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/incazteca/photos/models"
)

type PhotoPersistence struct {
	db *sql.DB
}

// Exif Exif data for a photo
type Exif map[string]interface{}

func NewPhotoPersistence(db *sql.DB) PhotoPersistence {
	return PhotoPersistence{db}
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

func (p *PhotoPersistence) FetchAll() ([]*models.Photo, error) {
	rows, err := p.db.Query("SELECT id, location, exif FROM photos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	photos := make([]*models.Photo, 0)
	for rows.Next() {
		photo := new(models.Photo)
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
func (p *PhotoPersistence) Fetch(id int) (*models.Photo, error) {
	row := p.db.QueryRow("SELECT id, location, exif FROM photos WHERE id = ?", id)

	photo := new(models.Photo)
	err := row.Scan(&photo.ID, &photo.Location, &photo.Exif)
	if err != nil {
		return nil, err
	}

	return photo, nil
}

// Save Places photo in file system and stores info in database
func (p *PhotoPersistence) Save(photo models.Photo) (int, error) {
	id := 0
	err := p.db.QueryRow(
		"INSERT INTO photos (location) VALUES ($1) RETURNING id",
		photo.Location,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
