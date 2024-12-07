package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"seed-sync/common"
	"strings"
)

//qbittorrent交互层
//ref:
//https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)

const (
	AUTH_URL    = "/api/v2/auth"
	APP_URL     = "/api/v2/app"
	TORRENT_URL = "/api/v2/torrents"
)

type QbittorrentClient struct {
	config  *DownloaderConfig
	version string
	cookie  string
}

type QbittorrentTorrentSeedHash struct {
	InfoHash    string   `json:"hash"`
	DownloadDir string   `json:"save_path"`
	Tags        []string `json:"tags"`
	Size        int64    `json:"size"`
}

func NewQbittorrentClient(config *DownloaderConfig) (*QbittorrentClient, error) {
	client := &QbittorrentClient{
		config:  config,
		version: "",
	}
	return client, nil
}

func (q *QbittorrentClient) Update(config *DownloaderConfig) error {
	q.config = config
	return nil
}

func (q *QbittorrentClient) Type() string {
	return common.DOWNLOADER_TYPE_QBITTORRENT
}

func (q *QbittorrentClient) Ping() error {
	_, err := q.getVersion()
	return err
}

func (q *QbittorrentClient) getVersion() (string, error) {
	//GET /api/v2/app/version
	body, err := q.doRequest("GET", APP_URL+"/version", nil)
	if err != nil {
		return "", err
	}
	q.version = string(body)
	return q.version, nil
}
func (q *QbittorrentClient) GetSeedsHash() ([]SeedHash, error) {
	//GET /api/v2/torrents/info
	body, err := q.doRequest("GET", TORRENT_URL+"/info", nil)
	if err != nil {
		return nil, err
	}
	//json 反序列化
	var torrents []QbittorrentTorrentSeedHash
	err = json.Unmarshal(body, &torrents)
	if err != nil {
		return nil, err
	}
	//转换为SeedHash
	seeds := make([]SeedHash, len(torrents))
	for i, torrent := range torrents {
		seeds[i] = SeedHash{
			InfoHash:    torrent.InfoHash,
			Size:        torrent.Size,
			Tags:        torrent.Tags,
			DownloadDir: torrent.DownloadDir,
		}
	}
	return seeds, nil
}

func (q *QbittorrentClient) AddTorrent(addTorrentRequest *AddTorrentRequest) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if len(addTorrentRequest.TorrentFile) > 0 {
		// 优先使用种子文件
		part, err := writer.CreateFormFile("torrents", "file.torrent")
		if err != nil {
			return fmt.Errorf("create form file failed: %w", err)
		}
		if _, err := part.Write(addTorrentRequest.TorrentFile); err != nil {
			return fmt.Errorf("write torrent file failed: %w", err)
		}
	} else if addTorrentRequest.TorrentUrl != "" {
		// 使用URL
		if err := writer.WriteField("urls", addTorrentRequest.TorrentUrl); err != nil {
			return fmt.Errorf("write urls field failed: %w", err)
		}
	} else {
		return fmt.Errorf("neither torrent file nor url provided")
	}

	// 设置下载目录
	if addTorrentRequest.DownloadDir != "" {
		if err := writer.WriteField("savepath", addTorrentRequest.DownloadDir); err != nil {
			return fmt.Errorf("write savepath field failed: %w", err)
		}
	}

	// 设置暂停状态
	if addTorrentRequest.Paused {
		if err := writer.WriteField("paused", "true"); err != nil {
			return fmt.Errorf("write paused field failed: %w", err)
		}
	}
	writer.Close()

	request, err := http.NewRequest("POST", q.config.Url+TORRENT_URL+"/add", body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	responseBody, err := q.doRequestWithCookie(request)
	if err != nil {
		return fmt.Errorf("add torrent request failed: %w", err)
	}

	if len(responseBody) > 0 {
		return fmt.Errorf("add torrent failed: %s", string(responseBody))
	}
	return nil
}

// 登录
func (q *QbittorrentClient) login() error {
	//curl -i --header 'Referer: http://localhost:8080' --data 'username=admin&password=adminadmin' http://localhost:8080/api/v2/auth/login
	form := url.Values{}
	form.Set("username", q.config.Username)
	form.Set("password", q.config.Password)
	request, err := http.NewRequest(
		"POST",
		q.config.Url+AUTH_URL+"/login",
		strings.NewReader(form.Encode()),
	)
	request.Header.Set("Referer", q.config.Url)
	if err != nil {
		return err
	}
	//发起请求
	client := common.DefaultHttpClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	//如果响应码是200 获取cookie
	if response.StatusCode == 200 {
		q.cookie = response.Header.Get("Set-Cookie")
	}
	return nil
}

func (q *QbittorrentClient) doRequest(method string, requestUrl string, params map[string]any) ([]byte, error) {
	var form url.Values
	if params != nil {
		form = url.Values{}
		for key, value := range params {
			form.Set(key, fmt.Sprintf("%v", value))
		}
	}
	request, err := http.NewRequest(method, q.config.Url+requestUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	return q.doRequestWithCookie(request)
}

// 统一请求处。设置cookie，如果cookie失效就重新登录
func (q *QbittorrentClient) doRequestWithCookie(request *http.Request) ([]byte, error) {
	request.Header.Set("Cookie", q.cookie)
	client := common.DefaultHttpClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	//如果响应是403 重新登录
	if response.StatusCode == 403 {
		err = q.login()
		if err != nil {
			return nil, err
		}
		return q.doRequestWithCookie(request)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
