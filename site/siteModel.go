package site

import (
	"time"
)

type AddSiteRequest struct {
	SiteName      string   `json:"siteName"`
	ShowName      string   `json:"-"`
	Host          string   `json:"host"`
	Cookie        string   `json:"cookie"`
	ApiToken      string   `json:"apiToken"`
	Passkey       string   `json:"passkey"`
	RssKey        string   `json:"rssKey"`
	UserAgent     string   `json:"userAgent"`
	CustomHeaders []string `json:"customHeaders"`
	Proxy         bool     `json:"proxy"`
	// 流控配置
	MaxPerMin  int  `json:"maxPerMin"`
	MaxPerHour int  `json:"maxPerHour"`
	MaxPerDay  int  `json:"maxPerDay"`
	IsActive   bool `json:"isActive"`
	Timeout    int  `json:"timeout"`
}

type SiteInfo struct {
	// SiteTable 字段
	ID           int64  `json:"id"`
	SiteName     string `json:"siteName"`
	ShowName     string `json:"showName"`
	Order        int    `json:"order"`
	Host         string `json:"host"`
	Cookie       string `json:"cookie"`
	ApiToken     string `json:"apiToken"`
	Passkey      string `json:"passkey"`
	RssKey       string `json:"rssKey"`
	UserAgent    string `json:"userAgent"`
	CustomHeader string `json:"customHeader"`
	Proxy        bool   `json:"proxy"`
	Timeout      int    `json:"timeout"`
	IsActive     bool   `json:"isActive"`

	// SiteFlowControl 字段
	MaxPerMin  int `json:"maxPerMin"`
	MaxPerHour int `json:"maxPerHour"`
	MaxPerDay  int `json:"maxPerDay"`

	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
