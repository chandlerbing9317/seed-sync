package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Log LogConfig `mapstructure:"log"`
}
type LogConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
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
		err = viper.Unmarshal(Conf)
		if err != nil {
			panic("配置文件读取失败" + err.Error()) // 映射过程中的错误处理
		}
	})
}
