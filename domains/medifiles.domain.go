package domains

import (
	"parse_photo_go/models"

	"github.com/uptrace/bun"
)

type MedifilesDbService struct {
	db *bun.DB
}

func NewMediafilesDbService(db *bun.DB) *MedifilesDbService {
	return &MedifilesDbService{
		db: db,
	}
}

func (m *MedifilesDbService) Create(mediafile models.CreateMediafileDto) error {
	return nil
}

func (m *MedifilesDbService) Remove(id int) error {
	return nil
}

func (m *MedifilesDbService) GetAllByLinkId(linkId int) ([]models.Mediafile, error) {
	return nil, nil
}
