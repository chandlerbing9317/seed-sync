package db

import (
	"time"
)

type SystemParamDAO struct {
}

var SystemParamDao = &SystemParamDAO{}

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

func (dao *SystemParamDAO) GetSystemParam(key string) (string, error) {
	var systemParam SystemParam
	err := DB.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		return "", err
	}
	return systemParam.Value, nil
}

// 保存
func (dao *SystemParamDAO) CreateSystemParam(key, value string) error {
	var systemParam SystemParam
	systemParam.Key = key
	systemParam.Value = value
	systemParam.CreateTime = time.Now()
	return DB.Create(&systemParam).Error
}

// 更新
func (dao *SystemParamDAO) UpdateSystemParam(key, value string) error {
	var systemParam SystemParam
	err := DB.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		return err
	}
	systemParam.Value = value
	systemParam.UpdateTime = time.Now()
	return DB.Save(&systemParam).Error
}

// 删除系统参数
func (dao *SystemParamDAO) DeleteSystemParam(key string) error {
	return DB.Model(&SystemParam{}).Where("key = ?", key).Delete(&SystemParam{}).Error
}
