package domains

import (
	"context"
	"errors"
	"fmt"
	"parse_photo_go/models"
	"parse_photo_go/utils"
	"time"

	"github.com/uptrace/bun"
	"modernc.org/sqlite"
)

type LinksDbService struct {
	db *bun.DB
}

func NewLinksDbService(db *bun.DB) *LinksDbService {
	return &LinksDbService{
		db: db,
	}
}

func (s *LinksDbService) CreateLink(path string, filename string) error {
	ctx := context.Background()

	link := models.Link{
		Path:         path,
		Name:         filename,
		IsDownloaded: false,
		DateCreate:   time.Now(),
		DateUpdate:   time.Now(),
	}
	query := s.db.NewInsert().Model(&link)

	_, err := query.Exec(ctx)

	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == 2067 { // 2067 = SQLITE_CONSTRAINT_UNIQUE
				fmt.Println("Link already exists:", path)
				return fmt.Errorf("link already exists")
			}
		}
		// Обработка других ошибок
		return fmt.Errorf("unexpected db error: %w", err)
	}

	return nil
}

func (s *LinksDbService) GetAll(isReachable, showDuplicate bool) ([]models.LinkWithDuplicatePath, error) {
	ctx := context.Background()

	var links []models.LinkWithDuplicatePath

	query := s.db.NewSelect().
		TableExpr("links AS l").
		ColumnExpr("l.*, d.path AS duplicate_path").
		Join("LEFT JOIN links AS d ON l.duplicate_id = d.id").
		Where("l.is_reachable = ?", isReachable).
		OrderExpr("l.date_update DESC, l.is_downloaded ASC")

	if showDuplicate {
		query = query.Where("l.duplicate_id IS NOT NULL")
	} else {
		query = query.Where("l.duplicate_id IS NULL")
	}

	err := query.Scan(ctx, &links)
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return nil, err
	}

	return links, nil
}

func (s *LinksDbService) Remove(id int64) error {
	ctx := context.Background()

	query := s.db.NewDelete().Model(&models.Link{}).Where("id= ?", id)
	_, err := query.Exec(ctx)
	if err != nil {
		fmt.Println("Error deleting link:", err)
	}
	return nil
}

func (s *LinksDbService) TagUnreachable(id int64, reachable bool) error {
	ctx := context.Background()
	query := s.db.NewUpdate().Model(&models.Link{}).
		Set("is_reachable = ?", reachable).
		Set("date_update = ?", utils.DateConvert()).
		Where("id = ?", id)
	_, err := query.Exec(ctx)
	if err != nil {
		fmt.Println("Error updating link:", err)
		return fmt.Errorf("link not found")
	}

	return nil
}

func (s *LinksDbService) GetOne(id int64) (*models.Link, error) {
	ctx := context.Background()
	var link models.Link
	err := s.db.NewSelect().Model(&link).Where("id=?", id).Scan(ctx)
	if err != nil {
		fmt.Println("Error fetching link:", err)
		return nil, fmt.Errorf("link not found")
	}

	return &link, nil
}

func (s *LinksDbService) UpdateFilesNumber(id int64, link models.UpdateLinkDto) error {
	ctx := context.Background()
	query := s.db.NewUpdate().Table("links").
		Set("is_downloaded = ?", link.IsDownloaded).
		Set("progress = ?", link.Progress).
		Set("mediafiles = ?", link.Mediafiles).
		Set("downloaded_mediafiles = ?", link.DownloadedMediafiles).
		Set("date_update = ?", utils.DateConvert()).
		Where("id = ?", id)
	_, err := query.Exec(ctx)
	if err != nil {
		fmt.Println("Error updating link:", err)
		return fmt.Errorf("link not found")
	}

	return nil
}
