package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	LogConfig    LogConfig    `mapstructure:"log"`
	ProxyConfig  ProxyConfig  `mapstructure:"proxy"`
	ServerConfig ServerConfig `mapstructure:"server"`
	// QbittorrentClientConfig  HttpClientConfig `mapstructure:"qbittorrentClient"`
	// TransmissionClientConfig HttpClientConfig `mapstructure:"transmissionClient"`
	// SiteClientConfig         HttpClientConfig `mapstructure:"siteClient"`
	ServerClientConfig HttpClientConfig `mapstructure:"serverClient"`
}

func init() {
	InitConfig()
}

var once sync.Once
var Conf *Config

// 配置文件初始化
// todo 从环境变量读取数据 便于docker等部署环境
func InitConfig() {
	once.Do(func() {
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		// 读取log配置
		viper.SetConfigName("log")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		Conf = &Config{}
		err = viper.Unmarshal(&Conf.LogConfig)
		if err != nil {
			panic("配置文件读取失败" + err.Error()) // 映射过程中的错误处理
		}
		//读取proxy配置
		viper.SetConfigName("proxy")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(&Conf.ProxyConfig)
		if err != nil {
			panic(err)
		}
		//读取server配置
		viper.SetConfigName("server")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(&Conf.ServerConfig)
		if err != nil {
			panic(err)
		}
		// //读取qbittorrentClient配置
		// viper.SetConfigName("qbittorrentClient")
		// err = viper.ReadInConfig()
		// if err != nil {
		// 	panic(err)
		// }
		// err = viper.Unmarshal(&Conf.QbittorrentClientConfig)
		// if err != nil {
		// 	panic(err)
		// }
		// //读取transmissionClient配置
		// viper.SetConfigName("transmissionClient")
		// err = viper.ReadInConfig()
		// if err != nil {
		// 	panic(err)
		// }
		// err = viper.Unmarshal(&Conf.TransmissionClientConfig)
		// if err != nil {
		// 	panic(err)
		// }
		// //读取siteClient配置
		// viper.SetConfigName("siteClient")
		// err = viper.ReadInConfig()
		// if err != nil {
		// 	panic(err)
		// }
		// err = viper.Unmarshal(&Conf.SiteClientConfig)
		// if err != nil {
		// 	panic(err)
		// }
		//读取serverClient配置
		viper.SetConfigName("serverClient")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(&Conf.ServerClientConfig)
		if err != nil {
			panic(err)
		}
	})
}
