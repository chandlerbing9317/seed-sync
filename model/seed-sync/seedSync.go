package seedSyncModel

type SeedSyncRequest struct {
	//要辅种的种子
	Torrents []TorrentForSeedSyncRequest `json:"torrents"`
	//辅种的站点
	Sites []string `json:"sites"`
}
type TorrentForSeedSyncRequest struct {
	PiecesHash string `json:"piecesHash"`
	FilesHash  string `json:"filesHash"`
}

type SeedSyncResponse struct {
	Torrents []*TorrentForSeedSyncResponse `json:"torrents"`
}
type TorrentForSeedSyncResponse struct {
	SiteName  string `json:"siteName"`
	TorrentId int    `json:"torrentId"`
	InfoHash  string `json:"infoHash"`
}
