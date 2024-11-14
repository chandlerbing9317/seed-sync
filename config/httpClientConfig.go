package config

import "time"

// http client 配置
type HttpClientConfig struct {
	Proxy               bool          `mapstructure:"proxy"`
	ConnTimeout         time.Duration `mapstructure:"connTimeout"`
	ReadTimeout         time.Duration `mapstructure:"readTimeout"`
	WriteTimeout        time.Duration `mapstructure:"writeTimeout"`
	KeepAlive           time.Duration `mapstructure:"keepAlive"`
	MaxIdleConns        int           `mapstructure:"maxIdleConns"`
	MaxIdleConnsPerHost int           `mapstructure:"maxIdleConnsPerHost"`
	IdleConnTimeout     time.Duration `mapstructure:"idleConnTimeout"`
	TLSHandshakeTimeout time.Duration `mapstructure:"tlsHandshakeTimeout"`
	ExpectContinue      time.Duration `mapstructure:"expectContinue"`
	WriteBufferSize     int           `mapstructure:"writeBufferSize"`
	ReadBufferSize      int           `mapstructure:"readBufferSize"`
}
