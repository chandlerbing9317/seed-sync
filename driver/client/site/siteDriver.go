package site

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"seed-sync/common"
	"seed-sync/config"
	"seed-sync/driver/db"
	"strings"
	"time"
)

// 站点接口 所有站点都实现该接口
type Site interface {
	SiteName() string
	DownloadSeed(torrentId int) ([]byte, error)
	Ping() error
}

// 基础站点，对BaseSite接口的实现
type BaseSite struct {
	SiteInfo *db.Site
}

func NewBaseSite(siteInfo *db.Site) *BaseSite {
	return &BaseSite{
		SiteInfo: siteInfo,
	}
}

// 根据种子id下载种子文件
func (baseSite *BaseSite) DownloadSeed(torrentId int) ([]byte, error) {
	requestUrl := baseSite.GetDownloadUrl(torrentId)
	requestClient, err := baseSite.GetRequestClient()
	if err != nil {
		return nil, err
	}
	resp, err := requestClient.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 站点ping功能，检测站点是否可用
func (baseSite *BaseSite) Ping() error {
	requestUrl := baseSite.GetPingUrl()
	requestClient, err := baseSite.GetRequestClient()
	if err != nil {
		return err
	}
	resp, err := requestClient.Get(requestUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//todo: 处理cookie失效异常
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed, status code: %d", resp.StatusCode)
	}
	return nil
}

// 默认的http header
func (baseSite *BaseSite) GetHttpHeader() map[string]string {
	header := make(map[string]string)
	header["Cookie"] = baseSite.SiteInfo.Cookie
	header["User-Agent"] = baseSite.SiteInfo.UserAgent
	if baseSite.SiteInfo.ApiToken != "" {
		header["X-API-Token"] = baseSite.SiteInfo.ApiToken
	}
	//处理自定义header，自定义header的格式为：key1=value1${common.CustomHeaderSeparator}key2=value2
	//自定义header的优先级高于上述明确设置的header
	if baseSite.SiteInfo.CustomHeader != "" {
		for _, customHeader := range strings.Split(baseSite.SiteInfo.CustomHeader, common.CustomHeaderSeparator) {
			headerPair := strings.Split(customHeader, ":")
			if len(headerPair) == 2 {
				header[strings.TrimSpace(headerPair[0])] = strings.TrimSpace(headerPair[1])
			}
		}
	}
	return header
}

// 默认http_client
func (baseSite *BaseSite) GetRequestClient() (*http.Client, error) {
	var transport *http.Transport

	//处理代理
	if baseSite.SiteInfo.Proxy {
		proxyURL, err := url.Parse(config.Conf.ProxyConfig.ProxyURL)
		if err != nil {
			return nil, err
		}

		if config.Conf.ProxyConfig.ProxyUsername != "" {
			proxyURL.User = url.UserPassword(
				config.Conf.ProxyConfig.ProxyUsername,
				config.Conf.ProxyConfig.ProxyPassword,
			)
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	} else {
		transport = &http.Transport{}
	}

	//处理header
	headers := baseSite.GetHttpHeader()
	customTransport := &roundTripperWithHeaders{
		headers:  headers,
		original: transport,
	}

	return &http.Client{
		Transport: customTransport,
		Timeout:   time.Duration(baseSite.SiteInfo.Timeout) * time.Second,
	}, nil
}

// 父类接口
func (baseSite *BaseSite) SiteName() string {
	return ""
}

// 父类接口
func (baseSite *BaseSite) GetDownloadUrl(torrentId int) string {
	return ""
}

// 父类接口
func (baseSite *BaseSite) GetPingUrl() string {
	return ""
}

// 定义一个自定义的 RoundTripper
type roundTripperWithHeaders struct {
	headers  map[string]string
	original http.RoundTripper
}

func (rth *roundTripperWithHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	// 添加默认请求头
	for key, value := range rth.headers {
		req.Header.Set(key, value)
	}
	return rth.original.RoundTrip(req)
}
