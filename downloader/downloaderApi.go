package downloader

import (
	"net/http"
	"seed-sync/common"

	"github.com/gin-gonic/gin"
)

func CreateDownloader(c *gin.Context) {
	var request DownloaderCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, common.FailResult("创建下载器失败"+err.Error()))
		return
	}
	err := DownloaderService.CreateDownloader(&request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("创建下载器失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("创建下载器成功"))
}

func GetDownloaderList(c *gin.Context) {
	downloaders := DownloaderService.GetDownloaderList()
	c.JSON(http.StatusOK, common.SuccessResult(downloaders))
}

