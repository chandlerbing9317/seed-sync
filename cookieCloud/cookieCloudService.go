package cookieCloud

import (
	"errors"
	"fmt"
	"seed-sync/log"
	"seed-sync/site"
	"sync"

	"go.uber.org/zap"
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
	var errs []error
	for _, siteInfo := range siteList {
		cookieStr, ok := cookie.GetCookieByDomain(siteInfo.Host)
		if ok {
			err = site.SiteService.UpdateCookie(siteInfo.SiteName, cookieStr)
			if err != nil {
				log.Error("更新站点cookie失败", zap.String("siteName", siteInfo.SiteName), zap.Error(err))
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("更新站点cookie失败: %v", errs)
	}
	return nil
}
