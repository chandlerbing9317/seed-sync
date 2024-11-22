package cookieCloud

import (
	"encoding/json"
	"seed-sync/db"

	"gorm.io/gorm"
)

const (
	CookieCloudConfigKey = "cookie_cloud_config_key"
)

type CookieCloudDAO struct {
	db *gorm.DB
}

var cookieCloudDAO = &CookieCloudDAO{
	db: db.DB,
}



func (dao *CookieCloudDAO) CreateCookieCloudConfig(config *CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return db.SystemParamDao.CreateSystemParam(CookieCloudConfigKey, string(configBytes))
}

func (dao *CookieCloudDAO) UpdateCookieCloudConfig(config *CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return db.SystemParamDao.UpdateSystemParam(CookieCloudConfigKey, string(configBytes))
}

func (dao *CookieCloudDAO) DeleteCookieCloudConfig() error {
	return db.SystemParamDao.DeleteSystemParam(CookieCloudConfigKey)
}
