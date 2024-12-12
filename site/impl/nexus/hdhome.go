package nexus

import "seed-sync/site"

const (
    HdhomeSiteName = "hdhome"
)

type HdhomeSite struct {
    *NexusSite
}

func (hdhomeSite *HdhomeSite) SiteName() string {
    return HdhomeSiteName
}

func NewHdhomeSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &HdhomeSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(HdhomeSiteName, NewHdhomeSite)
} 