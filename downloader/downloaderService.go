package downloader

import (
	"fmt"
	"seed-sync/log"
	"sync"

	"go.uber.org/zap"
)

type downloaderService struct {
	downloaderDAO *DownloaderDAO
	downloaders   map[string]*Downloader
	lock          sync.Mutex
}

var DownloaderService *downloaderService

// 初始化的时候，从库里查出所有的下载器并初始化保存
func init() {
	DownloaderService = &downloaderService{
		downloaderDAO: downloaderDAO,
		downloaders:   make(map[string]*Downloader),
		lock:          sync.Mutex{},
	}
	downloaders, err := DownloaderService.downloaderDAO.GetAllDownloaders()
	if err != nil {
		log.Error("get all downloaders error", zap.Error(err))
		panic(err)
	}
	for _, downloader := range downloaders {
		DownloaderService.CreateDownloader(&DownloaderCreateRequest{
			Name:     downloader.Name,
			Url:      downloader.Url,
			Username: downloader.Username,
			Password: downloader.Password,
			Type:     downloader.Type,
		})
	}
}

// 创建下载器
func (service *downloaderService) CreateDownloader(request *DownloaderCreateRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[request.Name]; ok {
		return fmt.Errorf("downloader already exists: %s", request.Name)
	}
	tx := service.downloaderDAO.db.Begin()
	service.downloaderDAO.AddDownloaderTx(tx, &DownloaderTable{
		Name:     request.Name,
		Url:      request.Url,
		Username: request.Username,
		Password: request.Password,
		Type:     request.Type,
	})
	downloader, err := NewDownloader(&DownloaderConfig{
		Type:     request.Type,
		Url:      request.Url,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	service.downloaders[request.Name] = &downloader
	return nil
}

// 删除下载器
func (service *downloaderService) DeleteDownloader(name string) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[name]; !ok {
		return fmt.Errorf("downloader not found: %s", name)
	}
	service.downloaderDAO.DeleteDownloader(name)
	delete(service.downloaders, name)
	return nil
}

// 更新下载器
func (service *downloaderService) UpdateDownloader(request *DownloaderCreateRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[request.Name]; !ok {
		return fmt.Errorf("downloader not found: %s", request.Name)
	}
	//先删除再下载
	err := service.DeleteDownloader(request.Name)
	if err != nil {
		return err
	}
	return service.CreateDownloader(request)
}
