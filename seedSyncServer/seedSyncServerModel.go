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
	SiteName     string   `json:"siteName"`
	ShowName     string   `json:"showName"`
	Host         string   `json:"host"`
	Domain       []string `json:"domain"`
	UserAgent    string   `json:"userAgent"`
	CustomHeader []string `json:"customHeader"`
	SeedListUrl  string   `json:"seedListUrl"`
	RssUrl       string   `json:"rssUrl"`
	DetailUrl    string   `json:"detailUrl"`
	DownloadUrl  string   `json:"downloadUrl"`
	PingUrl      string   `json:"pingUrl"`
}
type UserAuthRequest struct {
	Token      string `json:"token"`
	SiteName   string `json:"siteName"`
	SiteUserID string `json:"siteUserId"`
}
