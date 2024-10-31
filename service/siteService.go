package service

import (
	"time"

	"gorm.io/gorm"
)

type SiteService struct {
	db *gorm.DB
}

// 站点表，存储的是用户提交的站点信息
type Site struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name" gorm:"uniqueIndex;not null;column:name"`
	Order          int       `json:"order" gorm:"index;not null;column:order"`
	URL            string    `json:"url" gorm:"column:url"`
	Cookie         string    `json:"cookie" gorm:"column:cookie"`
	APIKey         string    `json:"api_key" gorm:"column:api_key"`
	Token          string    `json:"token" gorm:"column:token"`
	CustomHeader   string    `json:"custom_header" gorm:"column:custom_header"`
	Passkey        string    `json:"passkey" gorm:"column:passkey"`
	RSS            string    `json:"rss" gorm:"column:rss"`
	Domains        string    `json:"domains" gorm:"column:domains"`
	DownloadURL    string    `json:"download_url" gorm:"column:download_url"`
	TorrentListURL string    `json:"torrent_list_url" gorm:"column:torrent_list_url"`
	Proxy          bool      `json:"proxy" gorm:"column:proxy"`
	Timeout        int       `json:"timeout" gorm:"column:timeout"`
	IsOverride     bool      `json:"is_override" gorm:"column:is_override"` // 是否使用服务端配置覆盖
	IsActive       bool      `json:"is_active" gorm:"column:is_active"`
	CreateTime     time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime     time.Time `json:"update_time" gorm:"column:update_time"`
}

// TableName 指定表名
func (Site) TableName() string {
	return "seed_sync_site"
}

func (service *SiteService) GetAllSites() ([]*Site, error) {
	var sites []*Site
	if err := service.db.Find(&sites).Error; err != nil {
		return nil, err
	}
	return sites, nil
}
