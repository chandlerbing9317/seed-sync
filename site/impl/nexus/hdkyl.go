package nexus

import "seed-sync/site"

const (
	HdkylSiteName = "hdkyl"
)

type HdkylSite struct {
	*NexusSite
}

func (hdkylSite *HdkylSite) SiteName() string {
	return HdkylSiteName
}

func NewHdkylSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	hdkyl := &HdkylSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	hdkyl.BaseSite.SetImplementor(hdkyl)
	return hdkyl, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(HdkylSiteName, NewHdkylSite)
}
