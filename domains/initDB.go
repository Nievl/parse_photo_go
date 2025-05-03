package domains

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"parse_photo_go/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func CheckAndCreateTables(dbPath string) (*bun.DB, error) {

	log.Println("Reading DB file, DB_NAME:", dbPath)

	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath)
	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	log.Println("Connected to database")
	ctx := context.Background()

	_, err = db.NewCreateTable().Model((*models.Link)(nil)).Exec(ctx)
	if err != nil {
		log.Println("Error creating Link table:", err)
	}
	_, err = db.NewCreateTable().Model((*models.MediafilesLinks)(nil)).Exec(ctx)
	if err != nil {
		log.Println("Error creating MediafilesLinks table:", err)
	}
	_, err = db.NewCreateTable().Model((*models.Mediafile)(nil)).Exec(ctx)
	if err != nil {
		log.Println("Error creating Mediafile table:", err)
	}

	log.Println("All tables checked/created successfully.")
	return db, nil
}
