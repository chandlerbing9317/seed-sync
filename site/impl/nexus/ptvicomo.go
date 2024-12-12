package nexus

import "seed-sync/site"

const (
    PtvicomoSiteName = "ptvicomo"
)

type PtvicomoSite struct {
    *NexusSite
}

func (ptvicomoSite *PtvicomoSite) SiteName() string {
    return PtvicomoSiteName
}

func NewPtvicomoSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &PtvicomoSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(PtvicomoSiteName, NewPtvicomoSite)
} 