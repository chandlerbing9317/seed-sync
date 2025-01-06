package nexus

import "seed-sync/site"

const (
	CrabptSiteName = "crabpt"
)

type CrabptSite struct {
	*NexusSite
}

func (crabptSite *CrabptSite) SiteName() string {
	return CrabptSiteName
}

func NewCrabptSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	crabpt := &CrabptSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	crabpt.BaseSite.SetImplementor(crabpt)
	return crabpt, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(CrabptSiteName, NewCrabptSite)
}
