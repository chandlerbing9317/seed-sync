package nexus

import "seed-sync/site"

const (
    HdfansSiteName = "hdfans"
)

type HdfansSite struct {
    *NexusSite
}

func (hdfansSite *HdfansSite) SiteName() string {
    return HdfansSiteName
}

func NewHdfansSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &HdfansSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(HdfansSiteName, NewHdfansSite)
} 