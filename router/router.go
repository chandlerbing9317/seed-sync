package router

import (
	"seed-sync/cookieCloud"
	"seed-sync/downloader"
	"seed-sync/router/middleware"
	"seed-sync/site"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// 设置 gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(middleware.Cors())
	router.Use(middleware.TraceLogger())
	router.Use(middleware.Recovery())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	//cookie cloud相关api
	router.GET("/cookie-cloud/get", cookieCloud.GetCookieCloudConfig)
	router.POST("/cookie-cloud/create", cookieCloud.CreateCookieCloud)
	router.POST("/cookie-cloud/update", cookieCloud.UpdateCookieCloud)

	//下载器相关api
	router.POST("/downloader/create", downloader.CreateDownloader)
	router.POST("/downloader/update", downloader.UpdateDownloader)
	router.POST("/downloader/delete", downloader.DeleteDownloader)
	router.GET("/downloader/list", downloader.GetDownloaderList)

	//站点相关api
	router.POST("/site/create", site.AddSite)
	router.POST("/site/update", site.UpdateSite)
	router.POST("/site/delete/", site.DeleteSite)
	router.POST("/site/batch-update-orders", site.BatchUpdateSiteOrders)
	router.GET("/site/list", site.GetSiteList)

	return router
}
