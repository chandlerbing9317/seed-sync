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
	qingwapt := &QingwaptSite{
		NexusSite: nexusSite.(*NexusSite),
	}
	qingwapt.BaseSite.SetImplementor(qingwapt)
	return qingwapt, nil
}

// 注册站点
func init() {
	site.Factory.RegisterSite(QingwaptSiteName, NewQingwaptSite)
}
