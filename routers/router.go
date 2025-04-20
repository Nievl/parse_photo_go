package routers

import (
	"parse_photo_go/controllers"
	"parse_photo_go/domains"
	"parse_photo_go/services"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func InitRoutes(app *gin.Engine, db *bun.DB) {
	linkDbService := domains.NewLinksDbService(db)
	mediafilesDbService := domains.NewMediafilesDbService(db)
	mediafilesService := services.NewMediafilesService(*mediafilesDbService)
	linkService := services.NewLinkService(*linkDbService, *mediafilesService)
	linkController := controllers.NewLinkController(*linkService)

	app.POST("/links", linkController.Create)
	app.GET("/links", linkController.GetAll)
	app.DELETE("/links", linkController.Remove)
	app.GET("/links/download", linkController.DownloadFiles)
	app.GET("/links/scan_files_for_link", linkController.ScanFilesForLink)
	app.GET("/links/check_downloaded", linkController.CheckDownloaded)
	app.GET("/links/tag_unreachable", linkController.TagUnreachable)

}
