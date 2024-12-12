package nexus

import "seed-sync/site"

const (
    BtschoolSiteName = "btschool"
)

type BtschoolSite struct {
    *NexusSite
}

func (btschoolSite *BtschoolSite) SiteName() string {
    return BtschoolSiteName
}

func NewBtschoolSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &BtschoolSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(BtschoolSiteName, NewBtschoolSite)
} 