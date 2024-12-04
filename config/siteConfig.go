package config

type SiteBaseConfig struct {
	UserAgent          string `mapstructure:"userAgent"`
	DownloadTorrentUrl string `mapstructure:"downloadTorrentUrl"`
	SeedDetailUrl      string `mapstructure:"seedDetailUrl"`
	PingUrl            string `mapstructure:"pingUrl"`
}

type SiteDriverConfig struct {
	SiteBaseConfig `mapstructure:",squash"` // 继承基础配置
	DriverType     string                   `mapstructure:"driverType"` // 如 "nexus", "unit3d"
}

type SpecificSiteConfig struct {
	SiteDriverConfig `mapstructure:",squash"` // 继承驱动配置
	SiteName         string                   `mapstructure:"siteName"`
}

// 继承形式的配置管理器 主要是考虑到很多站点的配置基本相同
type SiteConfigManager struct {
	BaseConfig    SiteBaseConfig                `mapstructure:"baseConfig"`    // 基础配置
	DriverConfigs map[string]SiteDriverConfig   `mapstructure:"driverConfigs"` // 驱动配置
	SiteConfigs   map[string]SpecificSiteConfig `mapstructure:"siteConfigs"`   // 具体站点配置
}

// 获取最终配置，实现配置继承
func (m *SiteConfigManager) GetSiteConfig(siteName string) SiteBaseConfig {
	siteConfig := m.SiteConfigs[siteName]
	driverConfig := m.DriverConfigs[siteConfig.DriverType]
	return m.mergeConfig(m.BaseConfig, driverConfig.SiteBaseConfig, siteConfig.SiteBaseConfig)
}

// 合并配置，后面的配置会覆盖前面的配置
func (m *SiteConfigManager) mergeConfig(base, driver, site SiteBaseConfig) SiteBaseConfig {
	result := base

	// 如果驱动层配置了值，则覆盖基础配置
	if driver.UserAgent != "" {
		result.UserAgent = driver.UserAgent
	}
	if driver.DownloadTorrentUrl != "" {
		result.DownloadTorrentUrl = driver.DownloadTorrentUrl
	}
	if driver.SeedDetailUrl != "" {
		result.SeedDetailUrl = driver.SeedDetailUrl
	}
	if site.PingUrl != "" {
		result.PingUrl = site.PingUrl
	}

	// 如果站点层配置了值，则覆盖驱动层配置
	if site.UserAgent != "" {
		result.UserAgent = site.UserAgent
	}
	if site.DownloadTorrentUrl != "" {
		result.DownloadTorrentUrl = site.DownloadTorrentUrl
	}
	if site.SeedDetailUrl != "" {
		result.SeedDetailUrl = site.SeedDetailUrl
	}
	if site.PingUrl != "" {
		result.PingUrl = site.PingUrl
	}

	return result
}
