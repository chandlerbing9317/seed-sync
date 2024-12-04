package seedSyncServer

type SeedSyncRequest struct {
	//要辅种的种子
	InfoHash []string `json:"infoHash"`
	//辅种的站点
	Sites []string `json:"sites"`
}
type SeedSyncTorrentInfoResponse struct {
	SiteName  string `json:"siteName"`
	TorrentId int    `json:"torrentId"`
	InfoHash  string `json:"infoHash"`
}
type SupportSiteResponse struct {
	SiteName string   `json:"siteName"`
	ShowName string   `json:"showName"`
	Hosts    []string `json:"hosts"`
}
type UserAuthRequest struct {
	Token      string `json:"token"`
	SiteName   string `json:"siteName"`
	SiteUserID string `json:"siteUserId"`
}
