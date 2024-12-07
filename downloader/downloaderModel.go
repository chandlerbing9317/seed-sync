package downloader

type DownloaderCreateRequest struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type DownloaderUpdateRequest struct {
	DownloaderCreateRequest
	Id int `json:"id"`
}
