package common

import (
	"net"
	"net/http"
	"seed-sync/config"
	"time"
)

const (
	DefaultConnTimeout         = 30 * time.Second
	DefaultReadTimeout         = 30 * time.Second
	DefaultWriteTimeout        = 30 * time.Second
	DefaultKeepAlive           = 30 * time.Second
	DefaultIdleConnTimeout     = 2 * time.Minute
	DefaultTLSHandshakeTimeout = 30 * time.Second
	DefaultExpectContinue      = 1 * time.Second
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	DefaultWriteBufferSize     = 64 * 1024
	DefaultReadBufferSize      = 64 * 1024
)

// 按配置创建http client
func NewHttpClient(config config.HttpClientConfig) *http.Client {
	config = getHttpClientConfig(config)

	transport := &http.Transport{
		Proxy: GetProxyFunc(config.Proxy),
		DialContext: (&net.Dialer{
			Timeout:   config.ConnTimeout * time.Second, // 连接建立超时
			KeepAlive: config.KeepAlive * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          config.MaxIdleConns,                  // 最大空闲连接数
		IdleConnTimeout:       config.IdleConnTimeout * time.Second, // 空闲连接超时
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: config.ExpectContinue * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost, // 每个host最大空闲连接数

		ResponseHeaderTimeout: config.ReadTimeout * time.Second, // 等待服务器响应头的超时时间
		WriteBufferSize:       config.WriteBufferSize,
		ReadBufferSize:        config.ReadBufferSize,
	}
	return &http.Client{
		Transport: transport,
	}
}

// 根据已有config 填充默认值
func getHttpClientConfig(config config.HttpClientConfig) config.HttpClientConfig {
	if config.ConnTimeout == 0 {
		config.ConnTimeout = DefaultConnTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = DefaultReadTimeout
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = DefaultWriteTimeout
	}
	if config.KeepAlive == 0 {
		config.KeepAlive = DefaultKeepAlive
	}
	if config.IdleConnTimeout == 0 {
		config.IdleConnTimeout = DefaultIdleConnTimeout
	}
	if config.TLSHandshakeTimeout == 0 {
		config.TLSHandshakeTimeout = DefaultTLSHandshakeTimeout
	}
	if config.ExpectContinue == 0 {
		config.ExpectContinue = DefaultExpectContinue
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = DefaultMaxIdleConns
	}
	if config.MaxIdleConnsPerHost == 0 {
		config.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}
	if config.WriteBufferSize == 0 {
		config.WriteBufferSize = DefaultWriteBufferSize
	}
	if config.ReadBufferSize == 0 {
		config.ReadBufferSize = DefaultReadBufferSize
	}
	return config
}
