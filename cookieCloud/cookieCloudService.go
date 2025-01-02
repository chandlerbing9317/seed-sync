package cookieCloud

import (
	"errors"
	"fmt"
	"seed-sync/common"
	"seed-sync/log"
	"seed-sync/scheduler"
	"seed-sync/seedSyncServer"
	"seed-sync/site"
	"sync"

	"go.uber.org/zap"
)

// 初始化且对外暴露的单例service
var CookieCloudService *cookieCloudService
var once sync.Once

func InitCookieCloud() {
	once.Do(func() {
		CookieCloudService = &cookieCloudService{
			cookieCloudDAO: cookieCloudDAO,
			lock:           sync.Mutex{},
		}
		//查询库中的cookie cloud配置 如果存在就初始化client，用于后续使用
		config := CookieCloudService.cookieCloudDAO.GetCookieCloudConfig()
		if config != nil {
			client, err := NewCookieCloudClient(config)
			if err != nil {
				log.Fatal("初始化cookie cloud失败", zap.Error(err))
			}
			CookieCloudService.client = client
		}
	})
}

type cookieCloudService struct {
	client         *CookieCloudClient
	cookieCloudDAO *CookieCloudDAO
	lock           sync.Mutex
}

func (service *cookieCloudService) CreateCookieCloud(config *CookieCloudConfig) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client != nil {
		return errors.New("cookie cloud已配置")
	}
	//初始化客户端
	client, err := NewCookieCloudClient(config)
	if err != nil {
		return err
	}
	service.client = client
	//保存配置到数据库
	return service.cookieCloudDAO.CreateCookieCloudConfig(config)
}

// 添加或更新cookie cloud配置
func (service *cookieCloudService) UpdateCookieCloud(config *CookieCloudConfig) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client == nil {
		return errors.New("未配置cookie cloud")
	}
	//更新客户端
	client, err := service.client.Update(config)
	if err != nil {
		return err
	}
	service.client = client
	//更新配置到数据库
	return service.cookieCloudDAO.UpdateCookieCloudConfig(config)
}

// 删除cookie cloud配置
func (service *cookieCloudService) DeleteCookieCloud() error {
	service.lock.Lock()
	defer service.lock.Unlock()

	if service.client != nil {
		service.client.Destroy()
		service.cookieCloudDAO.DeleteCookieCloudConfig()
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

// 同步cookieCloud的cookie到站点
// 如果站点不存在，就自动创建站点

func (service *cookieCloudService) SyncCookie() error {
	if service.client == nil {
		return errors.New("未配置cookie cloud")
	}
	cookie, err := service.client.GetCookie()
	if err != nil {
		return err
	}
	//todo: 删除直接获取，改为定时轮询从内存中获取支持的站点
	err = seedSyncServer.SeedSyncServerService.GetSupportedSite()
	if err != nil {
		return err
	}
	supportedSiteMap, exists := common.CacheGetObject[map[string]seedSyncServer.SupportSiteResponse](common.SUPPORT_SITE_CACHE_KEY)
	if !exists {
		return errors.New("获取支持的站点失败")
	}

	var errs []error
	for _, siteInfo := range supportedSiteMap {
		for _, host := range siteInfo.Hosts {
			cookieStr, ok := cookie.GetCookieByDomain(host)
			if ok {
				err = site.SiteService.UpdateCookie(siteInfo.SiteName, cookieStr, host)
				if err != nil {
					log.Error("更新站点cookie失败", zap.String("siteName", siteInfo.SiteName), zap.Error(err))
					errs = append(errs, err)
				}
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("更新站点cookie失败: %v", errs)
	}
	return nil
}

func (service *cookieCloudService) SyncCookieScheduler(schedulerTask *scheduler.SchedulerTaskTable) error {
	return service.SyncCookie()
}
