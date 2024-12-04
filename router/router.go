package router

import (
	"seed-sync/cookieCloud"
	"seed-sync/downloader"
	"seed-sync/router/middleware"

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
	router.POST("/cookie-cloud/create-or-update", cookieCloud.CreateOrUpdateCookieCloud)

	//下载器相关api
	router.POST("/downloader/create", downloader.CreateDownloader)
	router.GET("/downloader/list", downloader.GetDownloaderList)

	return router
}
