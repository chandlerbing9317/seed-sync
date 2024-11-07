package config

type ProxyConfig struct {
	ProxyURL      string `mapstructure:"proxyUrl"`
	ProxyUsername string `mapstructure:"proxyUsername"`
	ProxyPassword string `mapstructure:"proxyPassword"`
}
