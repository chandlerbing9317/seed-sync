package db

import (
	"time"

	"gorm.io/gorm"
)

type SystemParamDAO struct {
	db *gorm.DB
}

var SystemParamDao = &SystemParamDAO{
	db: DB,
}

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
	err := dao.db.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		return "", err
	}
	return systemParam.Value, nil
}

// 带事务的保存或更新系统参数
func (dao *SystemParamDAO) SaveOrUpdateSystemParamWithTx(tx *gorm.DB, key string, value string) error {
	return tx.Where("`key` = ?", key).
		Assign(map[string]interface{}{
			"value":       value,
			"update_time": time.Now(),
		}).
		FirstOrCreate(&SystemParam{
			Key:   key,
			Value: value,
		}).Error
}

// 带事务的删除系统参数
func (dao *SystemParamDAO) DeleteSystemParamWithTx(tx *gorm.DB, key string) error {
	return tx.Where("`key` = ?", key).Delete(&SystemParam{}).Error
}

func (dao *SystemParamDAO) SaveOrUpdateSystemParam(key, value string) error {
	var systemParam SystemParam
	err := dao.db.Model(&SystemParam{}).Where("key = ?", key).First(&systemParam).Error
	if err != nil {
		// 如果记录不存在，则创建新记录
		if err == gorm.ErrRecordNotFound {
			systemParam.Key = key
			systemParam.Value = value
			systemParam.CreateTime = time.Now()
			return dao.db.Create(&systemParam).Error
		}
		return err
	}
	//否则更新记录
	systemParam.Value = value
	systemParam.UpdateTime = time.Now()
	return dao.db.Save(&systemParam).Error
}

// 删除系统参数
func (dao *SystemParamDAO) DeleteSystemParam(key string) error {
	return dao.db.Model(&SystemParam{}).Where("key = ?", key).Delete(&SystemParam{}).Error
}
