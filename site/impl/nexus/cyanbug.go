package nexus

import "seed-sync/site"

const (
	CyanbugSiteName = "cyanbug"
)

type CyanbugSite struct {
	*NexusSite
}

func (cyanbugSite *CyanbugSite) SiteName() string {
	return CyanbugSiteName
}

func NewCyanbugSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	cyanbug := &CyanbugSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	cyanbug.BaseSite.SetImplementor(cyanbug)
	return cyanbug, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(CyanbugSiteName, NewCyanbugSite)
}
