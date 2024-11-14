package seedSyncModel

type UserAuthRequest struct {
	Token      string `json:"token"`
	SiteName   string `json:"siteName"`
	SiteUserID string `json:"siteUserId"`
}
