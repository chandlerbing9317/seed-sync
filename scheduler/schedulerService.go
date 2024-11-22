package scheduler

import (
	"fmt"
	"seed-sync/common"
	"seed-sync/cookieCloud"
	"seed-sync/log"
	"seed-sync/seedSyncServer"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	SchedulerTaskStatusNotExecuted = "not_executed"
	SchedulerTaskStatusExecuting   = "executing"
)

type schedulerService struct {
	schedulerTaskDAO *SchedulerTaskDAO
	lock             sync.Mutex
}

var SchedulerService = &schedulerService{
	schedulerTaskDAO: schedulerTaskDAO,
	lock:             sync.Mutex{},
}

// 维护一个executeContent与实际执行函数的map
var executeFuncMap map[string]*SchedulerTask

// 目前在这里统一管控，而非采用注册式，避免杂乱不好找
func initExecuteFuncMap() {
	executeFuncMap = map[string]*SchedulerTask{
		common.CHECK_USER_EXECUTE_CONTENT: {
			ExecuteContent: common.CHECK_USER_EXECUTE_CONTENT,
			ExecuteFunc:    seedSyncServer.SeedSyncServerService.CheckUserForSchedulerTask,
		},
		common.GET_SITE_EXECUTE_CONTENT: {
			ExecuteContent: common.GET_SITE_EXECUTE_CONTENT,
			ExecuteFunc:    seedSyncServer.SeedSyncServerService.GetSiteForSchedulerTask,
		},
		common.SYNC_COOKIE_CLOUD_EXECUTE_CONTENT: {
			ExecuteContent: common.SYNC_COOKIE_CLOUD_EXECUTE_CONTENT,
			ExecuteFunc:    cookieCloud.CookieCloudService.SyncCookieForSchedulerTask,
		},
	}
}

func init() {
	initExecuteFuncMap()
}

func (service *schedulerService) CreateOrUpdateSchedulerTask(task *CreateSchedulerTaskRequest) error {
	//根据cron表达式计算下一次执行时间:
	nextExecuteTime, err := common.GetNextExecuteTime(task.Cron)
	if err != nil {
		return err
	}
	service.lock.Lock()
	defer service.lock.Unlock()

	if schedulerTask := schedulerTaskDAO.GetSchedulerTaskByName(task.TaskName); schedulerTask != nil {
		//更新
		schedulerTask.Cron = task.Cron
		schedulerTask.ExecuteContent = task.ExecuteContent
		schedulerTask.Active = task.Active
		schedulerTask.NextExecuteTime = nextExecuteTime
		schedulerTask.UpdateTime = time.Now()
		return service.schedulerTaskDAO.UpdateSchedulerTask(schedulerTask)
	} else {
		//创建
		schedulerTask := &SchedulerTaskTable{
			TaskName:        task.TaskName,
			Cron:            task.Cron,
			ExecuteContent:  task.ExecuteContent,
			Active:          task.Active,
			NextExecuteTime: nextExecuteTime,
			CreateUser:      task.CreateUser,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		return service.schedulerTaskDAO.CreateSchedulerTask(schedulerTask)
	}
}

func (service *schedulerService) GetAllSchedulerTask() ([]*SchedulerTaskTable, error) {
	return service.schedulerTaskDAO.GetAllSchedulerTask()
}

// 执行定时任务
func (service *schedulerService) ScheduleTask() error {
	service.lock.Lock()
	defer service.lock.Unlock()
	tasks, err := service.schedulerTaskDAO.GetActiveSchedulerTask()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		//判断task的时间到了执行时间，且状态是未执行
		if task.NextExecuteTime.Before(time.Now()) && task.ExecuteStatus == SchedulerTaskStatusNotExecuted {
			if schedulerTask, ok := executeFuncMap[task.ExecuteContent]; ok {
				//注：这里更新状态不能放在go routine中，
				//因为先判断未执行再更新为正执行本质是竞态条件，要在同一锁中
				task.ExecuteStatus = SchedulerTaskStatusExecuting
				task.LastExecuteTime = time.Now()
				service.schedulerTaskDAO.UpdateSchedulerTask(task)
				service.executeTask(task, schedulerTask)
			} else {
				log.Error("未找到执行函数", zap.String("executeContent", task.ExecuteContent))
			}
		}
	}
	return nil
}

// 执行任务，放在go routine中执行
func (service *schedulerService) executeTask(task *SchedulerTaskTable, schedulerTask *SchedulerTask) {
	go func(task *SchedulerTaskTable, schedulerTask *SchedulerTask) {
		defer func() {
			if r := recover(); r != nil {
				service.updateSchedulerTaskResult(task, fmt.Errorf("%v", r))
			}
		}()
		err := schedulerTask.ExecuteFunc()
		service.updateSchedulerTaskResult(task, err)
	}(task, schedulerTask)
}

func (service *schedulerService) updateSchedulerTaskResult(task *SchedulerTaskTable, err error) {
	result := ""
	if err != nil {
		result = "执行失败:" + err.Error()
		log.Error("执行定时任务失败", zap.String("taskName", task.TaskName), zap.Error(err))
	} else {
		result = "执行成功"
		log.Info("执行定时任务成功", zap.String("taskName", task.TaskName))
	}
	task.LastExecuteResult = result
	nextExecuteTime, _ := common.GetNextExecuteTime(task.Cron)
	task.NextExecuteTime = nextExecuteTime
	task.ExecuteStatus = SchedulerTaskStatusNotExecuted
	service.schedulerTaskDAO.UpdateSchedulerTask(task)
}
