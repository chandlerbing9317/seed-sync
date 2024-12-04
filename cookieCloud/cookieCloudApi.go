package cookieCloud

import (
	"net/http"
	"seed-sync/common"

	"github.com/gin-gonic/gin"
)

func CreateOrUpdateCookieCloud(c *gin.Context) {
	var config CookieCloudConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	err := CookieCloudService.CreateOrUpdateCookieCloud(&config)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("cookie cloud配置成功"))
}

func GetCookieCloudConfig(c *gin.Context) {
	config, err := CookieCloudService.GetCookieCloudConfig()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("获取cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(config))
}
