package nexus

import (
	"seed-sync/driver/client/site"
	"seed-sync/driver/db"
)

type HhanSite struct {
	*NexusSite
}

const (
	HhanSiteName = "hhan"
)

// hhan
func (hhanSite *HhanSite) SiteName() string {
	return HhanSiteName
}

func NewHhanSite(siteInfo *db.Site) (site.Site, error) {
	//先实例化父类
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	return &HhanSite{
		NexusSite: nexusSite.(*NexusSite),
	}, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(HhanSiteName, NewHhanSite)
}
