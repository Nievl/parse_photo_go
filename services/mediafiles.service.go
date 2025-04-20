package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"parse_photo_go/models"
)

type MediafilesService struct{}

func NewMediafilesService() *MediafilesService {
	return &MediafilesService{}
}

func (s *MediafilesService) Create(mediafile models.CreateMediafileDto) error {
	// implementation for creating a mediafile
	return nil
}

func (s *MediafilesService) Remove(id int64) error {
	// implementation for removing a mediafile by id
	return nil
}

func (s *MediafilesService) GetAllByLinkId(linkId int) ([]models.Mediafile, error) {
	// implementation for getting all mediafiles by link id
	// For now, just returning an empty slice and nil error
	return []models.Mediafile{}, nil
}

func (s *MediafilesService) DownloadFile(url string, filePath string, linkId int) (models.CreateMediafileDto, error) {

	// implementation for downloading a file from url to filePath
	// and associating it with the linkId
	return models.CreateMediafileDto{}, nil
}

func Download(url string, pathName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	err = os.WriteFile(pathName, body, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
