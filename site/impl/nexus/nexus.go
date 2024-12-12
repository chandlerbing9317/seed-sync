package nexus

import (
	"seed-sync/common"
	"seed-sync/site"
	"strconv"
)

// nexus站点，对BaseSite接口的实现
type NexusSite struct {
	*site.BaseSite
}

// nexus实现接口
func (nexusSite *NexusSite) GetDownloadUrl(torrentId int) string {
	return nexusSite.SiteInfo.Host +
		common.FormatUrlTemplate(nexusSite.Config.DownloadTorrentUrl, map[string]string{"torrentId": strconv.Itoa(torrentId)})
}
func (nexusSite *NexusSite) GetPingUrl() string {
	return nexusSite.SiteInfo.Host + nexusSite.Config.PingUrl
}

// nexus header加入cookie
func (nexusSite *NexusSite) GetHttpHeader() map[string]string {
	header := nexusSite.BaseSite.GetHttpHeader()
	header["Cookie"] = nexusSite.SiteInfo.Cookie
	return header
}

func NewNexusSite(siteInfo *site.SiteInfo) (site.SiteClient, error) {
	return &NexusSite{
		BaseSite: site.NewBaseSite(siteInfo),
	}, nil
}
