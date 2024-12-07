package downloader

import (
	"fmt"
	"sync"
)

type downloaderService struct {
	downloaderDAO *DownloaderDAO
	downloaders   map[string]Downloader
	lock          sync.Mutex
}

var DownloaderService *downloaderService

// 初始化的时候，从库里查出所有的下载器并初始化保存
func init() {
	DownloaderService = &downloaderService{
		downloaderDAO: downloaderDAO,
		downloaders:   make(map[string]Downloader),
		lock:          sync.Mutex{},
	}
	downloaderTables := DownloaderService.downloaderDAO.GetAllDownloaders()
	//遍历数据库的下载器并初始化保存
	for _, downloaderTable := range downloaderTables {
		downloader, err := NewDownloader(&DownloaderConfig{
			Type:     downloaderTable.Type,
			Url:      downloaderTable.Url,
			Username: downloaderTable.Username,
			Password: downloaderTable.Password,
		})
		if err != nil {
			panic(err)
		}
		DownloaderService.downloaders[downloaderTable.Name] = downloader
	}
}

// 创建下载器
func (service *downloaderService) CreateDownloader(request *DownloaderCreateRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[request.Name]; ok {
		return fmt.Errorf("下载器已存在: %s", request.Name)
	}
	tx := service.downloaderDAO.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := service.downloaderDAO.AddDownloaderTx(tx, &DownloaderTable{
		Name:     request.Name,
		Url:      request.Url,
		Username: request.Username,
		Password: request.Password,
		Type:     request.Type,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
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
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	service.downloaders[request.Name] = downloader
	return nil
}

// 删除下载器
func (service *downloaderService) DeleteDownloader(name string) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[name]; !ok {
		return fmt.Errorf("未找到下载器: %s", name)
	}
	service.downloaderDAO.DeleteDownloader(name)
	delete(service.downloaders, name)
	return nil
}

// 更新下载器
func (service *downloaderService) UpdateDownloader(request *DownloaderUpdateRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	if _, ok := service.downloaders[request.Name]; !ok {
		return fmt.Errorf("未找到下载器: %s", request.Name)
	}
	//开启事务
	tx := service.downloaderDAO.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := service.downloaderDAO.UpdateDownloaderTx(tx, &DownloaderTable{
		Id:       request.Id,
		Name:     request.Name,
		Url:      request.Url,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	err = service.downloaders[request.Name].Update(&DownloaderConfig{
		Type:     request.Type,
		Url:      request.Url,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (service *downloaderService) GetDownloaderList() []DownloaderTable {
	service.lock.Lock()
	defer service.lock.Unlock()
	return service.downloaderDAO.GetAllDownloaders()
}

// 获取下载器
func (service *downloaderService) GetDownloaderById(id int64) (Downloader, error) {
	service.lock.Lock()
	defer service.lock.Unlock()
	downloader := service.downloaderDAO.GetDownloaderById(id)
	if downloader == nil {
		return nil, fmt.Errorf("未找到下载器，下载器id: %d", id)
	}
	d, ok := service.downloaders[downloader.Name]
	if !ok {
		return nil, fmt.Errorf("未找到下载器，下载器id: %d", id)
	}
	return d, nil
}
