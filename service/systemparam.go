package service

import (
	"seed-sync/driver"
	"time"

	"gorm.io/gorm"
)

type SystemParam struct {
	ID         int       `gorm:"column:id"`
	Key        string    `gorm:"column:key"`
	Value      string    `gorm:"column:value"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (s *SystemParam) TableName() string {
	return "seed_sync_system_param"
}

func GetSystemParam(key string) (string, error) {
	var systemParam SystemParam
	err := driver.DB.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		return "", err
	}
	return systemParam.Value, nil
}

func SaveOrUpdateSystemParam(key, value string) error {
	var systemParam SystemParam
	err := driver.DB.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		// 如果记录不存在，则创建新记录
		if err == gorm.ErrRecordNotFound {
			systemParam.Key = key
			systemParam.Value = value
			systemParam.CreateTime = time.Now()
			return driver.DB.Create(&systemParam).Error
		}
		return err
	}
	//否则更新记录
	systemParam.Value = value
	systemParam.UpdateTime = time.Now()
	return driver.DB.Save(&systemParam).Error
}
