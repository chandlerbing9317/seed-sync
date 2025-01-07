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
	Update(config *DownloaderConfig) error
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
	ID          int64
	InfoHash    string
	Size        int64
	Tags        []string
	DownloadDir string
	Status      string
}

const (
	SEED_STATUS_STOPPED           = "stopped"
	SEED_STATUS_VERIFY_LOCAL_DATA = "verify_local_data"
	SEED_STATUS_QUEUE_TO_DOWNLOAD = "queue_to_download"
	SEED_STATUS_DOWNLOADING       = "downloading"
	SEED_STATUS_QUEUE_TO_SEED     = "queue_to_seed"
	SEED_STATUS_SEEDING           = "seeding"
)

func NewDownloader(config *DownloaderConfig) (Downloader, error) {
	switch config.Type {
	case common.DOWNLOADER_TYPE_TRANSMISSION:
		return NewTransmissionClient(config)
	case common.DOWNLOADER_TYPE_QBITTORRENT:
		return NewQbittorrentClient(config)
	}
	return nil, fmt.Errorf("unsupported downloader type: %s", config.Type)
}
