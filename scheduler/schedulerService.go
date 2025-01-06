package scheduler

import (
	"fmt"
	"seed-sync/common"
	"seed-sync/log"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	SchedulerTaskStatusNotExecuted = "not_executed"
	SchedulerTaskStatusExecuting   = "executing"
)

type schedulerService struct {
	schedulerTaskDAO *SchedulerTaskDAO
	executeFuncMap   map[string]*SchedulerTask
	lock             sync.Mutex
}

var SchedulerService *schedulerService

// 初始化的时候，从库里查出所有的定时任务并初始化保存
var once sync.Once

func InitSchedulerService() {
	once.Do(func() {
		SchedulerService = &schedulerService{
			schedulerTaskDAO: schedulerTaskDAO,
			executeFuncMap:   make(map[string]*SchedulerTask),
			lock:             sync.Mutex{},
		}
		SchedulerService.resetSchedulerService()
	})
}

// 服务初始化的时候，重置所有定时任务的下一次执行时间和状态
func (service *schedulerService) resetSchedulerService() {
	tasks, err := service.schedulerTaskDAO.GetActiveSchedulerTask()
	if err != nil {
		log.Fatal("初始化获取定时任务失败", zap.Error(err))
	}
	for _, task := range tasks {
		nextExecuteTime, err := common.GetNextExecuteTime(task.Cron)
		if err != nil {
			log.Fatal("初始化cron计算下一次执行时间失败", zap.String("cron", task.Cron), zap.Error(err))
		}
		task.NextExecuteTime = nextExecuteTime
		task.ExecuteStatus = SchedulerTaskStatusNotExecuted
		service.schedulerTaskDAO.UpdateSchedulerTask(task)
	}
}

func (service *schedulerService) RegisterExecuteFunc(executeContent string, executeFunc SchedulerTaskExecuteFunc) {
	service.lock.Lock()
	defer service.lock.Unlock()
	service.executeFuncMap[executeContent] = &SchedulerTask{
		ExecuteContent: executeContent,
		ExecuteFunc:    executeFunc,
	}
}

func (service *schedulerService) CreateOrUpdateSchedulerTask(task *CreateSchedulerTaskRequest) error {
	//根据cron表达式计算下一次执行时间:
	nextExecuteTime, err := common.GetNextExecuteTime(task.Cron)
	if err != nil {
		return fmt.Errorf("cron计算下一次执行时间失败，cron: %s, 错误: %v", task.Cron, err)
	}
	service.lock.Lock()
	defer service.lock.Unlock()

	if schedulerTask := schedulerTaskDAO.GetSchedulerTaskByTaskIDAndExecuteContent(task.TaskID, task.ExecuteContent); schedulerTask != nil {
		//更新
		schedulerTask.TaskName = task.TaskName
		schedulerTask.Cron = task.Cron
		schedulerTask.Active = task.Active
		schedulerTask.NextExecuteTime = nextExecuteTime
		schedulerTask.UpdateTime = time.Now()
		return service.schedulerTaskDAO.UpdateSchedulerTask(schedulerTask)
	} else {
		//创建
		schedulerTask := &SchedulerTaskTable{
			TaskID:          task.TaskID,
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

func (service *schedulerService) CreateSchedulerTaskWithTx(tx *gorm.DB, task *CreateSchedulerTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	//校验taskID+executeContent 唯一
	schedulerTask := service.GetSchedulerTaskByTaskIDAndExecuteContent(task.TaskID, task.ExecuteContent)
	if schedulerTask != nil {
		return fmt.Errorf("定时任务taskID:%d,executeContent:%s已存在", task.TaskID, task.ExecuteContent)
	}
	return service.createOrUpdateSchedulerTaskWithTx(tx, task)
}

func (service *schedulerService) UpdateSchedulerTaskWithTx(tx *gorm.DB, task *CreateSchedulerTaskRequest) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	schedulerTask := service.GetSchedulerTaskByTaskIDAndExecuteContent(task.TaskID, task.ExecuteContent)
	if schedulerTask == nil {
		return fmt.Errorf("定时任务taskID:%d,executeContent:%s不存在", task.TaskID, task.ExecuteContent)
	}
	return service.createOrUpdateSchedulerTaskWithTx(tx, task)
}

func (service *schedulerService) createOrUpdateSchedulerTaskWithTx(tx *gorm.DB, task *CreateSchedulerTaskRequest) error {
	//根据cron表达式计算下一次执行时间:
	nextExecuteTime, err := common.GetNextExecuteTime(task.Cron)
	if err != nil {
		return fmt.Errorf("cron计算下一次执行时间失败，cron: %s, 错误: %v", task.Cron, err)
	}
	if schedulerTask := schedulerTaskDAO.GetSchedulerTaskByTaskIDAndExecuteContent(task.TaskID, task.ExecuteContent); schedulerTask != nil {
		//更新
		schedulerTask.TaskName = task.TaskName
		schedulerTask.Cron = task.Cron
		schedulerTask.Active = task.Active
		schedulerTask.NextExecuteTime = nextExecuteTime
		schedulerTask.UpdateTime = time.Now()
		return service.schedulerTaskDAO.UpdateSchedulerTaskWithTx(tx, schedulerTask)
	} else {
		//创建
		schedulerTask := &SchedulerTaskTable{
			TaskID:          task.TaskID,
			TaskName:        task.TaskName,
			ExecuteContent:  task.ExecuteContent,
			ExecuteStatus:   SchedulerTaskStatusNotExecuted,
			Cron:            task.Cron,
			Active:          task.Active,
			NextExecuteTime: nextExecuteTime,
			CreateUser:      task.CreateUser,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		return service.schedulerTaskDAO.CreateSchedulerTaskWithTx(tx, schedulerTask)
	}
}

func (service *schedulerService) GetSchedulerTaskByTaskIDAndExecuteContent(taskID int64, executeContent string) *SchedulerTaskTable {
	return service.schedulerTaskDAO.GetSchedulerTaskByTaskIDAndExecuteContent(taskID, executeContent)
}

// 执行定时任务
func (service *schedulerService) ExecuteSchedulerTask() error {
	service.lock.Lock()
	defer service.lock.Unlock()
	tasks, err := service.schedulerTaskDAO.GetActiveSchedulerTask()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		//判断task的时间到了执行时间，且状态是未执行
		if task.NextExecuteTime.Before(time.Now()) && task.ExecuteStatus == SchedulerTaskStatusNotExecuted {
			if schedulerTask, ok := service.executeFuncMap[task.ExecuteContent]; ok {
				//注：这里更新状态不能放在go routine中，
				//因为先判断未执行再更新为正执行属于竞态条件，要在同一锁中
				task.ExecuteStatus = SchedulerTaskStatusExecuting
				task.LastExecuteTime = time.Now()
				service.schedulerTaskDAO.UpdateSchedulerTask(task)
				service.doExecute(task, schedulerTask)
			} else {
				log.Error("未找到执行函数", zap.String("executeContent", task.ExecuteContent))
			}
		}
	}
	return nil
}

// 执行任务，放在go routine中执行
func (service *schedulerService) doExecute(task *SchedulerTaskTable, schedulerTask *SchedulerTask) {
	go func(task *SchedulerTaskTable, schedulerTask *SchedulerTask) {
		defer func() {
			if r := recover(); r != nil {
				service.updateSchedulerTaskResult(task, fmt.Errorf("%v", r))
			}
		}()
		err := schedulerTask.ExecuteFunc(task)
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
