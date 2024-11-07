package db

import (
	"encoding/json"

	"gorm.io/gorm"
)

const (
	CookieCloudConfigKey = "cookie_cloud_config_key"
)

type CookieCloudDAO struct {
	db *gorm.DB
}

var CookieCloudDao = &CookieCloudDAO{
	db: DB,
}

type CookieCloudConfig struct {
	Url string `json:"url"`
	//用户KEY
	UserKey string `json:"user_key"`
	//端对端加密密码
	P2pPassword string `json:"p2p_password"`
	//同步cron表达式
	SyncCron string `json:"sync_cron"`
}

func (dao *CookieCloudDAO) AddOrUpdateCookieCloudConfigWithTx(tx *gorm.DB, config *CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return SystemParamDao.SaveOrUpdateSystemParamWithTx(tx, CookieCloudConfigKey, string(configBytes))
}

func (dao *CookieCloudDAO) DeleteCookieCloudConfigWithTx(tx *gorm.DB) error {
	return SystemParamDao.DeleteSystemParamWithTx(tx, CookieCloudConfigKey)
}
