package seedSync

import "time"

type CreateSeedSyncTaskRequest struct {
	SeedSyncTaskBaseInfo
}

type UpdateSeedSyncTaskRequest struct {
	ID int64 `json:"id"`
	SeedSyncTaskBaseInfo
}

// 辅种任务基础信息
type SeedSyncTaskBaseInfo struct {
	TaskName string `json:"taskName"`
	//要辅种的站点
	SiteList []string `json:"siteList"`
	//要辅种的下载器id
	DownloaderId int64 `json:"downloaderId"`
	//辅种排除路径 下载器内的路径
	ExcludePath []string `json:"excludePath"`
	//辅种排除标签 下载器内的标签
	ExcludeTag []string `json:"excludeTag"`
	//辅种最小种子大小 
	MinSize int64 `json:"minSize"`
	//辅种后添加标签
	AddTag []string `json:"addTag"`
	//辅种任务状态，一般是开启或停止
	Status string `json:"status"`
	//辅种定时任务cron表达式
	Cron string `json:"cron"`
}

type SeedSyncTaskInfo struct {
	ID int64 `json:"id"`
	SeedSyncTaskBaseInfo
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
