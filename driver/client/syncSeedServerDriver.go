package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"seed-sync/common"
	"seed-sync/config"
	seedSyncModel "seed-sync/model/seed-sync"
)

type ServerClient struct {
	ServerConfig *config.ServerConfig
	client       *http.Client
}

var SeedSyncServerClient *ServerClient

const (
	GET_SUPPORTED_SITES_URL = "/sites/supported"
	SEED_SYNC_URL           = "/seed/sync"
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

// 获取支持的站点
func (s *ServerClient) GetSupportedSites(page int, pageSize int) (*common.PageResponse[seedSyncModel.SupportSiteResponse], error) {
	url := fmt.Sprintf("%s%s?page=%d&size=%d", s.ServerConfig.Url, GET_SUPPORTED_SITES_URL, page, pageSize)
	header, err := getHttpHeader()
	if err != nil {
		return nil, err
	}
	req, err := common.GetRequest("GET", url, header, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result common.Result[common.PageResponse[seedSyncModel.SupportSiteResponse]]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Msg)
	}
	return &result.Data, nil
}

func (s *ServerClient) SyncSeed(request *seedSyncModel.SeedSyncRequest) (*seedSyncModel.SeedSyncResponse, error) {
	url := fmt.Sprintf("%s%s", s.ServerConfig.Url, SEED_SYNC_URL)
	header, err := getHttpHeader()
	if err != nil {
		return nil, err
	}
	header["Content-Type"] = "application/json"
	req, err := common.GetRequest("POST", url, header, request)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result common.Result[seedSyncModel.SeedSyncResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Msg)
	}
	return &result.Data, nil
}

func getHttpHeader() (map[string]string, error) {
	header := make(map[string]string)
	header["X-Token"] = "123456"
	return header, nil
}
