package seedSyncServer

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
}

var SeedSyncServerClient *ServerClient = &ServerClient{
	ServerConfig: &config.Conf.ServerConfig,
}

const (
	GET_SUPPORTED_SITES_URL = "/sites/supported"
	SEED_SYNC_URL           = "/seed/sync"
	CHECK_USER_URL          = "/user/check"
)

// 获取支持的站点
func (s *ServerClient) GetSupportedSites(page int, pageSize int) (*common.PageResponse[SupportSiteResponse], error) {
	url := fmt.Sprintf("%s%s?page=%d&size=%d", s.ServerConfig.Url, GET_SUPPORTED_SITES_URL, page, pageSize)
	req, err := getHttpRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := common.DefaultHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result common.Result[common.PageResponse[SupportSiteResponse]]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Msg)
	}
	return &result.Data, nil
}

// 辅种 查询可辅种的种子
func (s *ServerClient) SyncSeed(request *SeedSyncRequest) (map[string][]SeedSyncTorrentInfoResponse, error) {
	url := fmt.Sprintf("%s%s", s.ServerConfig.Url, SEED_SYNC_URL)
	req, err := getHttpRequest("POST", url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := common.DefaultHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result common.Result[map[string][]SeedSyncTorrentInfoResponse]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Msg)
	}
	return result.Data, nil
}

// 检查用户
func (s *ServerClient) CheckUser() (bool, error) {
	url := fmt.Sprintf("%s%s", s.ServerConfig.Url, CHECK_USER_URL)
	req, err := getHttpRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	resp, err := common.DefaultHttpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result common.Result[any]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}
	if !result.Success {
		return false, nil
	}
	return true, nil
}

func getHttpHeader() (map[string]string, error) {
	header := make(map[string]string)
	header["X-Token"] = "123456"
	return header, nil
}

func getHttpRequest(method string, url string, body any) (*http.Request, error) {
	header, err := getHttpHeader()
	if err != nil {
		return nil, err
	}
	return common.GetRequest(method, url, header, body)
}
