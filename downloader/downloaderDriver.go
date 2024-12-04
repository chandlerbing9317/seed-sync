package downloader

import (
	"fmt"
	"seed-sync/common"
)

const (
	DOWNLOADER_STATUS_AVAILABLE   = "available"
	DOWNLOADER_STATUS_UNAVAILABLE = "unavailable"
)

type Downloader interface {
	Type() string
	Ping() error
	GetSeedsHash() ([]SeedHash, error)
	AddTorrent(request *AddTorrentRequest) error
}

type AddTorrentRequest struct {
	TorrentUrl  string
	DownloadDir string
	TorrentFile []byte
	Paused      bool
}

type DownloaderConfig struct {
	Type     string
	Url      string
	Username string
	Password string
}

type SeedHash struct {
	InfoHash    string
	Size        int64
	Tags        []string
	DownloadDir string
}

func NewDownloader(config *DownloaderConfig) (Downloader, error) {
	switch config.Type {
	case common.DOWNLOADER_TYPE_TRANSMISSION:
		return NewTransmissionClient(config)
	case common.DOWNLOADER_TYPE_QBITTORRENT:
		return NewQbittorrentClient(config)
	}
	return nil, fmt.Errorf("unsupported downloader type: %s", config.Type)
}
