package controllers

import (
	"parse_photo_go/models"
	"parse_photo_go/services"
	"parse_photo_go/utils"
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
		ctx.JSON(400, utils.ResultMaker("Invalid request body"+err.Error()))
		return
	}

	err := c.linkService.Create(link)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Failed to create link: \n"+err.Error()))
		return
	} else {
		ctx.JSON(200, utils.ResultMaker("Link created successfully"))
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
		ctx.JSON(500, utils.ResultMaker("Failed to get links"))
		return
	}
	ctx.JSON(200, links)

}

func (c *LinkController) Remove(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.Remove(id)
}

func (c *LinkController) DownloadFiles(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Invalid id parameter"))
		return
	}
	err = c.linkService.DownloadFiles(id)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Failed to download files: \n"+err.Error()))
		return
	} else {
		ctx.JSON(200, utils.ResultMaker("Files downloaded successfully"))
		return
	}
}

func (c *LinkController) ScanFilesForLink(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Invalid id parameter"))
		return
	}
	result, err := c.linkService.ScanFilesForLink(id)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Failed to scan files: \n"+err.Error()))
		return
	} else {
		ctx.JSON(200, utils.ResultMaker(result))
		return
	}
}

func (c *LinkController) CheckDownloaded(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Invalid id parameter"))
		return
	}
	result, err := c.linkService.CheckDownloaded(id)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Failed to check downloaded status: \n"+err.Error()))
	} else {
		ctx.JSON(200, utils.ResultMaker(result))
	}

}

func (c *LinkController) TagUnreachable(ctx *gin.Context) {
	raw_id := ctx.Query("id")
	id, err := strconv.ParseInt(raw_id, 10, 64)
	isReachable := ctx.DefaultQuery("isReachable", "true")
	reachable, _ := strconv.ParseBool(isReachable)
	if err != nil {
		ctx.JSON(400, utils.ResultMaker("Invalid id parameter"))
		return
	}
	c.linkService.TagUnreachable(id, reachable)
}
