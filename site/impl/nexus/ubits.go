package nexus

import "seed-sync/site"

const (
	UbitsSiteName = "ubits"
)

type UbitsSite struct {
	*NexusSite
}

func (ubitsSite *UbitsSite) SiteName() string {
	return UbitsSiteName
}

func NewUbitsSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	ubits := &UbitsSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	ubits.BaseSite.SetImplementor(ubits)
	return ubits, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(UbitsSiteName, NewUbitsSite)
}
