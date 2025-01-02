package seedSync

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

type SeedSyncDAO struct {
	DB *gorm.DB
}

var seedSyncDAO = &SeedSyncDAO{
	DB: db.DB,
}

type SeedSyncTaskTable struct {
	ID           int64     `gorm:"column:id;"`
	TaskName     string    `gorm:"column:task_name"`
	SiteList     string    `gorm:"column:site_list"`
	DownloaderId int64     `gorm:"column:downloader_id"`
	ExcludePath  string    `gorm:"column:exclude_path"`
	ExcludeTag   string    `gorm:"column:exclude_tag"`
	MinSize      int64     `gorm:"column:min_size"`
	AddTag       string    `gorm:"column:add_tag"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (SeedSyncTaskTable) TableName() string {
	return "seed_sync_seed_task"
}

func (s *SeedSyncDAO) CreateSeedSyncTask(task *SeedSyncTaskTable) error {
	return s.DB.Create(task).Error
}

func (s *SeedSyncDAO) CreateSeedSyncTaskWithTx(tx *gorm.DB, task *SeedSyncTaskTable) error {
	return tx.Create(task).Error
}

func (s *SeedSyncDAO) UpdateSeedSyncTask(task *SeedSyncTaskTable) error {
	return s.DB.Model(task).Where("id = ?", task.ID).Updates(task).Error
}

func (s *SeedSyncDAO) UpdateSeedSyncTaskWithTx(tx *gorm.DB, task *SeedSyncTaskTable) error {
	return tx.Model(task).Where("id = ?", task.ID).Updates(task).Error
}

func (s *SeedSyncDAO) GetSeedSyncTask(id int64) *SeedSyncTaskTable {
	var task SeedSyncTaskTable
	err := s.DB.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil
	}
	return &task
}

func (s *SeedSyncDAO) GetSeedSyncTaskByTaskName(taskName string) *SeedSyncTaskTable {
	var task SeedSyncTaskTable
	err := s.DB.Where("task_name = ?", taskName).First(&task).Error
	if err != nil {
		return nil
	}
	return &task
}

func (s *SeedSyncDAO) GetAllSeedSyncTaskList() []*SeedSyncTaskTable {
	var tasks []*SeedSyncTaskTable
	err := s.DB.Find(&tasks).Error
	if err != nil {
		return nil
	}
	return tasks
}
