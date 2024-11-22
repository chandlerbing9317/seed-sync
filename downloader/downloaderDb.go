package downloader

import (
	"seed-sync/db"

	"gorm.io/gorm"
)

type DownloaderDAO struct {
	db *gorm.DB
}

var downloaderDAO = &DownloaderDAO{
	db: db.DB,
}

type DownloaderTable struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func (DownloaderTable) TableName() string {
	return "seed_sync_downloader"
}

func (d *DownloaderDAO) AddDownloader(downloader *DownloaderTable) error {
	return d.db.Create(downloader).Error
}

// 事务版本的添加
func (d *DownloaderDAO) AddDownloaderTx(tx *gorm.DB, downloader *DownloaderTable) error {
	return tx.Create(downloader).Error
}

func (d *DownloaderDAO) UpdateDownloader(downloader *DownloaderTable) error {
	return d.db.Save(downloader).Error
}

// 事务版本的更新
func (d *DownloaderDAO) UpdateDownloaderTx(tx *gorm.DB, downloader *Downloader) error {
	return tx.Save(downloader).Error
}

func (d *DownloaderDAO) DeleteDownloader(name string) error {
	return d.db.Where("name = ?", name).Delete(&DownloaderTable{}).Error
}

// 事务版本的删除
func (d *DownloaderDAO) DeleteDownloaderTx(tx *gorm.DB, name string) error {
	return tx.Where("name = ?", name).Delete(&DownloaderTable{}).Error
}

func (d *DownloaderDAO) GetDownloader(name string) (*DownloaderTable, error) {
	var downloader DownloaderTable
	err := d.db.Where("name = ?", name).First(&downloader).Error
	if err != nil {
		return nil, err
	}
	return &downloader, nil
}

func (d *DownloaderDAO) GetAllDownloaders() ([]*DownloaderTable, error) {
	var downloaders []*DownloaderTable
	err := d.db.Find(&downloaders).Error
	if err != nil {
		return nil, err
	}
	return downloaders, nil
}