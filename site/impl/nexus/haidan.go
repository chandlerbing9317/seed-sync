package nexus

import "seed-sync/site"

const (
	HaidanSiteName = "haidan"
)

type HaidanSite struct {
	*NexusSite
}

func (haidanSite *HaidanSite) SiteName() string {
	return HaidanSiteName
}

func NewHaidanSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	haidan := &HaidanSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	haidan.BaseSite.SetImplementor(haidan)
	return haidan, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(HaidanSiteName, NewHaidanSite)
}
