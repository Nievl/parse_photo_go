package main

import (
	"log"
	"net/http"
	"os"
	"parse_photo_go/domains"
	"parse_photo_go/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("server starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	db, err := domains.CheckAndCreateTables()
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) { c.File("./static/index.html") }) // Отдаём index.html при запросе на "/"
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "ok"}) })
	routers.InitRoutes(r, db)
	r.Run(":" + port)

}
