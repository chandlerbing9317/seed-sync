package site

import (
	"errors"
	"seed-sync/common"
	"strings"
	"sync"
	"time"
)

type siteService struct {
	siteDao       *SiteDAO
	siteClientMap map[string]SiteClient
	lock          sync.Mutex
}

var SiteService = &siteService{
	siteDao: siteDAO,
	lock:    sync.Mutex{},
}

// 添加站点
func (service *siteService) AddSite(request *AddSiteRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()

	// 开启事务
	tx := service.siteDao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查站点是否存在
	if service.siteDao.GetSiteInfo(request.SiteName) != nil {
		tx.Rollback()
		return errors.New("站点" + request.SiteName + "已存在")
	}

	// 处理自定义header
	customHeader := ""
	if len(request.CustomHeaders) > 0 {
		customHeader = strings.Join(request.CustomHeaders, common.CustomHeaderSeparator)
	}

	// 获取最大order
	maxOrder, err := service.siteDao.GetMaxOrderTx(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 准备站点数据
	siteTable := &SiteTable{
		SiteName:     request.SiteName,
		ShowName:     request.ShowName,
		Host:         request.Host,
		Cookie:       request.Cookie,
		ApiToken:     request.ApiToken,
		Passkey:      request.Passkey,
		RssKey:       request.RssKey,
		UserAgent:    request.UserAgent,
		CustomHeader: customHeader,
		Proxy:        request.Proxy,
		IsActive:     request.IsActive,
		Timeout:      request.Timeout,
		Order:        maxOrder + 1,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	siteFlowControl := &SiteFlowControl{
		SiteName:   request.SiteName,
		MaxPerMin:  request.MaxPerMin,
		MaxPerHour: request.MaxPerHour,
		MaxPerDay:  request.MaxPerDay,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 在事务中添加站点
	if err := service.siteDao.AddSiteTx(tx, siteTable, siteFlowControl); err != nil {
		tx.Rollback()
		return err
	}

	//创建站点
	siteClient, err := Factory.CreateSite(GenerateSiteInfo(siteTable, siteFlowControl))
	if err != nil {
		tx.Rollback()
		return err
	}
	service.siteClientMap[siteTable.SiteName] = siteClient

	// 提交事务
	return tx.Commit().Error
}

// 更新站点
func (service *siteService) UpdateSite(request *UpdateSiteRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()

	// 开启事务
	tx := service.siteDao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查站点是否存在
	if service.siteDao.GetSiteInfo(request.SiteName) == nil {
		tx.Rollback()
		return errors.New("站点" + request.SiteName + "不存在")
	}
	if _, ok := service.siteClientMap[request.SiteName]; !ok {
		tx.Rollback()
		return errors.New("站点" + request.SiteName + "不存在")
	}

	// 处理自定义header
	customHeader := ""
	if len(request.CustomHeaders) > 0 {
		customHeader = strings.Join(request.CustomHeaders, common.CustomHeaderSeparator)
	}
	siteTable := &SiteTable{
		SiteName:     request.SiteName,
		ShowName:     request.ShowName,
		Order:        request.Order,
		Host:         request.Host,
		Cookie:       request.Cookie,
		ApiToken:     request.ApiToken,
		Passkey:      request.Passkey,
		RssKey:       request.RssKey,
		UserAgent:    request.UserAgent,
		CustomHeader: customHeader,
		Proxy:        request.Proxy,
		IsActive:     request.IsActive,
		Timeout:      request.Timeout,
		UpdateTime:   time.Now(),
	}

	siteFlowControl := &SiteFlowControl{
		SiteName:   request.SiteName,
		MaxPerMin:  request.MaxPerMin,
		MaxPerHour: request.MaxPerHour,
		MaxPerDay:  request.MaxPerDay,
		UpdateTime: time.Now(),
	}

	if err := service.siteDao.UpdateSiteBySiteNameTx(tx, request.SiteName, siteTable, siteFlowControl); err != nil {
		tx.Rollback()
		return err
	}
	siteClient := service.siteClientMap[request.SiteName]
	err := siteClient.Update(GenerateSiteInfo(siteTable, siteFlowControl))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 批量更新站点顺序
func (service *siteService) BatchUpdateSiteOrders(updates []SiteOrderUpdateRequest) error {
	return service.siteDao.BatchUpdateSiteOrders(updates)
}

// 获取所有站点
func (service *siteService) GetSiteList() ([]*SiteInfo, error) {
	return service.siteDao.GetAllSites()
}

// 更新站点cookie
func (service *siteService) UpdateCookie(siteName string, cookie string) error {
	return service.siteDao.UpdateCookie(siteName, cookie)
}

func (service *siteService) Ping(siteName string) error {
	siteClient := service.siteClientMap[siteName]
	return siteClient.Ping()
}
