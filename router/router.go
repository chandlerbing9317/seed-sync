package router

import (
	"seed-sync/cookieCloud"
	"seed-sync/router/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	// 设置 gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中间件
	router.Use(middleware.Cors())
	router.Use(middleware.TraceLogger())
	router.Use(middleware.Recovery())

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	// @Summary      健康检查
	// @Description  服务健康检查接口
	// @Tags         系统
	// @Produce      json
	// @Success      200  {object}  gin.H
	// @Router       /ping [get]
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	//cookie cloud相关api
	router.GET("/cookie-cloud/get", cookieCloud.GetCookieCloudConfig)
	router.POST("/cookie-cloud/create-or-update", cookieCloud.CreateOrUpdateCookieCloud)

	return router
}
