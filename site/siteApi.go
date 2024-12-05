package site

import (
	"seed-sync/common"
	"seed-sync/seedSyncServer"

	"github.com/gin-gonic/gin"
)

func AddSite(ctx *gin.Context) {
	var request AddSiteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(200, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	//参数校验
	if request.SiteName == "" {
		ctx.JSON(200, common.FailResult("站点名不能为空"))
		return
	}
	//todo: host(http格式校验)
	if request.Host == "" {
		ctx.JSON(200, common.FailResult("host不能为空"))
		return
	}

	if request.MaxPerMin <= 0 || request.MaxPerHour <= 0 || request.MaxPerDay <= 0 {
		ctx.JSON(200, common.FailResult("每分钟、每小时、每天的最大请求数必须大于0"))
		return
	}

	supportedSiteMap, exists := common.CacheGetObject[map[string]seedSyncServer.SupportSiteResponse](common.SUPPORT_SITE_CACHE_KEY)
	if !exists {
		ctx.JSON(200, common.FailResult("获取支持的站点失败"))
		return
	}
	if _, exists := supportedSiteMap[request.SiteName]; !exists {
		ctx.JSON(200, common.FailResult("站点"+request.SiteName+"不支持"))
		return
	}
	request.ShowName = supportedSiteMap[request.SiteName].ShowName

	//添加站点
	if err := SiteService.AddSite(&request); err != nil {
		ctx.JSON(200, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(200, common.SuccessResult("添加站点成功"))
}

func UpdateSite(ctx *gin.Context) {
	var request UpdateSiteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(200, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	//参数校验
	if request.SiteName == "" {
		ctx.JSON(200, common.FailResult("站点名不能为空"))
		return
	}

	// todo: host(http格式校验)
	if request.Host == "" {
		ctx.JSON(200, common.FailResult("host不能为空"))
		return
	}
	if request.MaxPerMin <= 0 || request.MaxPerHour <= 0 || request.MaxPerDay <= 0 {
		ctx.JSON(200, common.FailResult("每分钟、每小时、每天的最大请求数必须大于0"))
		return
	}
	supportedSiteMap, exists := common.CacheGetObject[map[string]seedSyncServer.SupportSiteResponse](common.SUPPORT_SITE_CACHE_KEY)
	if !exists {
		ctx.JSON(200, common.FailResult("获取支持的站点失败"))
		return
	}
	if _, exists := supportedSiteMap[request.SiteName]; !exists {
		ctx.JSON(200, common.FailResult("站点"+request.SiteName+"不支持"))
		return
	}
	request.ShowName = supportedSiteMap[request.SiteName].ShowName

	if err := SiteService.UpdateSite(&request); err != nil {
		ctx.JSON(200, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(200, common.SuccessResult("更新站点成功"))
}

// 批量更新站点顺序
func BatchUpdateSiteOrders(ctx *gin.Context) {
	var request []SiteOrderUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(200, common.FailResult("请求参数错误:"+err.Error()))
		return
	}
	if err := SiteService.BatchUpdateSiteOrders(request); err != nil {
		ctx.JSON(200, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(200, common.SuccessResult("更新站点顺序成功"))
}

// 获取站点列表
func GetSiteList(ctx *gin.Context) {
	siteList, err := SiteService.GetSiteList()
	if err != nil {
		ctx.JSON(200, common.FailResult(err.Error()))
		return
	}
	ctx.JSON(200, common.SuccessResult(siteList))
}
