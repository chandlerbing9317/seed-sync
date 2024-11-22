package seedSyncServer

import (
	"math"
	"seed-sync/common"
	"seed-sync/site"
	"strings"
)

type seedSyncServerService struct {
	syncSeedServerDriver *ServerClient
}

var SeedSyncServerService = &seedSyncServerService{}

// 定时任务，定时从服务器同步用户是否可用
func (service *seedSyncServerService) CheckUserForSchedulerTask() error {
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

// 定时任务，定时从服务器同步支持的站点
func (service *seedSyncServerService) GetSiteForSchedulerTask() error {
	supportedSites, err := service.syncSeedServerDriver.GetSupportedSites(0, math.MaxInt)
	if err != nil {
		return err
	}
	sites, err := site.SiteService.GetSiteList()
	if err != nil {
		return err
	}
	supportedSiteMap := make(map[string]SupportSiteResponse)
	for _, site := range supportedSites.Records {
		supportedSiteMap[site.SiteName] = site
	}
	//更新从服务端同步过来的站点信息
	//主要是一些站点的爬虫配置
	sitesToUpdate := make([]*site.SiteTable, 0, len(sites))
	for _, site := range sites {
		if supportedSite, ok := supportedSiteMap[site.SiteName]; ok && site.IsOverride {
			site.Domains = strings.Join(supportedSite.Domain, ";")
			site.UserAgent = supportedSite.UserAgent
			//尽量保留用户自定义的header，但服务端返回的优先级更高
			headerMap := make(map[string]string)
			for _, header := range strings.Split(site.CustomHeader, common.CustomHeaderSeparator) {
				headerMap[header] = header
			}
			for _, header := range supportedSite.CustomHeader {
				headerMap[header] = header
			}
			headers := make([]string, 0, len(headerMap))
			for header := range headerMap {
				headers = append(headers, header)
			}
			site.CustomHeader = strings.Join(headers, common.CustomHeaderSeparator)
			site.SeedListUrl = supportedSite.SeedListUrl
			site.RssUrl = supportedSite.RssUrl
			site.DownloadUrl = supportedSite.DownloadUrl
			site.DetailUrl = supportedSite.DetailUrl
			site.PingUrl = supportedSite.PingUrl
			sitesToUpdate = append(sitesToUpdate, site)
		}
	}
	if len(sitesToUpdate) > 0 {
		err = site.SiteService.UpdateBatchSite(sitesToUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}
