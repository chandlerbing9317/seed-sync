package downloader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"seed-sync/common"
	"time"
)

// transmission交互层
// ref:
// https://github.com/transmission/transmission/blob/4.0.5/docs/rpc-spec.md
// https://github.com/hekmon/transmissionrpc

const (
	csrfHeader          = "X-Transmission-Session-Id"
	authorizationPrefix = "Basic "
)

// method
const (
	MethodTorrentStart      = "torrent-start"
	MethodTorrentStartNow   = "torrent-start-now"
	MethodTorrentStop       = "torrent-stop"
	MethodTorrentVerify     = "torrent-verify"
	MethodTorrentReannounce = "torrent-reannounce"
	MethodTorrentSet        = "torrent-set"
	MethodTorrentGet        = "torrent-get"
	MethodTorrentAdd        = "torrent-add"
	MethodSessionGet        = "session-get"
	MethodSessionSet        = "session-set"
	MethodSessionStats      = "session-stats"
)

type requestBody struct {
	Method    string         `json:"method"`
	Arguments map[string]any `json:"arguments"`
	Tag       int            `json:"tag,omitempty"`
}

type responseBody struct {
	Arguments map[string]any `json:"arguments"`
	Result    string         `json:"result"`
	Tag       *int           `json:"tag"`
}

type TransmissionClient struct {
	config       *DownloaderConfig
	version      string
	tagGenerator *rand.Rand
	csrfToken    string
}

type TransmissionTorrentSeedHash struct {
	InfoHash    string   `json:"hashString"`
	DownloadDir string   `json:"downloadDir"`
	TotalSize   int64    `json:"totalSize"`
	Labels      []string `json:"labels"`
}

func NewTransmissionClient(config *DownloaderConfig) (*TransmissionClient, error) {
	client := &TransmissionClient{
		config:       config,
		tagGenerator: rand.New(common.NewLockedRandomSource(time.Now().Unix())),
		csrfToken:    "",
	}
	return client, nil
}
func (t *TransmissionClient) Update(config *DownloaderConfig) error {
	t.config = config
	return nil
}

func (t *TransmissionClient) Type() string {
	return common.DOWNLOADER_TYPE_TRANSMISSION
}

// 使用session-get 方法测试连接是否可连接
func (t *TransmissionClient) Ping() error {
	arguments := map[string]any{
		"fields": []string{"rpc-version", "rpc-version-minimum", "rpc-version-semver", "version"},
	}
	response, err := t.doRequest(MethodSessionGet, arguments)
	if err != nil {
		return err
	}
	//处理response
	if response.Result != "success" {
		return fmt.Errorf("session-get failed: %s", response.Result)
	}
	//获取版本
	t.version = response.Arguments["version"].(string)
	return nil
}

func (t *TransmissionClient) GetSeedsHash() ([]SeedHash, error) {
	arguments := map[string]any{
		"fields": []string{"id", "hashString", "downloadDir", "totalSize", "labels"},
	}
	response, err := t.doRequest(MethodTorrentGet, arguments)
	if err != nil {
		return nil, err
	}
	//处理response
	hashes := make([]SeedHash, 0)
	for _, arg := range response.Arguments {
		hash := arg.(TransmissionTorrentSeedHash)
		hashes = append(hashes, SeedHash{
			InfoHash:    hash.InfoHash,
			Size:        hash.TotalSize,
			Tags:        hash.Labels,
			DownloadDir: hash.DownloadDir,
		})
	}
	return hashes, nil
}

func (t *TransmissionClient) AddTorrent(AddTorrentRequest *AddTorrentRequest) error {
	arguments := map[string]any{
		"downloadDir": AddTorrentRequest.DownloadDir,
		"paused":      AddTorrentRequest.Paused,
	}
	//优先使用torrentFile
	if AddTorrentRequest.TorrentFile != nil {
		arguments["metainfo"] = base64.StdEncoding.EncodeToString(AddTorrentRequest.TorrentFile)
	} else if AddTorrentRequest.TorrentUrl != "" {
		arguments["filename"] = AddTorrentRequest.TorrentUrl
	}

	response, err := t.doRequest(MethodTorrentAdd, arguments)
	if err != nil {
		return err
	}
	if response.Result != "success" {
		return fmt.Errorf("torrent-add failed: %s", response.Result)
	}
	return nil
}

func (t *TransmissionClient) getTag() int {
	return t.tagGenerator.Int()
}
func (t *TransmissionClient) getRequestHeader() map[string]string {
	header := map[string]string{}
	header["Content-Type"] = "application/json"
	header[csrfHeader] = t.csrfToken
	return header
}

// 发送请求
func (t *TransmissionClient) doRequest(method string, arguments map[string]any) (*responseBody, error) {
	request := requestBody{
		Method:    method,
		Arguments: arguments,
		Tag:       t.getTag(),
	}
	req, err := common.GetRequest("POST", t.config.Url+"/transmission/rpc", t.getRequestHeader(), request)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.config.Username, t.config.Password)
	//todo 这里可能考虑使用可复用的http_client
	client := common.DefaultHttpClient
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//如果响应码是409，则需要更新csrfToken 然后重新请求
	if httpResp.StatusCode == 409 {
		t.csrfToken = httpResp.Header.Get(csrfHeader)
		return t.doRequest(method, arguments)
	} else {
		return t.parseResponse(httpResp)
	}
}

// 解析response
func (t *TransmissionClient) parseResponse(httpResp *http.Response) (*responseBody, error) {
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	var response responseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
