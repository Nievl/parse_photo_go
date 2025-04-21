package domains

import (
	"context"
	"fmt"
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

func (m *MedifilesDbService) Create(mediafile models.CreateMediafileDto, linkId int) error {
	ctx := context.Background()

	tx, _ := m.db.BeginTx(ctx, nil)
	defer tx.Commit()
	_, err := tx.Exec("INSERT INTO mediafiles (name, path, hash, size) VALUES (?, ?, ?, ?)",
		mediafile.Name,
		mediafile.Path,
		mediafile.Hash,
		mediafile.Size)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert mediafile: %s", err.Error())
	}
	_, err = tx.Exec(`
		INSERT INTO links_mediafiles (link_id, mediafile_id) 
		VALUES (?, 
			(SELECT id as mediafile_id
			FROM mediafiles
			WHERE path = ?)
		)`, linkId, mediafile.Path)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert link_mediafile: %s", err.Error())
	}

	return nil
}

func (m *MedifilesDbService) Remove(id int) error {
	return nil
}

func (m *MedifilesDbService) GetAllByLinkId(linkId int) ([]models.Mediafile, error) {
	ctx := context.Background()
	result := []models.Mediafile{}
	query := m.db.NewSelect().Model(&result).
		TableExpr("mediafiles AS m").
		Join("JOIN links_medifiles AS lm ON lm.medifile_id = m.id").
		Where("lm.link_id = ?", linkId)

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mediafiles: %s", err.Error())
	}
	return result, nil
}
