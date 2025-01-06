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
	systemParamDao *db.SystemParamDAO
}

var cookieCloudDAO *CookieCloudDAO

func InitCookieCloudDAO() {
	cookieCloudDAO = &CookieCloudDAO{
		db: db.DB,
		systemParamDao: db.SystemParamDao,
	}
}

func (dao *CookieCloudDAO) CreateCookieCloudConfig(config *CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return dao.systemParamDao.CreateSystemParam(CookieCloudConfigKey, string(configBytes))
}

func (dao *CookieCloudDAO) UpdateCookieCloudConfig(config *CookieCloudConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return dao.systemParamDao.UpdateSystemParam(CookieCloudConfigKey, string(configBytes))
}

func (dao *CookieCloudDAO) DeleteCookieCloudConfig() error {
	return dao.systemParamDao.DeleteSystemParam(CookieCloudConfigKey)
}

func (dao *CookieCloudDAO) GetCookieCloudConfig() *CookieCloudConfig {
	configBytes, err := dao.systemParamDao.GetSystemParam(CookieCloudConfigKey)
	if err != nil {
		return nil
	}
	config := &CookieCloudConfig{}
	err = json.Unmarshal([]byte(configBytes), config)
	if err != nil {
		return nil
	}
	return config
}
