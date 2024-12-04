package seedSync

type CreateSeedSyncTaskRequest struct {
	TaskName     string   `json:"taskName"`
	SiteList     []string `json:"siteList"`
	DownloaderId int64    `json:"downloaderId"`
	ExcludePath  []string `json:"excludePath"`
	ExcludeTag   []string `json:"excludeTag"`
	MinSize      int64    `json:"minSize"`
	AddTag       string   `json:"addTag"`
	Status       string   `json:"status"`
}

type UpdateSeedSyncTaskRequest struct {
	Id           int64    `json:"id"`
	TaskName     string   `json:"taskName"`
	SiteList     []string `json:"siteList"`
	DownloaderId int64    `json:"downloaderId"`
	ExcludePath  []string `json:"excludePath"`
	ExcludeTag   []string `json:"excludeTag"`
	MinSize      int64    `json:"minSize"`
	AddTag       string   `json:"addTag"`
	Status       string   `json:"status"`
}
