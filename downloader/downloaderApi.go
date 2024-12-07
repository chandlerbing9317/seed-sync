package downloader

import (
	"fmt"
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
	if err := paramCheck(&request); err != nil {
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

// 删除
func DeleteDownloader(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusOK, common.FailResult("删除下载器失败: 名称不能为空"))
		return
	}
	err := DownloaderService.DeleteDownloader(name)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("删除下载器失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("删除下载器成功"))
}

func UpdateDownloader(c *gin.Context) {
	var request DownloaderUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, common.FailResult("更新下载器失败"+err.Error()))
		return
	}
	if err := paramCheck(&request.DownloaderCreateRequest); err != nil {
		c.JSON(http.StatusOK, common.FailResult("更新下载器失败"+err.Error()))
		return
	}
	err := DownloaderService.UpdateDownloader(&request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("更新下载器失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("更新下载器成功"))
}

func GetDownloaderList(c *gin.Context) {
	downloaders := DownloaderService.GetDownloaderList()
	c.JSON(http.StatusOK, common.SuccessResult(downloaders))
}

func paramCheck(request *DownloaderCreateRequest) error {
	if request.Name == "" {
		return fmt.Errorf("名称不能为空")
	}
	if err := common.ValidateURL(request.Url); err != nil {
		return fmt.Errorf("url格式不正确: %s", err.Error())
	}
	url, err := common.NormalizeURL(request.Url)
	if err != nil {
		return fmt.Errorf("url格式不正确: %s", err.Error())
	}
	request.Url = url
	if request.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if request.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if request.Type != common.DOWNLOADER_TYPE_TRANSMISSION && request.Type != common.DOWNLOADER_TYPE_QBITTORRENT {
		return fmt.Errorf("下载器类型不正确")
	}
	return nil
}
