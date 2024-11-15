package site

import (
	"fmt"
	"seed-sync/driver/db"
	"sync"
)

type SiteFactory struct {
	siteConstructors map[string]SiteConstructor
	mu               sync.RWMutex
}

type SiteConstructor func(siteInfo *db.Site) (Site, error)

var Factory *SiteFactory = &SiteFactory{
	siteConstructors: make(map[string]SiteConstructor),
}

// 注册站点
func (f *SiteFactory) RegisterSite(siteName string, constructor SiteConstructor) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.siteConstructors[siteName] = constructor
}

// 创建站点
func (f *SiteFactory) CreateSite(siteInfo *db.Site) (Site, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	constructor, ok := f.siteConstructors[siteInfo.SiteName]
	if !ok {
		return nil, fmt.Errorf("不支持的站点: %s", siteInfo.SiteName)
	}
	return constructor(siteInfo)
}
