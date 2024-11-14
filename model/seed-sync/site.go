package seedSyncModel

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
