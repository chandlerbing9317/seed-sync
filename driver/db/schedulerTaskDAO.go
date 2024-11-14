package db

import (
	"seed-sync/common"
	schedulerModel "seed-sync/model/scheduler"
	"time"

	"gorm.io/gorm"
)

type SchedulerTaskDAO struct {
	db *gorm.DB
}

var SchedulerTaskDao = &SchedulerTaskDAO{
	db: DB,
}

type SchedulerTask struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskName          string    `json:"task_name" gorm:"column:task_name"`
	Cron              string    `json:"cron" gorm:"column:cron"`
	ExecuteContent    string    `json:"execute_content" gorm:"column:execute_content"`
	ExecuteStatus     string    `json:"execute_status" gorm:"column:execute_status"`
	LastExecuteTime   time.Time `json:"last_execute_time" gorm:"column:last_execute_time"`
	NextExecuteTime   time.Time `json:"next_execute_time" gorm:"column:next_execute_time"`
	LastExecuteResult string    `json:"last_execute_result" gorm:"column:last_execute_result"`
	Active            bool      `json:"active" gorm:"column:active"`
	CreateTime        time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime        time.Time `json:"update_time" gorm:"column:update_time"`
}

func (dao *SchedulerTaskDAO) GetSchedulerTaskByName(taskName string) *SchedulerTask {
	var task SchedulerTask
	err := dao.db.Where("task_name = ?", taskName).First(&task).Error
	//处理一下找不到的返回nil
	if err != nil {
		return nil
	}
	return &task
}

func (dao *SchedulerTaskDAO) AddOrUpdateSchedulerTask(task *schedulerModel.CreateSchedulerTaskRequest) error {
	//根据cron表达式计算下一次执行时间:
	nextExecuteTime, err := common.GetNextExecuteTime(task.Cron)
	if err != nil {
		return err
	}

	if schedulerTask := dao.GetSchedulerTaskByName(task.TaskName); schedulerTask != nil {
		//更新
		schedulerTask.Cron = task.Cron
		schedulerTask.ExecuteContent = task.ExecuteContent
		schedulerTask.Active = task.Active
		schedulerTask.NextExecuteTime = nextExecuteTime
		schedulerTask.UpdateTime = time.Now()
		return dao.db.Save(schedulerTask).Error
	} else {
		//创建
		schedulerTask := &SchedulerTask{
			TaskName:        task.TaskName,
			Cron:            task.Cron,
			ExecuteContent:  task.ExecuteContent,
			Active:          task.Active,
			NextExecuteTime: nextExecuteTime,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		return dao.db.Create(schedulerTask).Error
	}
}
