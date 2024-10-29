package api

import (
	"net/http"
	"seed-sync/common"
	"seed-sync/service"

	"github.com/gin-gonic/gin"
)

const (
	CookieCloudConfigKey = "cookie_cloud_config_key"
)

// @Summary      添加或更新cookie cloud配置
// @Description  添加或更新cookie cloud配置
// @Tags         cookie cloud
// @Accept       json
// @Produce      json
// @Param        config  body  service.CookieCloudConfig  true  "cookie cloud配置"
// @Success      200  
// @Router       /cookie-cloud/add-or-update [post]
func AddOrUpdateCookieCloudConfig(c *gin.Context) {
	var config service.CookieCloudConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	err := service.AddOrUpdateCookieCloudConfig(&config)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("cookie cloud配置成功"))
}

// @Summary      获取cookie cloud配置
// @Description  获取cookie cloud配置
// @Tags         cookie cloud
// @Produce      json
// @Success      200 
// @Router       /cookie-cloud/get [get]
func GetCookieCloudConfig(c *gin.Context) {
	config, err := service.GetCookieCloudConfig()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("获取cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(config))
}

// @Summary      同步站点cookie
// @Description  同步站点cookie
// @Tags         cookie cloud
// @Produce      json
// @Success      200 
// @Router       /cookie-cloud/sync-site-cookie [get]
func SyncSiteCookie(c *gin.Context) {
	cookie, err := service.SyncSiteCookie()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("同步站点cookie失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(cookie))
}
