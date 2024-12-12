package site

import (
	"fmt"
	"net/http"
	"seed-sync/common"
	"seed-sync/seedSyncServer"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_TIMEOUT      = 60
	DEFAULT_MAX_PER_MIN  = 10
	DEFAULT_MAX_PER_HOUR = 50
	DEFAULT_MAX_PER_DAY  = 200
	DEFAULT_IS_ACTIVE    = true
	DEFAULT_PROXY        = false
	DEFAULT_USER_AGENT   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

func AddSite(ctx *gin.Context) {
	var request AddSiteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	//参数校验
	if err := paramCheck(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult("添加站点失败:"+err.Error()))
		return
	}
	//添加站点
	if err := SiteService.AddSite(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.SuccessResult("添加站点成功"))
}

func UpdateSite(ctx *gin.Context) {
	var request UpdateSiteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	//参数校验
	if err := paramCheck(request.AddSiteRequest); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult("更新站点失败:"+err.Error()))
		return
	}
	//更新站点
	if err := SiteService.UpdateSite(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.SuccessResult("更新站点成功"))
}

// 删除站点
func DeleteSite(ctx *gin.Context) {
	siteName := ctx.Param("siteName")
	if err := SiteService.DeleteSite(siteName); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.SuccessResult("删除站点成功"))
}

func GetAvailableSites(ctx *gin.Context) {
	supportedSiteMap, exists := common.CacheGetObject[map[string]seedSyncServer.SupportSiteResponse](common.SUPPORT_SITE_CACHE_KEY)
	if !exists {
		ctx.JSON(http.StatusOK, common.FailResult("获取支持的站点失败"))
		return
	}
	siteList := make([]seedSyncServer.SupportSiteResponse, 0, len(supportedSiteMap))
	for _, site := range supportedSiteMap {
		siteList = append(siteList, site)
	}
	ctx.JSON(http.StatusOK, common.SuccessResult(siteList))
}

// 批量更新站点顺序
func BatchUpdateSiteOrders(ctx *gin.Context) {
	var request []SiteOrderUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	if err := SiteService.BatchUpdateSiteOrders(request); err != nil {
		ctx.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.SuccessResult("更新站点顺序成功"))
}

// 获取站点列表
func GetSiteList(ctx *gin.Context) {
	siteList, err := SiteService.GetSiteList()
	if err != nil {
		ctx.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.SuccessResult(siteList))
}

func paramCheck(request *AddSiteRequest) error {
	if request.SiteName == "" {
		return fmt.Errorf("站点名不能为空")
	}
	//这里的host改为删除开头的http://或https:// 强制改为https
	request.Host = strings.TrimPrefix(request.Host, common.Http)
	request.Host = strings.TrimPrefix(request.Host, common.Https)

	request.Host = common.Https + request.Host
	if err := common.ValidateURL(request.Host); err != nil {
		return fmt.Errorf("host格式不正确: %s", err.Error())
	}
	host, err := common.NormalizeURL(request.Host)
	if err != nil {
		return fmt.Errorf("host格式不正确: %s", err.Error())
	}
	request.Host = host

	if request.Timeout < 0 {
		return fmt.Errorf("超时时间不能小于0")
	}
	if request.MaxPerMin < 0 || request.MaxPerHour < 0 || request.MaxPerDay < 0 {
		return fmt.Errorf("每分钟、每小时、每天的最大请求数必须大于0")
	}
	if request.Timeout == 0 {
		request.Timeout = DEFAULT_TIMEOUT
	}
	if request.MaxPerMin == 0 {
		request.MaxPerMin = DEFAULT_MAX_PER_MIN
	}
	if request.MaxPerHour == 0 {
		request.MaxPerHour = DEFAULT_MAX_PER_HOUR
	}
	if request.MaxPerDay == 0 {
		request.MaxPerDay = DEFAULT_MAX_PER_DAY
	}
	if request.UserAgent == "" {
		request.UserAgent = DEFAULT_USER_AGENT
	}

	supportedSiteMap, exists := common.CacheGetObject[map[string]seedSyncServer.SupportSiteResponse](common.SUPPORT_SITE_CACHE_KEY)
	if !exists {
		return fmt.Errorf("获取支持的站点失败")
	}
	if _, exists := supportedSiteMap[request.SiteName]; !exists {
		return fmt.Errorf("站点%s不支持", request.SiteName)
	}
	request.ShowName = supportedSiteMap[request.SiteName].ShowName
	return nil
}
