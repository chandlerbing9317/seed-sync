package seedSyncServer

import (
	"math"
	"sync"
)

type seedSyncServerService struct {
	syncSeedServerDriver *ServerClient
}

var SeedSyncServerService *seedSyncServerService

// 初始化的时候，从库里查出所有的辅种任务并初始化保存
var once sync.Once

func InitSeedSyncServerService() {
	once.Do(func() {
		SeedSyncServerService = &seedSyncServerService{
			syncSeedServerDriver: InitSeedSyncServerDriver(),
		}
	})
}

func (service *seedSyncServerService) GetUserStatus(token string) (string, error) {
	return service.syncSeedServerDriver.GetUserStatus(token)
}

func (service *seedSyncServerService) GetSupportedSite(token string) ([]SupportSiteResponse, error) {
	supportedSites, err := service.syncSeedServerDriver.GetSupportedSites(0, math.MaxInt, token)
	if err != nil {
		return nil, err
	}
	return supportedSites, nil
}

// 辅种
func (service *seedSyncServerService) SeedSync(request *SeedSyncRequest, token string) (map[string][]SeedSyncTorrentInfoResponse, error) {
	return service.syncSeedServerDriver.SyncSeed(request, token)
}
