package site

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

// ===== 站点数据库相关 =====

type SiteDAO struct {
	db *gorm.DB
}

var siteDAO = &SiteDAO{
	db: db.DB,
}

// 站点表，存储的是用户提交的站点信息
type SiteTable struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SiteName     string    `json:"siteName" gorm:"uniqueIndex;not null;column:site_name"`
	ShowName     string    `json:"showName" gorm:"column:show_name"`
	Order        int       `json:"order" gorm:"index;not null;column:order"`
	Host         string    `json:"host" gorm:"column:host"`
	Domains      string    `json:"domains" gorm:"column:domains"`
	Cookie       string    `json:"cookie" gorm:"column:cookie"`
	ApiToken     string    `json:"apiToken" gorm:"column:api_token"`
	Passkey      string    `json:"passkey" gorm:"column:passkey"`
	RssKey       string    `json:"rssKey" gorm:"column:rss_key"`
	UserAgent    string    `json:"userAgent" gorm:"column:user_agent"`
	CustomHeader string    `json:"customHeader" gorm:"column:custom_header"`
	SeedListUrl  string    `json:"seedListUrl" gorm:"column:seed_list_url"`
	RssUrl       string    `json:"rssUrl" gorm:"column:rss_url"`
	DetailUrl    string    `json:"detailUrl" gorm:"column:detail_url"`
	DownloadUrl  string    `json:"downloadUrl" gorm:"column:download_url"`
	PingUrl      string    `json:"pingUrl" gorm:"column:ping_url"`
	Proxy        bool      `json:"proxy" gorm:"column:proxy"`
	Timeout      int       `json:"timeout" gorm:"column:timeout"`        // 单位：秒
	IsOverride   bool      `json:"isOverride" gorm:"column:is_override"` // 是否使用服务端配置覆盖
	IsActive     bool      `json:"isActive" gorm:"column:is_active"`     // 是否启用
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`
}

// TableName 指定表名
func (SiteTable) TableName() string {
	return "seed_sync_site"
}

// 站点流控表
type SiteFlowControl struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SiteName   string    `json:"site_name" gorm:"column:site_name"`
	MaxPerMin  int       `json:"max_per_min" gorm:"column:max_per_min"`
	MaxPerHour int       `json:"max_per_hour" gorm:"column:max_per_hour"`
	MaxPerDay  int       `json:"max_per_day" gorm:"column:max_per_day"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func DefaultSiteFlowControl() *SiteFlowControl {
	return &SiteFlowControl{
		MaxPerMin:  8,
		MaxPerHour: 40,
		MaxPerDay:  400,
	}
}

// TableName 指定表名
func (SiteFlowControl) TableName() string {
	return "seed_sync_site_flow_control"
}

// 获取所有站点
func (dao *SiteDAO) GetAllSites() ([]*SiteTable, error) {
	var sites []*SiteTable
	if err := dao.db.Find(&sites).Error; err != nil {
		return nil, err
	}
	return sites, nil
}

// 根据站点名查询站点流控
func (dao *SiteDAO) GetSiteFlowControl(siteName string) (*SiteFlowControl, error) {
	var siteFlowControl SiteFlowControl
	if err := dao.db.Where("site_name = ?", siteName).First(&siteFlowControl).Error; err != nil {
		return nil, err
	}
	return &siteFlowControl, nil
}

func (dao *SiteDAO) UpdateSite(site *SiteTable) error {
	return dao.db.Save(site).Error
}

func (dao *SiteDAO) UpdateBatchSite(sites []*SiteTable) error {
	return dao.db.Save(sites).Error
}

// 创建/更新一个站点
func (dao *SiteDAO) CreateOrUpdateSiteWithTx(tx *gorm.DB, site *SiteTable, siteFlowControl *SiteFlowControl) error {
	if siteFlowControl == nil {
		siteFlowControl = DefaultSiteFlowControl()
	}
	siteFlowControl.SiteName = site.SiteName

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 使用 name 字段查找并更新站点
	result := tx.Where("name = ?", site.SiteName).Updates(site)
	if result.RowsAffected == 0 {
		// 如果没有更新到记录，说明是新站点，执行创建
		if err := tx.Create(site).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// 更新站点流控
	result = tx.Where("site_name = ?", siteFlowControl.SiteName).Updates(siteFlowControl)
	if result.RowsAffected == 0 {
		if err := tx.Create(siteFlowControl).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}
