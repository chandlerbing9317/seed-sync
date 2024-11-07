package service

import (
	"errors"
	"seed-sync/driver/client"
	"seed-sync/driver/db"
	"sync"

	"gorm.io/gorm"
)

// 初始化且对外暴露的单例service
var CookieCloud = &CookieCloudService{
	cookieCloudDAO: db.CookieCloudDao,
	db:             db.DB,
	lock:           sync.Mutex{},
}

type CookieCloudService struct {
	client         *client.CookieCloudClient
	cookieCloudDAO *db.CookieCloudDAO
	db             *gorm.DB
	lock           sync.Mutex
}

// 添加或更新cookie cloud配置
func (service *CookieCloudService) AddOrUpdateCookieCloud(config *db.CookieCloudConfig) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	// 开启事务
	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 操作库
	if err := service.cookieCloudDAO.AddOrUpdateCookieCloudConfigWithTx(tx, config); err != nil {
		tx.Rollback()
		return err
	}
	//操作client
	var cookieCloudClient *client.CookieCloudClient
	var err error
	if service.client != nil {
		// 走更新流程
		if cookieCloudClient, err = service.client.Update(config); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 走创建流程
		if cookieCloudClient, err = client.NewCookieCloudClient(config); err != nil {
			tx.Rollback()
			return err
		}
	}
	//更新client
	service.client = cookieCloudClient

	// 所有操作成功后提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
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
		if err := service.cookieCloudDAO.DeleteCookieCloudConfigWithTx(tx); err != nil {
			tx.Rollback()
			return err
		}
		// 删客户端
		service.client.Destroy()
		service.client = nil

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

// 获取cookie cloud配置
func (service *CookieCloudService) GetCookieCloudConfig() (*db.CookieCloudConfig, error) {
	service.lock.Lock()
	defer service.lock.Unlock()
	if service.client == nil {
		return nil, errors.New("未配置cookie cloud")
	}
	return service.client.GetConfig(), nil
}
