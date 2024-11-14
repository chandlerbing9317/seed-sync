package service

import (
	"seed-sync/driver/db"
)

type SiteService struct {
	siteDao *db.SiteDAO
}

var Site = &SiteService{
	siteDao: db.SiteDao,
}

func (s *SiteService) GetSiteList() ([]*db.Site, error) {
	return s.siteDao.GetAllSites()
}


