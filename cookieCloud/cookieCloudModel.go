package cookieCloud

type CookieCloudConfig struct {
	Url string `json:"url"`
	//用户KEY
	UserKey string `json:"user_key"`
	//端对端加密密码
	P2pPassword string `json:"p2p_password"`
	//同步cron表达式
	SyncCron string `json:"sync_cron"`
}
