package persistence

import (
	"database/sql"
	"github.com/incazteca/photos/models"
)

type PhotoPersistence struct {
	db *sql.DB
}

func NewPhotoPersistence(db *sql.DB) PhotoPersistence {
	return PhotoPersistence{db}
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
