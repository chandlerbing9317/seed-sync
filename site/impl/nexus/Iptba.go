package nexus

import "seed-sync/site"

const (
	IptbaSiteName = "1ptba"
)

type IptbaSite struct {
	*NexusSite
}

func (ptbaSite *IptbaSite) SiteName() string {
	return IptbaSiteName
}

func New1PtbaSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	return &IptbaSite{
		NexusSite: nexusSite.(*NexusSite),
	}, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(IptbaSiteName, New1PtbaSite)
}
