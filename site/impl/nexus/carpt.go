package nexus

import "seed-sync/site"

const (
	CarptSiteName = "carpt"
)

type CarptSite struct {
	*NexusSite
}

func (carptSite *CarptSite) SiteName() string {
	return CarptSiteName
}

func NewCarptSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	carpt := &CarptSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	carpt.BaseSite.SetImplementor(carpt)
	return carpt, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(CarptSiteName, NewCarptSite)
}
