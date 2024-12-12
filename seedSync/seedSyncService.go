package seedSync

import (
	"errors"
	"fmt"
	"seed-sync/common"
	"seed-sync/downloader"
	"seed-sync/log"
	"seed-sync/seedSyncServer"
	"seed-sync/site"
	"strings"
	"sync"
	"time"
)

const SEED_SYNC_BATCH_SIZE = 200

type seedSyncService struct {
	seedSyncDAO *SeedSyncDAO
	lock        sync.Mutex
}

var SeedSyncService = &seedSyncService{
	seedSyncDAO: seedSyncDAO,
	lock:        sync.Mutex{},
}

// 创建辅种任务
func (service *seedSyncService) CreateSeedSyncTask(request *CreateSeedSyncTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	err := service.checkCreateParam(request)
	if err != nil {
		return err
	}
	task := &SeedSyncTaskTable{
		TaskName:     request.TaskName,
		SiteList:     strings.Join(request.SiteList, ";"),
		DownloaderId: request.DownloaderId,
		ExcludePath:  strings.Join(request.ExcludePath, ";"),
		MinSize:      request.MinSize,
		AddTag:       request.AddTag,
		Status:       request.Status,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	return service.seedSyncDAO.CreateSeedSyncTask(task)
}

// 更新辅种任务
func (service *seedSyncService) UpdateSeedSyncTask(request *UpdateSeedSyncTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	//参数校验
	err := service.checkUpdateParam(request)
	if err != nil {
		return err
	}
	task := &SeedSyncTaskTable{
		Id:           request.Id,
		TaskName:     request.TaskName,
		SiteList:     strings.Join(request.SiteList, ";"),
		DownloaderId: request.DownloaderId,
		ExcludePath:  strings.Join(request.ExcludePath, ";"),
		MinSize:      request.MinSize,
		AddTag:       request.AddTag,
		Status:       request.Status,
		UpdateTime:   time.Now(),
	}
	return service.seedSyncDAO.UpdateSeedSyncTask(task)
}

// 辅种
func (service *seedSyncService) SeedSync(taskName string) error {
	task := service.seedSyncDAO.GetSeedSyncTaskByTaskName(taskName)
	if task == nil {
		return errors.New("辅种任务" + taskName + "不存在")
	}
	if task.Status != common.SEED_SYNC_TASK_STATUS_USED {
		return errors.New("辅种任务" + taskName + "状态未启用")
	}
	return service.doSeedSync(task)
}

func (service *seedSyncService) doSeedSync(task *SeedSyncTaskTable) error {
	//辅种流程：
	//1. 根据辅种的下载器，去查询下载器下所有的种子
	//2. 根据辅种配置，过滤部分不辅种的种子
	//3. 根据辅种配置，拿到要辅种的种子和站点，向服务端发请求进行辅种，得到可以辅种的种子
	//4. 判断可以辅种的种子不存在，去相应站点下载种子
	//5. 下载种子后，调用下载器的下载接口进行辅种
	downloaderClient, err := downloader.DownloaderService.GetDownloaderById(task.DownloaderId)
	if err != nil {
		return fmt.Errorf("辅种失败：获取下载器失败， 错误: %v", err)
	}
	seeds, err := downloaderClient.GetSeedsHash()
	if err != nil {
		return fmt.Errorf("辅种失败：从下载器获取种子失败， 错误: %v", err)
	}
	//转map，key为hash
	seedMap := make(map[string]downloader.SeedHash)
	for _, seed := range seeds {
		seedMap[seed.InfoHash] = seed
	}

	//分批
	batchSeeds := make([][]downloader.SeedHash, 0)
	for i := 0; i < len(seeds); i += SEED_SYNC_BATCH_SIZE {
		end := i + SEED_SYNC_BATCH_SIZE
		if end > len(seeds) {
			end = len(seeds)
		}
		batchSeeds = append(batchSeeds, seeds[i:end])
	}
	//分批请求和辅种
	for _, batch := range batchSeeds {
		request := service.getSeedSyncRequest(batch, task)
		//无可辅种的种子，跳过
		if request == nil {
			continue
		}
		response, err := seedSyncServer.SeedSyncServerClient.SyncSeed(request)
		if err != nil {
			return fmt.Errorf("辅种失败：向seedSyncServer请求辅种种子失败， 错误: %v", err)
		}
		err = service.handleSeedSyncResponse(response, seedMap, downloaderClient)
		if err != nil {
			return fmt.Errorf("辅种失败，错误: %v", err)
		}
	}
	return nil
}

func (service *seedSyncService) handleSeedSyncResponse(response map[string][]seedSyncServer.SeedSyncTorrentInfoResponse, seedMap map[string]downloader.SeedHash, downloaderClient downloader.Downloader) error {
	//处理返回结果，对于seedMap中不存在的种子，进行下载
	for srcHash, seedForSyncList := range response {
		for _, seedForSync := range seedForSyncList {
			if _, ok := seedMap[seedForSync.InfoHash]; !ok {
				err := service.downloadAndSyncSeed(seedMap[srcHash], seedForSync, downloaderClient)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				//将种子添加到seedMap
				seedMap[seedForSync.InfoHash] = downloader.SeedHash{
					InfoHash:    seedForSync.InfoHash,
					Size:        0,
					Tags:        []string{},
					DownloadDir: "",
				}
			}
		}
	}
	return nil
}

func (service *seedSyncService) downloadAndSyncSeed(srcSeed downloader.SeedHash, seedForSync seedSyncServer.SeedSyncTorrentInfoResponse, downloaderClient downloader.Downloader) error {
	//获取站点客户端
	siteClient := site.SiteService.GetSiteClient(seedForSync.SiteName)
	if siteClient == nil {
		return fmt.Errorf("辅种失败：站点客户端不存在，站点名称: %s", seedForSync.SiteName)
	}
	//下载种子
	bytes, err := siteClient.DownloadTorrent(seedForSync.TorrentId)
	if err != nil {
		return fmt.Errorf("辅种失败：下载种子失败，站点名称: %s,种子id: %d, 错误: %v", seedForSync.SiteName, seedForSync.TorrentId, err)
	}
	//辅种
	//创建request
	//todo: tag
	request := &downloader.AddTorrentRequest{
		DownloadDir: srcSeed.DownloadDir,
		TorrentFile: bytes,
		Paused:      true,
	}
	err = downloaderClient.AddTorrent(request)
	if err != nil {
		return fmt.Errorf("辅种失败：下载器添加种子失败， 错误: %v", err)
	}
	return nil
}

func (service *seedSyncService) checkCreateParam(request *CreateSeedSyncTaskRequest) error {
	return service.checkParam(&UpdateSeedSyncTaskRequest{
		TaskName:     request.TaskName,
		SiteList:     request.SiteList,
		DownloaderId: request.DownloaderId,
		ExcludePath:  request.ExcludePath,
		MinSize:      request.MinSize,
		AddTag:       request.AddTag,
		Status:       request.Status,
	}, true)
}
func (service *seedSyncService) checkUpdateParam(request *UpdateSeedSyncTaskRequest) error {
	return service.checkParam(request, false)
}

// 参数校验
func (service *seedSyncService) checkParam(request *UpdateSeedSyncTaskRequest, create bool) error {
	//参数校验
	//0. 更新流程任务得存在
	if !create {
		task := service.seedSyncDAO.GetSeedSyncTask(request.Id)
		if task == nil {
			return errors.New("任务不存在")
		}
	}

	//1.任务名称不能为空
	if request.TaskName == "" {
		return errors.New("任务名称不能为空")
	}
	//2. 任务名不能重复
	task := service.seedSyncDAO.GetSeedSyncTaskByTaskName(request.TaskName)
	if create && task != nil {
	}
	if create && task != nil {
		return errors.New("任务名" + request.TaskName + "已存在")
	} else if !create && task != nil && task.Id != request.Id {
		return errors.New("任务名" + request.TaskName + "已存在")
	}
	//3. 站点名合法
	if len(request.SiteList) == 0 {
		return errors.New("站点名不能为空")
	}
	siteList, err := site.SiteService.GetSiteList()
	if err != nil {
		return err
	}
	siteMap := make(map[string]bool)
	for _, site := range siteList {
		siteMap[site.SiteName] = true
	}
	for _, site := range request.SiteList {
		if !siteMap[site] {
			return errors.New("站点名不合法")
		}
	}
	//4. 下载器id合法
	downloader, err := downloader.DownloaderService.GetDownloaderById(request.DownloaderId)
	if err != nil {
		return err
	}
	if downloader == nil {
		return errors.New("下载器不存在")
	}
	//5. status合法
	if request.Status != common.SEED_SYNC_TASK_STATUS_USED && request.Status != common.SEED_SYNC_TASK_STATUS_STOP {
		return errors.New("任务状态不合法")
	}
	return nil
}
func (service *seedSyncService) getSeedSyncRequest(seeds []downloader.SeedHash, task *SeedSyncTaskTable) *seedSyncServer.SeedSyncRequest {
	//向服务端请求可辅种的种子
	infoHashList := make([]string, 0)
	for _, seed := range seeds {
		//过滤掉不辅种的种子
		if seed.Size < task.MinSize {
			continue
		}
		if strings.Contains(seed.DownloadDir, task.ExcludePath) {
			continue
		}
		if common.HasSameElement(strings.Split(task.ExcludeTag, ";"), seed.Tags) {
			continue
		}
		infoHashList = append(infoHashList, seed.InfoHash)
	}
	if len(infoHashList) == 0 {
		return nil
	}
	return &seedSyncServer.SeedSyncRequest{
		InfoHash: infoHashList,
		Sites:    strings.Split(task.SiteList, ";"),
	}
}
