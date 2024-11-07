package nexus

import (
	"fmt"
	"seed-sync/common"
	"seed-sync/driver/client/site"
	"seed-sync/driver/db"
)

// nexus站点，对BaseSite接口的实现
type NexusSite struct {
	*site.BaseSite
}

// nexus实现接口
func (nexusSite *NexusSite) GetDownloadUrl(torrentId int) string {
	return common.Https + nexusSite.SiteInfo.Host + fmt.Sprintf(nexusSite.SiteInfo.DownloadUrl, torrentId)
}
func (nexusSite *NexusSite) GetPingUrl() string {
	return common.Https + nexusSite.SiteInfo.Host + nexusSite.SiteInfo.PingUrl
}

func NewNexusSite(siteInfo *db.Site) (site.Site, error) {
	return &NexusSite{
		BaseSite: site.NewBaseSite(siteInfo),
	}, nil
}
