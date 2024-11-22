package cookieCloud

import (
	"net/http"
	"seed-sync/common"

	"github.com/gin-gonic/gin"
)


// @Summary      添加或更新cookie cloud配置
// @Description  添加或更新cookie cloud配置
// @Tags         cookie cloud
// @Accept       json
// @Produce      json
// @Param        config  body  CookieCloudConfig  true  "cookie cloud配置"
// @Success      200
// @Router       /cookie-cloud/create-or-update [post]
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

// @Summary      获取cookie cloud配置
// @Description  获取cookie cloud配置
// @Tags         cookie cloud
// @Produce      json
// @Success      200
// @Router       /cookie-cloud/get [get]
func GetCookieCloudConfig(c *gin.Context) {
	config, err := CookieCloudService.GetCookieCloudConfig()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("获取cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(config))
}
