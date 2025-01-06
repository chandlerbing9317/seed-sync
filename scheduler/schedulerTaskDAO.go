package scheduler

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

type SchedulerTaskDAO struct {
	db *gorm.DB
}

var schedulerTaskDAO *SchedulerTaskDAO

func InitSchedulerTaskDAO() {
	schedulerTaskDAO = &SchedulerTaskDAO{
		db: db.DB,
	}
}

type SchedulerTaskTable struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	//关联的外键任务id，这个id不唯一，id+executeContent才唯一，因为定时任务可能来自不同的表
	TaskID            int64     `gorm:"column:task_id"`
	TaskName          string    `gorm:"column:task_name"`
	ExecuteContent    string    `gorm:"column:execute_content"`
	Cron              string    `gorm:"column:cron"`
	ExecuteStatus     string    `gorm:"column:execute_status"`
	LastExecuteTime   time.Time `gorm:"column:last_execute_time"`
	NextExecuteTime   time.Time `gorm:"column:next_execute_time"`
	LastExecuteResult string    `gorm:"column:last_execute_result"`
	Active            bool      `gorm:"column:active"`
	CreateUser        string    `gorm:"column:create_user"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (*SchedulerTaskTable) TableName() string {
	return "seed_sync_schedule_task"
}

func (dao *SchedulerTaskDAO) GetSchedulerTaskByTaskIDAndExecuteContent(taskID int64, executeContent string) *SchedulerTaskTable {
	var task SchedulerTaskTable
	err := dao.db.Where("task_id = ? AND execute_content = ?", taskID, executeContent).First(&task).Error
	//处理一下找不到的返回nil
	if err != nil {
		return nil
	}
	return &task
}

func (dao *SchedulerTaskDAO) CreateSchedulerTask(task *SchedulerTaskTable) error {
	return dao.db.Create(task).Error
}

func (dao *SchedulerTaskDAO) CreateSchedulerTaskWithTx(tx *gorm.DB, task *SchedulerTaskTable) error {
	return tx.Create(task).Error
}

func (dao *SchedulerTaskDAO) UpdateSchedulerTask(task *SchedulerTaskTable) error {
	return dao.db.Save(task).Error
}

func (dao *SchedulerTaskDAO) UpdateSchedulerTaskWithTx(tx *gorm.DB, task *SchedulerTaskTable) error {
	return tx.Save(task).Error
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
