package seedSync

import (
	"fmt"
	"net/http"
	"seed-sync/common"
	"seed-sync/downloader"
	"seed-sync/site"

	"github.com/gin-gonic/gin"
)

func CreateSeedSyncTask(c *gin.Context) {
	request := &CreateSeedSyncTaskRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	if err := checkCreateParam(request); err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	err = SeedSyncService.CreateSeedSyncTask(request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("辅种任务创建成功"))
}

// 更新辅种任务
func UpdateSeedSyncTask(c *gin.Context) {
	request := &UpdateSeedSyncTaskRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	if err := checkUpdateParam(request); err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	err = SeedSyncService.UpdateSeedSyncTask(request)
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult("辅种任务更新成功"))
}

func GetSeedSyncTaskList(c *gin.Context) {
	taskList, err := SeedSyncService.GetSeedSyncTaskList()
	if err != nil {
		c.JSON(http.StatusOK, common.FailResult(err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResult(taskList))
}

func checkCreateParam(request *CreateSeedSyncTaskRequest) error {
	return checkParam(&UpdateSeedSyncTaskRequest{
		SeedSyncTaskBaseInfo: request.SeedSyncTaskBaseInfo,
	}, true)
}

func checkUpdateParam(request *UpdateSeedSyncTaskRequest) error {
	return checkParam(request, false)
}

// 参数校验
func checkParam(request *UpdateSeedSyncTaskRequest, create bool) error {
	//参数校验
	//0. 更新流程任务得存在
	if !create {
		task := SeedSyncService.GetSeedSyncTaskById(request.ID)
		if task == nil {
			return fmt.Errorf("%s: %d", "任务不存在", request.ID)
		}
	}
	//1.任务名称不能为空
	if request.TaskName == "" {
		return fmt.Errorf("%s: %s", "任务名称不能为空", request.TaskName)
	}
	//2. 任务名不能重复
	task := SeedSyncService.GetSeedSyncTaskByName(request.TaskName)
	if create && task != nil {
		return fmt.Errorf("%s: %s", "任务名已存在", request.TaskName)
	} else if !create && task != nil && task.ID != request.ID {
		return fmt.Errorf("%s: %s", "任务名已存在", request.TaskName)
	}
	//3. 站点名合法
	if len(request.SiteList) == 0 {
		return fmt.Errorf("%s: %s", "站点名不能为空", request.SiteList)
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
			return fmt.Errorf("%s: %s", "站点未添加", site)
		}
	}
	//4. 下载器id合法
	downloader, err := downloader.DownloaderService.GetDownloaderById(request.DownloaderId)
	if err != nil {
		return err
	}
	if downloader == nil {
		return fmt.Errorf("%s: %d", "下载器不存在", request.DownloaderId)
	}
	//5. status合法
	if request.Status != common.SEED_SYNC_TASK_STATUS_USED && request.Status != common.SEED_SYNC_TASK_STATUS_STOP {
		return fmt.Errorf("%s: %s", "任务状态不合法", request.Status)
	}

	//6. cron表达式合法
	if err := common.CheckCronExpr(request.Cron); err != nil {
		return fmt.Errorf("%s: %s", "cron表达式不合法", err.Error())
	}
	return nil
}
