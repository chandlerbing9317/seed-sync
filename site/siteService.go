package site

type siteService struct {
	siteDao *SiteDAO
}

var SiteService = &siteService{
	siteDao: siteDAO,
}

func (service *siteService) GetSiteList() ([]*SiteTable, error) {
	return service.siteDao.GetAllSites()
}

func (service *siteService) UpdateSite(site *SiteTable) error {
	return service.siteDao.UpdateSite(site)
}

func (service *siteService) UpdateBatchSite(sites []*SiteTable) error {
	return service.siteDao.UpdateBatchSite(sites)
}
