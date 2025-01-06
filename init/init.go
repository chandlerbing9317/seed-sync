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
	_ "seed-sync/site/impl/nexus"
	"seed-sync/user"
	"time"
)

const (
	EXEC_SCHEDULER_TASK_INTERVAL = time.Second * 5
)

//统一管控的初始化，因为当前项目有很多初始化需要做的内容，但这些初始化又有先后顺序，所以在这里统一管控
//还一个目的是便于包的解耦，避免包间循环依赖问题

func Init() {
	//初始化配置
	initConfig()
	//初始化日志
	initLogger()
	// 初始化http client
	initHttpClient()

	//初始化数据库
	initDb()
	//初始化DAO
	initDAO()
	//初始化服务
	initService()
	//初始化注册定时任务
	initRegisterScheduler()
	//初始化启动任务
	initStartTask()
	//启动定时任务
	execSchedulerTask()
	log.Info("服务初始化成功")
}

func initConfig() {
	config.InitConfig()
}

func initLogger() {
	log.InitLogger()
	defer log.Sugar.Sync()
}

func initHttpClient() {
	log.Info("服务启动===>   初始化http client")
	common.InitHttpClient()
	log.Info("初始化http client完成")
}
func initDb() {
	log.Info("服务启动====>   初始化数据库")
	db.InitDb()
	log.Info("初始化数据库完成")
}
func initDAO() {
	log.Info("服务启动=====>   初始化DAO层")
	db.InitSystemParamDAO()
	cookieCloud.InitCookieCloudDAO()
	downloader.InitDownloaderDAO()
	scheduler.InitSchedulerTaskDAO()
	seedSync.InitSeedSyncDAO()
	site.InitSiteDAO()
	user.InitUserDAO()
	log.Info("初始化DAO完成")
}

func initService() {
	log.Info("服务启动======>   初始化Service层")
	user.InitUserService()
	cookieCloud.InitCookieCloudService()
	downloader.InitDownloaderService()
	scheduler.InitSchedulerService()
	seedSync.InitSeedSyncService()
	site.InitSiteService()
	seedSyncServer.InitSeedSyncServerService()
	log.Info("初始化Service层完成")
}
func initRegisterScheduler() {
	log.Info("服务启动=======>   初始化定时任务")
	//定时任务 检查用户状态
	scheduler.SchedulerService.RegisterExecuteFunc(common.CHECK_USER_EXECUTE_CONTENT, user.UserService.CheckUserScheduler)
	//定时任务 同步支持的站点
	scheduler.SchedulerService.RegisterExecuteFunc(common.GET_SITE_EXECUTE_CONTENT, site.SiteService.GetSupportedSiteScheduler)
	//定时任务 通过cookieCloud同步cookie
	scheduler.SchedulerService.RegisterExecuteFunc(common.SYNC_COOKIE_CLOUD_EXECUTE_CONTENT, cookieCloud.CookieCloudService.SyncCookieScheduler)

	//定时任务 辅种
	scheduler.SchedulerService.RegisterExecuteFunc(common.SYNC_SEED_EXECUTE_CONTENT, seedSync.SeedSyncService.SeedSyncScheduler)
	log.Info("初始化定时任务完成")
}

// 一些初始化需要执行的任务
func initStartTask() {
	//1. 校验用户是否可用
	log.Info("服务启动========>   校验用户是否可用")
	user.UserService.CheckUser()
	log.Info("校验用户是否可用完成")
	//2. 同步支持的站点
	log.Info("服务启动========>   同步支持的站点")
	site.SiteService.GetSupportedSite()
	log.Info("同步支持的站点完成")
	//3. 同步站点cookie
	log.Info("服务启动========>   同步站点cookie")
	cookieCloud.CookieCloudService.SyncCookie()
	log.Info("同步站点cookie完成")
}

// 异步线程，轮询启动
func execSchedulerTask() {
	log.Info("服务启动=========>   启动定时任务")
	go func() {
		for {
			scheduler.SchedulerService.ExecuteSchedulerTask()
			time.Sleep(EXEC_SCHEDULER_TASK_INTERVAL)
		}
	}()
	log.Info("启动定时任务完成")
}
