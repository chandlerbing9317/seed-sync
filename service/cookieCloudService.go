package service

import (
	"encoding/json"
	"errors"
	"seed-sync/config"
	"seed-sync/driver"
	"sync"

	"gorm.io/gorm"
)

var (
	//对外暴露service
	CookieCloud *CookieCloudService
	once        sync.Once
)

type CookieCloudService struct {
	client *driver.CookieCloudClient
	db     *gorm.DB
	lock   sync.Mutex
}

func init() {
	once.Do(func() {
		CookieCloud = &CookieCloudService{
			db:   driver.DB,
			lock: sync.Mutex{},
		}
	})
}

const (
	CookieCloudConfigKey = "cookie_cloud_config_key"
)

// 添加或更新cookie cloud配置
func (service *CookieCloudService) AddOrUpdateCookieCloud(config *config.CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	service.lock.Lock()
	defer service.lock.Unlock()

	// 开启事务
	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新库
	if err := SaveOrUpdateSystemParamWithTx(tx, CookieCloudConfigKey, string(configBytes)); err != nil {
		tx.Rollback()
		return err
	}

	var cookieCloudClient *driver.CookieCloudClient
	if service.client != nil {
		// 走更新流程
		if cookieCloudClient, err = service.client.Update(config); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 走创建流程
		if cookieCloudClient, err = driver.NewCookieCloudClient(config); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 所有操作成功后提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	service.client = cookieCloudClient
	return nil
}

// 删除cookie cloud配置
func (service *CookieCloudService) DeleteCookieCloud() error {
	service.lock.Lock()
	defer service.lock.Unlock()

	if service.client != nil {
		// 开启事务
		tx := service.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 删库
		if err := DeleteSystemParamWithTx(tx, CookieCloudConfigKey); err != nil {
			tx.Rollback()
			return err
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}

		// 删客户端
		service.client.Destroy()
		service.client = nil
	}
	return nil
}

// 获取cookie cloud配置
func (service *CookieCloudService) GetCookieCloudConfig() (*config.CookieCloudConfig, error) {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client == nil {
		return nil, errors.New("未配置cookie cloud")
	}
	return service.client.GetConfig(), nil
}
