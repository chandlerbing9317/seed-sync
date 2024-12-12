package nexus

import "seed-sync/site"

const (
    QingwaptSiteName = "qingwapt"
)

type QingwaptSite struct {
    *NexusSite
}

func (qingwaptSite *QingwaptSite) SiteName() string {
    return QingwaptSiteName
}

func NewQingwaptSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
    nexusSite, err := NewNexusSite(siteInfo)
    if err != nil {
        return nil, err
    }
    return &QingwaptSite{
        NexusSite: nexusSite.(*NexusSite),
    }, nil
}

// 注册站点
func init() {
    site.Factory.RegisterSite(QingwaptSiteName, NewQingwaptSite)
} 