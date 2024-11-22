package common

const (
	CustomHeaderSeparator = ";seed-sync-custom-header-separator;"
	Http                  = "http://"
	Https                 = "https://"
)

// two factor type
const (
	TwoFactorTypeEmail = "email"
	TwoFactorTypeAuth  = "auth"
)

// scheduler task execute content
const (
	SYNC_COOKIE_CLOUD_EXECUTE_CONTENT = "sync-cookie-cloud-execute-content"
	CHECK_USER_EXECUTE_CONTENT        = "check-user-execute-content"
	GET_SITE_EXECUTE_CONTENT          = "get-site-execute-content"
	SYNC_SEED_EXECUTE_CONTENT         = "sync-seed-execute-content"
)

// user available status
const (
	USER_AVAILABLE_STATUS_AVAILABLE   = "AVAILABLE"
	USER_AVAILABLE_STATUS_UNAVAILABLE = "UNAVAILABLE"
)

// cache key
const (
	USER_AVAILABLE_CACHE_KEY = "user-available-cache-key"
	SUPPORT_SITE_CACHE_KEY   = "support-site-cache-key"
)

// downloader type
const (
	DOWNLOADER_TYPE_TRANSMISSION = "transmission"
	DOWNLOADER_TYPE_QBITTORRENT  = "qbittorrent"
)
