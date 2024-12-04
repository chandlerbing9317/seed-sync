package site

import (
	"errors"
	"seed-sync/common"
	"strings"
	"sync"
	"time"
)

type siteService struct {
	siteDao *SiteDAO
	lock    sync.Mutex
}

var SiteService = &siteService{
	siteDao: siteDAO,
	lock:    sync.Mutex{},
}

func (service *siteService) AddSite(request *AddSiteRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if site := service.siteDao.GetSiteInfo(request.SiteName); site != nil {
		return errors.New("站点" + site.ShowName + "已存在")
	}

	customHeader := ""
	if len(request.CustomHeaders) > 0 {
		customHeader = strings.Join(request.CustomHeaders, common.CustomHeaderSeparator)
	}

	//todo order
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
	if err := service.siteDao.AddSite(siteTable, siteFlowControl); err != nil {
		return err
	}
	return nil
}

func (service *siteService) GetSiteList() ([]*SiteInfo, error) {
	return service.siteDao.GetAllSites()
}
