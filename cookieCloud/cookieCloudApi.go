package cookieCloud

import (
	"fmt"
	"net/http"
	"seed-sync/common"

	"github.com/gin-gonic/gin"
)

func CreateCookieCloud(c *gin.Context) {
	var config CookieCloudConfig

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	if err := paramCheck(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	err := CookieCloudService.CreateCookieCloud(&config)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("cookie cloud配置成功"))
}

func UpdateCookieCloud(c *gin.Context) {
	var config CookieCloudConfig

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud更新失败"+err.Error()))
		return
	}
	if err := paramCheck(&config); err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud更新失败"+err.Error()))
		return
	}
	err := CookieCloudService.UpdateCookieCloud(&config)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud更新失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("cookie cloud更新成功"))
}

func DeleteCookieCloud(c *gin.Context) {
	err := CookieCloudService.DeleteCookieCloud()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("cookie cloud删除失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("cookie cloud删除成功"))
}

func GetCookieCloudConfig(c *gin.Context) {
	config, err := CookieCloudService.GetCookieCloudConfig()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult("获取cookie cloud配置失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(config))
}

func paramCheck(config *CookieCloudConfig) error {

	//参数校验
	if config.Url == "" || config.UserKey == "" || config.P2pPassword == "" || config.SyncCron == "" {
		return fmt.Errorf("包含未填的必填项")
	}
	//url合法性校验
	if err := common.ValidateURL(config.Url); err != nil {
		return fmt.Errorf("url不合法" + err.Error())
	}
	//规范化url
	url, err := common.NormalizeURL(config.Url)
	if err != nil {
		return fmt.Errorf("url不合法" + err.Error())
	}
	config.Url = url

	//cron表达式合法性校验
	if _, err := common.GetNextExecuteTime(config.SyncCron); err != nil {
		return fmt.Errorf("cron表达式不合法" + err.Error())
	}
	return nil
}
