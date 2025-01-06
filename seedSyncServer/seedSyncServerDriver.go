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

var seedSyncServerClient *ServerClient

func InitSeedSyncServerDriver() *ServerClient {
	seedSyncServerClient = &ServerClient{
		ServerConfig: &config.Conf.ServerConfig,
	}
	return seedSyncServerClient
}

const (
	GET_SUPPORTED_SITES_URL = "/sites/supported"
	SEED_SYNC_URL           = "/seed/sync"
	CHECK_USER_URL          = "/user/check"
	GET_USER_STATUS_URL     = "/user/status"
)

// 获取支持的站点
func (s *ServerClient) GetSupportedSites(page int, pageSize int, token string) ([]SupportSiteResponse, error) {
	url := fmt.Sprintf("%s%s?page=%d&size=%d", s.ServerConfig.Url, GET_SUPPORTED_SITES_URL, page, pageSize)
	req, err := getHttpRequest("GET", url, nil, token)
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
	return result.Data.Records, nil
}

// 辅种 查询可辅种的种子
func (s *ServerClient) SyncSeed(request *SeedSyncRequest, token string) (map[string][]SeedSyncTorrentInfoResponse, error) {
	url := fmt.Sprintf("%s%s", s.ServerConfig.Url, SEED_SYNC_URL)
	req, err := getHttpRequest("POST", url, request, token)
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

// 查询用户状态
func (s *ServerClient) GetUserStatus(token string) (string, error) {
	url := fmt.Sprintf("%s%s", s.ServerConfig.Url, GET_USER_STATUS_URL)
	req, err := getHttpRequest("GET", url, nil, token)
	if err != nil {
		return "", err
	}
	resp, err := common.DefaultHttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result common.Result[string]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", errors.New(result.Msg)
	}
	return result.Data, nil
}

func getHttpHeader(token string) (map[string]string, error) {
	header := make(map[string]string)
	header["X-Token"] = token
	return header, nil
}

func getHttpRequest(method string, url string, body any, token string) (*http.Request, error) {
	header, err := getHttpHeader(token)
	if err != nil {
		return nil, err
	}
	return common.GetRequest(method, url, header, body)
}
