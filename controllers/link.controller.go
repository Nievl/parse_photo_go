package controllers

import (
	"parse_photo_go/helpers"
	"parse_photo_go/models"
	"parse_photo_go/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	linkService services.LinkService
}

func NewLinkController(linkService services.LinkService) *LinkController {
	return &LinkController{linkService: linkService}
}

func (c *LinkController) Create(ctx *gin.Context) {
	var link models.CreateLinkDto
	if err := ctx.ShouldBindJSON(&link); err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid request body"+err.Error()))
		return
	}

	err := c.linkService.Create(link)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Failed to create link"+err.Error()))
		return
	} else {
		ctx.JSON(200, helpers.ResultMaker("Link created successfully"))
		return
	}

}

func (c *LinkController) GetAll(ctx *gin.Context) {
	isReachable := ctx.DefaultQuery("isReachable", "true")
	reachable, _ := strconv.ParseBool(isReachable)
	showDuplicate := ctx.DefaultQuery("showDuplicate", "false")
	duplicate, _ := strconv.ParseBool(showDuplicate)

	links, err := c.linkService.GetAll(reachable, duplicate)

	if err != nil {
		ctx.JSON(500, helpers.ResultMaker("Failed to get links"))
		return
	}
	ctx.JSON(200, links)

}

func (c *LinkController) Remove(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.Remove(id)
}

func (c *LinkController) DownloadFiles(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.DownloadFiles(id)
}

func (c *LinkController) ScanFilesForLink(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.ScanFilesForLink(id)
}

func (c *LinkController) CheckDownloaded(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.CheckDownloaded(id)
}

func (c *LinkController) TagUnreachable(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	isReachable := ctx.DefaultQuery("isReachable", "true")
	reachable, _ := strconv.ParseBool(isReachable)
	if err != nil {
		ctx.JSON(400, helpers.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.TagUnreachable(id, reachable)
}
