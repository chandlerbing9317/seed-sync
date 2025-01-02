package init

import (
	"seed-sync/common"
	"seed-sync/config"
	"seed-sync/cookieCloud"
	"seed-sync/db"
	"seed-sync/downloader"
	"seed-sync/log"
	"seed-sync/scheduler"
	"seed-sync/seedSync"
	"seed-sync/seedSyncServer"
	"seed-sync/site"
)

//统一管控的初始化，因为当前项目有很多初始化需要做的内容，但这些初始化又有先后顺序，所以在这里统一管控
//还一个目的是便于包的解耦，避免包间循环依赖问题

func Init() {
	//初始化日志
	initLogger()
	//初始化数据库
	initDb()
	//初始化配置
	initConfig()
	//初始化下载器
	initDownloader()
	//初始化站点客户端
	initSiteClient()
	//初始化cookie cloud
	initCookieCloud()
	//初始化注册定时任务
	initRegisterScheduler()
}

func initLogger() {
	log.InitLogger()
	defer log.Sugar.Sync()
}

func initDb() {
	log.Info("服务启动===>   初始化数据库")
	db.InitDb()
	log.Info("初始化数据库完成")
}

func initConfig() {
	log.Info("服务启动===>   初始化配置")
	config.InitConfig()
	log.Info("初始化配置完成")
}

func initDownloader() {
	log.Info("服务启动===>   初始化下载器")
	downloader.InitDownloader()
	log.Info("初始化下载器完成")
}

func initCookieCloud() {
	log.Info("服务启动===>   初始化cookie cloud")
	cookieCloud.InitCookieCloud()
	log.Info("初始化cookie cloud完成")
}

func initSiteClient() {
	log.Info("服务启动===>   初始化站点客户端")
	site.SiteService.InitSiteClient()
	log.Info("初始化站点客户端完成")
}

func initRegisterScheduler() {
	log.Info("服务启动===>   初始化定时任务")
	//定时任务 检查用户状态
	scheduler.SchedulerService.RegisterExecuteFunc(common.CHECK_USER_EXECUTE_CONTENT, seedSyncServer.SeedSyncServerService.CheckUserScheduler)
	//定时任务 通过cookieCloud同步cookie
	scheduler.SchedulerService.RegisterExecuteFunc(common.SYNC_COOKIE_CLOUD_EXECUTE_CONTENT, cookieCloud.CookieCloudService.SyncCookieScheduler)
	//定时任务 同步支持的站点
	scheduler.SchedulerService.RegisterExecuteFunc(common.GET_SITE_EXECUTE_CONTENT, seedSyncServer.SeedSyncServerService.GetSupportedSiteScheduler)
	//定时任务 辅种
	scheduler.SchedulerService.RegisterExecuteFunc(common.SYNC_SEED_EXECUTE_CONTENT, seedSync.SeedSyncService.SeedSyncScheduler)
	log.Info("初始化定时任务完成")
}
