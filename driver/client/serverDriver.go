package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"seed-sync/common"
	"seed-sync/config"
)

type ServerClient struct {
	ServerConfig *config.ServerConfig
	client       *http.Client
}

var SeedSyncServerClient *ServerClient

const (
	GET_SUPPORTED_SITES_URL = "/api/v1/sites/supported"
)

// 与seed-sync-server通信的客户端
// 需要建立可复用的http client 避免服务端网络开销过大
func init() {
	SeedSyncServerClient, _ = newServerClient()
}

func newServerClient() (*ServerClient, error) {
	return &ServerClient{
		ServerConfig: &config.Conf.ServerConfig,
		client:       common.NewHttpClient(config.Conf.ServerClientConfig),
	}, nil
}

type SupportSite struct {
	SiteName     string   `json:"siteName"`
	ShowName     string   `json:"showName"`
	Host         string   `json:"host"`
	Domain       []string `json:"domain"`
	UserAgent    string   `json:"userAgent"`
	CustomHeader []string `json:"customHeader"`
	SeedListUrl  string   `json:"seedListUrl"`
	RssUrl       string   `json:"rssUrl"`
	DetailUrl    string   `json:"detailUrl"`
	DownloadUrl  string   `json:"downloadUrl"`
	PingUrl      string   `json:"pingUrl"`
}

// 获取支持的站点
func (s *ServerClient) GetSupportedSites(page int, pageSize int) (*common.PageResponse[[]SupportSite], error) {
	url := fmt.Sprintf("%s%s?page=%d&size=%d", s.ServerConfig.Url, GET_SUPPORTED_SITES_URL, page, pageSize)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result common.Result[common.PageResponse[[]SupportSite]]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Msg)
	}
	return &result.Data, nil
}
