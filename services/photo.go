package services

import (
	"github.com/gofrs/uuid"
	"github.com/incazteca/photos/models"
	"github.com/incazteca/photos/persistence"
	"io"
	"os"
)

// PhotoService struct with info about persistence
type PhotoService struct {
	photoPersist persistence.PhotoPersistence
}

// NewPhotoService returns struct of photo service
func NewPhotoService(photoPersist persistence.PhotoPersistence) PhotoService {
	return PhotoService{photoPersist}
}

func (ps *PhotoService) FetchAll() ([]*models.Photo, error) {
	return ps.FetchAll()
}

// StorePhoto stores photo in database
func (ps *PhotoService) StorePhoto(rawPhoto io.ReadCloser) (int, error) {
	fileName := createFileName()

	err := createFile(fileName, rawPhoto)
	if err != nil {
		return 0, err
	}

	photo := models.Photo{Location: fileName}

	return ps.photoPersist.Save(photo)
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
