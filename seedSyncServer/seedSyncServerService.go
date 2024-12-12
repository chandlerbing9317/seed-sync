package seedSyncServer

import (
	"math"
	"seed-sync/common"
)

type seedSyncServerService struct {
	syncSeedServerDriver *ServerClient
}

var SeedSyncServerService = &seedSyncServerService{
	syncSeedServerDriver: SeedSyncServerClient,
}

// 定时任务，定时从服务器同步用户是否可用
func (service *seedSyncServerService) CheckUser() error {
	success, err := service.syncSeedServerDriver.CheckUser()
	if err != nil {
		return err
	}
	if !success {
		common.CacheSet(common.USER_AVAILABLE_CACHE_KEY, common.USER_AVAILABLE_STATUS_UNAVAILABLE)
	} else {
		common.CacheSet(common.USER_AVAILABLE_CACHE_KEY, common.USER_AVAILABLE_STATUS_AVAILABLE)
	}
	return nil
}

// 从服务器同步支持的站点
func (service *seedSyncServerService) GetSupportedSite() error {
	supportedSites, err := service.syncSeedServerDriver.GetSupportedSites(0, math.MaxInt)
	if err != nil {
		return err
	}
	supportedSiteMap := make(map[string]SupportSiteResponse)
	for _, site := range supportedSites.Records {
		supportedSiteMap[site.SiteName] = site
	}
	// 缓存支持的站点，用于添加站点时校验
	common.CacheSetObject(common.SUPPORT_SITE_CACHE_KEY, supportedSiteMap)
	return nil
}
