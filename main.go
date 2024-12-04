package main

import (
	"seed-sync/log"
	"seed-sync/router"

	"go.uber.org/zap"
)

// @title           Seed Sync API
// @version         1.0
// @description     This is a seed sync server.
// @BasePath        /
// @schemes         http
func main() {
	// 初始化日志
	log.InitLogger()
	defer log.Sugar.Sync()

	// 初始化路由
	r := router.InitRouter()

	// 启动服务器
	log.Info("Server starting on :8705")
	if err := r.Run(":8705"); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}
