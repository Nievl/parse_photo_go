package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"parse_photo_go/domains"
	"parse_photo_go/models"
	"parse_photo_go/routers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("server starting...")

	cfg := models.MustLoadConfig()
	log := Loger(cfg.Env)

	db, err := domains.CheckAndCreateTables(cfg.StoragePath)
	if err != nil {
		log.Error("Error creating tables: %v", err)
	}

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) { c.File("./static/index.html") }) // Отдаём index.html при запросе на "/"
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "ok"}) })
	routers.InitRoutes(r, db)
	r.Run(cfg.Address)

}

func Loger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
