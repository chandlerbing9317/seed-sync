package nexus

import "seed-sync/site"

const (
    OurbitsSiteName = "ourbits"
)

type OurbitsSite struct {
    *NexusSite
}

func (ourbitsSite *OurbitsSite) SiteName() string {
    return OurbitsSiteName
}

func NewOurbitsSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &OurbitsSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(OurbitsSiteName, NewOurbitsSite)
} 