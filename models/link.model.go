package models

import (
	"time"

	"github.com/uptrace/bun"
)

// CreateLinkDto — DTO для создания новой записи
type CreateLinkDto struct {
	Path string `json:"path"`
}

// UpdateLinkDto — DTO для обновления записи
type UpdateLinkDto struct {
	IsDownloaded         bool `json:"isDownloaded"`
	Progress             int  `json:"progress"`
	Mediafiles           int  `json:"mediafiles"`
	DownloadedMediafiles int  `json:"downloadedMediafiles"`
}

// Link — объединённая структура с полной информацией
type Link struct {
	bun.BaseModel        `bun:"table:links"`
	ID                   int       `bun:"id,pk,autoincrement" json:"id"`
	Path                 string    `bun:"path,unique" json:"path"`
	Name                 string    `bun:"name" json:"name"`
	IsDownloaded         bool      `bun:"is_downloaded,default:1" json:"isDownloaded"`
	Progress             int       `bun:"progress,default:0" json:"progress"`
	DownloadedMediafiles int       `bun:"downloaded_mediafiles,default:0" json:"downloadedMediafiles"`
	Mediafiles           int       `bun:"mediafiles,default:0" json:"mediafiles"`
	DateUpdate           time.Time `bun:"date_update,default:current_timestamp" json:"dateUpdate"`
	DateCreate           time.Time `bun:"date_create,default:current_timestamp" json:"dateCreate"`
	IsReachable          bool      `bun:"is_reachable,default:0" json:"isReachable"`
	DuplicateID          *int      `bun:"duplicate_id,nullzero" json:"duplicateId,omitempty"`
}

type MediafilesLinks struct {
	bun.BaseModel `bun:"table:mediafiles_links"`
	LinkID        int `bun:"link_id"`
	MediafileID   int `bun:"mediafile_id"`
}

type LinkWithDuplicatePath struct {
	Link          `bun:"embed"`
	DuplicatePath *string `bun:"duplicate_path" json:"duplicatePath"` // указываем d.path прямо здесь!
}
