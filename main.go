package main

import (
	"seed-sync/log"
	"seed-sync/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	log.InitLogger()
	defer log.Sugar.Sync()

	// 设置 gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	router := gin.New()

	// 使用自定义的 logger 中间件
	router.Use(middleware.TraceLogger())
	router.Use(middleware.Recovery())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 启动服务器
	log.Info("Server starting on :8705")
	if err := router.Run(":8705"); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}

}
