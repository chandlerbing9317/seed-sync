package main

import (
	initPkg "seed-sync/init"
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
	// 初始化
	initPkg.Init()

	// 初始化路由
	r := router.InitRouter()

	// 启动服务器
	log.Info("服务启动===>   启动服务器")
	if err := r.Run(":8705"); err != nil {
		log.Fatal("服务启动===>   启动服务器失败", zap.Error(err))
	}
}
