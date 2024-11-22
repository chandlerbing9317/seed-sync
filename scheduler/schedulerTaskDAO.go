package scheduler

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

type SchedulerTaskDAO struct {
	db *gorm.DB
}

var schedulerTaskDAO = &SchedulerTaskDAO{
	db: db.DB,
}

type SchedulerTaskTable struct {
	ID                int64     `gorm:"primaryKey;autoIncrement"`
	TaskName          string    `gorm:"column:task_name"`
	Cron              string    `gorm:"column:cron"`
	ExecuteContent    string    `gorm:"column:execute_content"`
	ExecuteStatus     string    `gorm:"column:execute_status"`
	LastExecuteTime   time.Time `gorm:"column:last_execute_time"`
	NextExecuteTime   time.Time `gorm:"column:next_execute_time"`
	LastExecuteResult string    `gorm:"column:last_execute_result"`
	Active            bool      `gorm:"column:active"`
	CreateUser        string    `gorm:"column:create_user"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (dao *SchedulerTaskDAO) GetSchedulerTaskByName(taskName string) *SchedulerTaskTable {
	var task SchedulerTaskTable
	err := dao.db.Where("task_name = ?", taskName).First(&task).Error
	//处理一下找不到的返回nil
	if err != nil {
		return nil
	}
	return &task
}

func (dao *SchedulerTaskDAO) CreateSchedulerTask(task *SchedulerTaskTable) error {
	return dao.db.Create(task).Error
}

func (dao *SchedulerTaskDAO) UpdateSchedulerTask(task *SchedulerTaskTable) error {
	return dao.db.Save(task).Error
}

func (dao *SchedulerTaskDAO) GetActiveSchedulerTask() ([]*SchedulerTaskTable, error) {
	var tasks []*SchedulerTaskTable
	err := dao.db.Where("active = ?", true).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (dao *SchedulerTaskDAO) GetAllSchedulerTask() ([]*SchedulerTaskTable, error) {
	var tasks []*SchedulerTaskTable
	err := dao.db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
