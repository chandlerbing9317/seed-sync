package nexus

import "seed-sync/site"

const (
	HhanSiteName = "hhan"
)

type HhanSite struct {
	*NexusSite
}

func (hhanSite *HhanSite) SiteName() string {
	return HhanSiteName
}

func NewHhanSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
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
