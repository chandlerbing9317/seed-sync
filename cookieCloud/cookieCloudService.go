package cookieCloud

import (
	"errors"
	"seed-sync/site"
	"sync"
)

// 初始化且对外暴露的单例service
var CookieCloudService = &cookieCloudService{
	cookieCloudDAO: cookieCloudDAO,
	lock:           sync.Mutex{},
}

type cookieCloudService struct {
	client         *CookieCloudClient
	cookieCloudDAO *CookieCloudDAO
	lock           sync.Mutex
}

// 添加或更新cookie cloud配置
func (service *cookieCloudService) CreateOrUpdateCookieCloud(config *CookieCloudConfig) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client != nil {
		service.cookieCloudDAO.UpdateCookieCloudConfig(config)
		service.client.Update(config)
	} else {
		var err error
		service.client, err = NewCookieCloudClient(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除cookie cloud配置
func (service *cookieCloudService) DeleteCookieCloud() error {
	service.lock.Lock()
	defer service.lock.Unlock()

	if service.client != nil {
		service.cookieCloudDAO.DeleteCookieCloudConfig()
		service.client.Destroy()
		service.client = nil
	}
	return nil
}

// 获取cookie cloud配置
func (service *cookieCloudService) GetCookieCloudConfig() (*CookieCloudConfig, error) {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client == nil {
		return nil, errors.New("未配置cookie cloud")
	}
	return service.client.GetConfig(), nil
}

// sync cookieCloud for scheduler task
func (service *cookieCloudService) SyncCookieForSchedulerTask() error {
	if service.client == nil {
		return errors.New("未配置cookie cloud")
	}
	cookie, err := service.client.GetCookie()
	if err != nil {
		return err
	}
	siteList, err := site.SiteService.GetSiteList()
	if err != nil {
		return err
	}
	sites := make([]*site.SiteTable, 0, len(siteList))
	for _, site := range siteList {
		cookieStr, ok := cookie.GetCookieByDomain(site.Host)
		if ok {
			site.Cookie = cookieStr
			sites = append(sites, site)
		}
	}
	if len(sites) > 0 {
		err = site.SiteService.UpdateBatchSite(sites)
		if err != nil {
			return err
		}
	}
	return nil
}
