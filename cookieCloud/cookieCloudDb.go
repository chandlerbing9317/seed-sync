package cookieCloud

import (
	"encoding/json"
	"seed-sync/db"
)

const (
	CookieCloudConfigKey = "cookie_cloud_config_key"
)

type CookieCloudDAO struct {
}

var cookieCloudDAO = &CookieCloudDAO{}

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

func (dao *CookieCloudDAO) GetCookieCloudConfig() (*CookieCloudConfig, error) {
	configBytes, err := db.SystemParamDao.GetSystemParam(CookieCloudConfigKey)
	if err != nil {
		return nil, err
	}
	config := &CookieCloudConfig{}
	err = json.Unmarshal([]byte(configBytes), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
