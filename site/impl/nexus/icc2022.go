package nexus

import "seed-sync/site"

const (
	Icc2022SiteName = "icc2022"
)

type Icc2022Site struct {
	*NexusSite
}

func (icc2022Site *Icc2022Site) SiteName() string {
	return Icc2022SiteName
}

func NewIcc2022Site(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	icc2022 := &Icc2022Site{
		NexusSite: nexusSite.(*NexusSite),
	}
	icc2022.BaseSite.SetImplementor(icc2022)
	return icc2022, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(Icc2022SiteName, NewIcc2022Site)
}
