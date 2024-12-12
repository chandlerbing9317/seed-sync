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
	Cookie       string    `json:"cookie" gorm:"column:cookie"`
	ApiToken     string    `json:"apiToken" gorm:"column:api_token"`
	Passkey      string    `json:"passkey" gorm:"column:passkey"`
	RssKey       string    `json:"rssKey" gorm:"column:rss_key"`
	UserAgent    string    `json:"userAgent" gorm:"column:user_agent"`
	CustomHeader string    `json:"customHeader" gorm:"column:custom_header"`
	Proxy        bool      `json:"proxy" gorm:"column:proxy"`
	IsActive     bool      `json:"isActive" gorm:"column:is_active"` // 是否启用
	Timeout      int       `json:"timeout" gorm:"column:timeout"`    //单位秒
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`
}

// TableName 指定表名
func (SiteTable) TableName() string {
	return "seed_sync_site"
}

// 站点流控表
type SiteFlowControl struct {
	ID         int64     `json:"-"`
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

func (dao *SiteDAO) AddSiteTx(tx *gorm.DB, site *SiteTable, siteFlowControl *SiteFlowControl) error {
	if err := tx.Create(site).Error; err != nil {
		return err
	}
	if err := tx.Create(siteFlowControl).Error; err != nil {
		return err
	}
	return nil
}

// 更新一个站点
func (dao *SiteDAO) UpdateSiteBySiteNameTx(tx *gorm.DB, siteName string, site *SiteTable, siteFlowControl *SiteFlowControl) error {
	if err := tx.Model(&SiteTable{}).Where("site_name = ?", siteName).Updates(site).Error; err != nil {
		return err
	}
	if err := tx.Model(&SiteFlowControl{}).Where("site_name = ?", siteName).Updates(siteFlowControl).Error; err != nil {
		return err
	}
	return nil
}

// 删除站点
func (dao *SiteDAO) DeleteSiteBySiteNameTx(tx *gorm.DB, siteName string) error {
	if err := tx.Where("site_name = ?", siteName).Delete(&SiteTable{}).Error; err != nil {
		return err
	}
	return tx.Where("site_name = ?", siteName).Delete(&SiteFlowControl{}).Error
}

// GetMaxOrderTx 在事务中获取最大order
func (dao *SiteDAO) GetMaxOrderTx(tx *gorm.DB) (int, error) {
	var maxOrder struct {
		MaxOrder int
	}

	err := tx.Model(&SiteTable{}).
		Select("COALESCE(MAX(`order`), 0) as max_order").
		Scan(&maxOrder).Error

	return maxOrder.MaxOrder, err
}

// 获取站点详情
func (dao *SiteDAO) GetSiteInfo(siteName string) *SiteInfo {
	var siteTable SiteTable
	if err := dao.db.Where("site_name = ?", siteName).First(&siteTable).Error; err != nil {
		return nil
	}
	var siteFlowControl SiteFlowControl
	if err := dao.db.Where("site_name = ?", siteName).First(&siteFlowControl).Error; err != nil {
		return nil
	}
	return GenerateSiteInfo(&siteTable, &siteFlowControl)
}

func (dao *SiteDAO) GetAllSites() ([]*SiteInfo, error) {
	var siteTables []*SiteTable
	if err := dao.db.Order("`order` ASC").Find(&siteTables).Error; err != nil {
		return nil, err
	}
	var siteInfos []*SiteInfo
	for _, siteTable := range siteTables {
		var siteFlowControl SiteFlowControl
		if err := dao.db.Where("site_name = ?", siteTable.SiteName).First(&siteFlowControl).Error; err != nil {
			return nil, err
		}
		siteInfos = append(siteInfos, GenerateSiteInfo(siteTable, &siteFlowControl))
	}
	return siteInfos, nil
}

func (dao *SiteDAO) BatchUpdateSiteOrders(updates []SiteOrderUpdateRequest) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			if err := tx.Model(&SiteTable{}).
				Where("site_name = ?", update.SiteName).
				Update("order", update.NewOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *SiteDAO) UpdateCookieAndHost(siteName string, cookie string, host string) error {
	return dao.db.Model(&SiteTable{}).Where("site_name = ?", siteName).Updates(map[string]interface{}{
		"cookie": cookie,
		"host":   host,
	}).Error
}
