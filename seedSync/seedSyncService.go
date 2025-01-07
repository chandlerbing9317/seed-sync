package seedSync

import (
	"errors"
	"fmt"
	"seed-sync/common"
	"seed-sync/downloader"
	"seed-sync/log"
	"seed-sync/scheduler"
	"seed-sync/seedSyncServer"
	"seed-sync/site"
	"seed-sync/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const SEED_SYNC_BATCH_SIZE = 200

type seedSyncService struct {
	seedSyncDAO *SeedSyncDAO
	lock        sync.Mutex
}

var SeedSyncService *seedSyncService

// 初始化的时候，从库里查出所有的辅种任务并初始化保存
var once sync.Once

func InitSeedSyncService() {
	once.Do(func() {
		SeedSyncService = &seedSyncService{
			seedSyncDAO: seedSyncDAO,
			lock:        sync.Mutex{},
		}
	})
}

// 创建辅种任务
func (service *seedSyncService) CreateSeedSyncTask(request *CreateSeedSyncTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	//开启事务
	tx := service.seedSyncDAO.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("创建辅种任务失败，错误: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	task := &SeedSyncTaskTable{
		TaskName:     request.TaskName,
		SiteList:     strings.Join(request.SiteList, ";"),
		DownloaderId: request.DownloaderId,
		ExcludePath:  strings.Join(request.ExcludePath, ";"),
		MinSize:      request.MinSize,
		AddTag:       strings.Join(request.AddTag, ";"),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	// 创建任务
	err := service.seedSyncDAO.CreateSeedSyncTaskWithTx(tx, task)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建辅种任务失败，错误: %v", err)
	}

	// 获取刚创建的任务ID
	newTask := service.seedSyncDAO.GetSeedSyncTaskByTaskNameWithTx(tx, task.TaskName)
	if newTask == nil {
		tx.Rollback()
		return fmt.Errorf("创建辅种任务失败，错误: %v", err)
	}

	//为辅种配置定时任务
	err = scheduler.SchedulerService.CreateSchedulerTaskWithTx(tx, &scheduler.CreateSchedulerTaskRequest{
		TaskID:         newTask.ID, // 使用获取到的ID
		TaskName:       task.TaskName,
		Cron:           request.Cron,
		ExecuteContent: common.SYNC_SEED_EXECUTE_CONTENT,
		Active:         request.Status == common.SEED_SYNC_TASK_STATUS_USED,
		CreateUser:     "user",
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 更新辅种任务
func (service *seedSyncService) UpdateSeedSyncTask(request *UpdateSeedSyncTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	//开启事务
	tx := service.seedSyncDAO.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	task := &SeedSyncTaskTable{
		ID:           request.ID,
		TaskName:     request.TaskName,
		SiteList:     strings.Join(request.SiteList, ";"),
		DownloaderId: request.DownloaderId,
		ExcludePath:  strings.Join(request.ExcludePath, ";"),
		MinSize:      request.MinSize,
		AddTag:       strings.Join(request.AddTag, ";"),
		UpdateTime:   time.Now(),
	}
	err := service.seedSyncDAO.UpdateSeedSyncTaskWithTx(tx, task)
	if err != nil {
		tx.Rollback()
		return err
	}
	//更新定时任务
	err = scheduler.SchedulerService.UpdateSchedulerTaskWithTx(tx, &scheduler.CreateSchedulerTaskRequest{
		TaskID:         task.ID,
		TaskName:       task.TaskName,
		Cron:           request.Cron,
		ExecuteContent: common.SYNC_SEED_EXECUTE_CONTENT,
		Active:         request.Status == common.SEED_SYNC_TASK_STATUS_USED,
		CreateUser:     "user",
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 查询所有辅种任务
func (service *seedSyncService) GetSeedSyncTaskList() ([]*SeedSyncTaskInfo, error) {
	taskList := service.seedSyncDAO.GetAllSeedSyncTaskList()
	taskInfoList := make([]*SeedSyncTaskInfo, 0)
	for _, task := range taskList {
		schedulerTask := scheduler.SchedulerService.GetSchedulerTaskByTaskIDAndExecuteContent(task.ID, common.SYNC_SEED_EXECUTE_CONTENT)
		taskInfo := service.getTaskInfo(task, schedulerTask)
		taskInfoList = append(taskInfoList, taskInfo)
	}
	return taskInfoList, nil
}

// 查询辅种任务
func (service *seedSyncService) GetSeedSyncTaskByName(taskName string) *SeedSyncTaskInfo {
	task := service.seedSyncDAO.GetSeedSyncTaskByTaskName(taskName)
	if task == nil {
		return nil
	}
	schedulerTask := scheduler.SchedulerService.GetSchedulerTaskByTaskIDAndExecuteContent(task.ID, common.SYNC_SEED_EXECUTE_CONTENT)
	return service.getTaskInfo(task, schedulerTask)
}

// 查询辅种任务
func (service *seedSyncService) GetSeedSyncTaskById(id int64) *SeedSyncTaskInfo {
	task := service.seedSyncDAO.GetSeedSyncTask(id)
	if task == nil {
		return nil
	}
	schedulerTask := scheduler.SchedulerService.GetSchedulerTaskByTaskIDAndExecuteContent(task.ID, common.SYNC_SEED_EXECUTE_CONTENT)
	return service.getTaskInfo(task, schedulerTask)
}

func (service *seedSyncService) getTaskInfo(task *SeedSyncTaskTable, schedulerTask *scheduler.SchedulerTaskTable) *SeedSyncTaskInfo {
	status := common.SEED_SYNC_TASK_STATUS_STOP
	if schedulerTask.Active {
		status = common.SEED_SYNC_TASK_STATUS_USED
	}
	return &SeedSyncTaskInfo{
		ID: task.ID,
		SeedSyncTaskBaseInfo: SeedSyncTaskBaseInfo{
			TaskName:     task.TaskName,
			SiteList:     strings.Split(task.SiteList, ";"),
			DownloaderId: task.DownloaderId,
			ExcludePath:  strings.Split(task.ExcludePath, ";"),
			MinSize:      task.MinSize,
			AddTag:       strings.Split(task.AddTag, ";"),
			Status:       status,
			Cron:         schedulerTask.Cron,
		},
		CreateTime: task.CreateTime,
		UpdateTime: task.UpdateTime,
	}
}

func (service *seedSyncService) SeedSyncScheduler(schedulerTask *scheduler.SchedulerTaskTable) error {
	return service.SeedSync(schedulerTask.TaskID)
}

// 辅种
func (service *seedSyncService) SeedSync(taskId int64) error {
	task := service.GetSeedSyncTaskById(taskId)
	if task == nil {
		return errors.New("辅种任务" + strconv.FormatInt(taskId, 10) + "不存在")
	}
	if task.Status != common.SEED_SYNC_TASK_STATUS_USED {
		return errors.New("辅种任务" + strconv.FormatInt(taskId, 10) + "状态未启用")
	}
	log.Info("开始辅种任务", zap.Int64("taskId", taskId), zap.String("taskName", task.TaskName))
	err := service.doSeedSync(task)
	if err != nil {
		log.Error("辅种任务"+strconv.FormatInt(taskId, 10)+"辅种失败", zap.String("taskName", task.TaskName), zap.Error(err))
		return err
	}
	log.Info("辅种任务"+strconv.FormatInt(taskId, 10)+"辅种成功", zap.String("taskName", task.TaskName))
	return nil
}

func (service *seedSyncService) doSeedSync(task *SeedSyncTaskInfo) error {
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

	//过滤
	filteredSeeds := make([]downloader.SeedHash, 0)
	for _, seed := range seeds {
		if seed.Status != downloader.SEED_STATUS_SEEDING {
			continue
		}
		if task.MinSize > 0 && seed.Size < task.MinSize {
			continue
		}
		if common.HasSameElement(task.ExcludePath, []string{seed.DownloadDir}) {
			continue
		}
		if common.HasSameElement(task.ExcludeTag, seed.Tags) {
			continue
		}
		filteredSeeds = append(filteredSeeds, seed)
	}
	//分批
	batchSeeds := make([][]downloader.SeedHash, 0)
	for i := 0; i < len(filteredSeeds); i += SEED_SYNC_BATCH_SIZE {
		end := i + SEED_SYNC_BATCH_SIZE
		if end > len(filteredSeeds) {
			end = len(filteredSeeds)
		}
		batchSeeds = append(batchSeeds, filteredSeeds[i:end])
	}
	//分批请求和辅种
	for _, batch := range batchSeeds {
		request := getSeedSyncRequest(batch, task)
		//无可辅种的种子，跳过
		if request == nil {
			continue
		}
		response, err := seedSyncServer.SeedSyncServerService.SeedSync(request, user.UserService.GetUser().Token)
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

// 下载种子并调用下载器客户端添加种子然后辅种
func (service *seedSyncService) downloadAndSyncSeed(srcSeed downloader.SeedHash, seedForSync seedSyncServer.SeedSyncTorrentInfoResponse, downloaderClient downloader.Downloader) error {
	//获取站点客户端
	siteClient := site.SiteService.GetSiteClient(seedForSync.SiteName)
	if siteClient == nil {
		return fmt.Errorf("辅种失败：站点客户端不存在，站点名称: %s", seedForSync.SiteName)
	}
	log.Info("站点种子不在下载器中，准备开始下载",
		zap.String("siteName", seedForSync.SiteName),
		zap.Int("torrentId", seedForSync.TorrentId),
		zap.String("infoHash", seedForSync.InfoHash),
		zap.Int64("id in downloader", srcSeed.ID),
		zap.String("downloadDir in downloader", srcSeed.DownloadDir),
	)
	//下载种子 //todo 这里要加一定的休眠 避免太快触发限流 但休眠按站点加，每个站点并行处理
	bytes, err := siteClient.DownloadTorrent(seedForSync.TorrentId)
	time.Sleep(5 * time.Second)
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

func getSeedSyncRequest(seeds []downloader.SeedHash, task *SeedSyncTaskInfo) *seedSyncServer.SeedSyncRequest {
	//向服务端请求可辅种的种子
	if len(seeds) == 0 {
		return nil
	}
	infoHashList := make([]string, 0)
	for _, seed := range seeds {
		infoHashList = append(infoHashList, seed.InfoHash)
	}

	return &seedSyncServer.SeedSyncRequest{
		InfoHash: infoHashList,
		Sites:    task.SiteList,
	}
}
