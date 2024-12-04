package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	LogConfig        LogConfig        `mapstructure:"log"`
	ProxyConfig      ProxyConfig      `mapstructure:"proxy"`
	ServerConfig     ServerConfig     `mapstructure:"server"`
	HttpClientConfig HttpClientConfig `mapstructure:"httpClient"`
	SiteConfig       SiteConfigManager `mapstructure:"site"`
}

func init() {
	InitConfig()
}

var once sync.Once
var Conf *Config

func InitConfig() {
	once.Do(func() {
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		// 读取各个配置文件
		configFiles := []string{"log", "proxy", "server", "httpClient", "site"}
		Conf = &Config{}

		for _, configName := range configFiles {
			viper.SetConfigName(configName)
			err := viper.ReadInConfig()
			if err != nil {
				panic(err)
			}

			switch configName {
			case "log":
				err = viper.Unmarshal(&Conf.LogConfig)
			case "proxy":
				err = viper.Unmarshal(&Conf.ProxyConfig)
			case "server":
				err = viper.Unmarshal(&Conf.ServerConfig)
			case "httpClient":
				err = viper.Unmarshal(&Conf.HttpClientConfig)
			case "site":
				err = viper.Unmarshal(&Conf.SiteConfig)
			}

			if err != nil {
				panic("配置文件读取失败: " + configName + ", " + err.Error())
			}
		}
	})
}
