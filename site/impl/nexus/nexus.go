package nexus

import (
	"fmt"
	"seed-sync/common"
	"seed-sync/site"
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

func (nexusSite *NexusSite) GetSeedListUrl(page int) string {
	return common.Https + nexusSite.SiteInfo.Host + fmt.Sprintf(nexusSite.SiteInfo.SeedListUrl, page)
}

func NewNexusSite(siteInfo *site.SiteTable) (site.Site, error) {
	return &NexusSite{
		BaseSite: site.NewBaseSite(siteInfo),
	}, nil
}
