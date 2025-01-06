package nexus

import "seed-sync/site"

const (
	AudiencesSiteName = "audiences"
)

type AudiencesSite struct {
	*NexusSite
}

func (audiencesSite *AudiencesSite) SiteName() string {
	return AudiencesSiteName
}

func NewAudiencesSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	nexusSite, err := NewNexusSite(siteInfo)
	if err != nil {
		return nil, err
	}
	audiences := &AudiencesSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	audiences.BaseSite.SetImplementor(audiences)
	return audiences, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(AudiencesSiteName, NewAudiencesSite)
}
