package unit3d

import (
	"fmt"
	"seed-sync/common"
	"seed-sync/driver/client/site"
)

// unit3d站点，对BaseSite接口的实现
type Unit3dSite struct {
	*site.BaseSite
}


// unit3d实现接口
func (unit3dSite *Unit3dSite) GetDownloadUrl(torrentId int) string {
	return common.Https + unit3dSite.SiteInfo.Host + fmt.Sprintf(unit3dSite.SiteInfo.DownloadUrl, torrentId, unit3dSite.SiteInfo.RssKey)
}
func (unit3dSite *Unit3dSite) GetPingUrl() string {
	return common.Https + unit3dSite.SiteInfo.Host + unit3dSite.SiteInfo.PingUrl
}
