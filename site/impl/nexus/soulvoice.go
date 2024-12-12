package nexus

import "seed-sync/site"

const (
    SoulvoiceSiteName = "soulvoice"
)

type SoulvoiceSite struct {
    *NexusSite
}

func (soulvoiceSite *SoulvoiceSite) SiteName() string {
    return SoulvoiceSiteName
}

func NewSoulvoiceSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &SoulvoiceSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(SoulvoiceSiteName, NewSoulvoiceSite)
} 