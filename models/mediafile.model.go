package models

import (
	"time"

	"github.com/uptrace/bun"
)

type CreateMediafileDto struct {
	Name   string `bun:"name" json:"name"`
	Path   string `bun:"path" json:"path"`
	Hash   string `bun:"hash" json:"hash"`
	Size   int64  `bun:"size" json:"size"`
	LinkID int    `bun:"link_id" json:"linkId"` // Преобразуем строку в int
}

type Mediafile struct {
	bun.BaseModel `bun:"table:mediafiles"`
	ID            int `bun:"id,pk,autoincrement" json:"id"`
	CreateMediafileDto
	DateAdded time.Time `bun:"date_added,default:current_timestamp" json:"dateAdded"`
}
