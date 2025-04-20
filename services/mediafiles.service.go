package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"parse_photo_go/domains"
	"parse_photo_go/models"
	"parse_photo_go/utils"
)

type MediafilesService struct{ mediafilesDbService domains.MedifilesDbService }

func NewMediafilesService(mediafilesDbService domains.MedifilesDbService) *MediafilesService {
	return &MediafilesService{
		mediafilesDbService: mediafilesDbService,
	}
}

func (s *MediafilesService) Create(mediafile models.CreateMediafileDto) error {
	return s.mediafilesDbService.Create(mediafile)
}

func (s *MediafilesService) Remove(id int) error {
	return s.mediafilesDbService.Remove(id)
}

func (s *MediafilesService) GetAllByLinkId(linkId int) ([]models.Mediafile, error) {
	return s.mediafilesDbService.GetAllByLinkId(linkId)
}

func (s *MediafilesService) DownloadFile(url string, filePath string, linkId int) (models.CreateMediafileDto, error) {

	err := Download(url, filePath)
	if err != nil {
		return models.CreateMediafileDto{}, fmt.Errorf("error downloading file: %v", err)
	}
	hash, err := utils.GetHashByPath(filePath)
	if err != nil {
		return models.CreateMediafileDto{}, fmt.Errorf("error getting hash: %v", err)
	}
	info, _ := os.Stat(filePath)
	mediafile := models.CreateMediafileDto{
		Hash:   hash,
		Name:   "",
		Path:   filePath,
		Size:   info.Size(),
		LinkID: linkId,
	}

	return mediafile, nil
}

func Download(url string, pathName string) error {
	fmt.Printf("Downloading file from %s to %s\n", url, pathName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching URL: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
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
