package scheduler

type CreateSchedulerTaskRequest struct {
	TaskName       string `json:"taskName"`
	TaskID         int64  `json:"taskID"`
	Cron           string `json:"cron"`
	ExecuteContent string `json:"executeContent"`
	Active         bool   `json:"active"`
	CreateUser     string `json:"createUser"`
}

type SchedulerTaskExecuteFunc func(schedulerTask *SchedulerTaskTable) error

type SchedulerTask struct {
	ExecuteContent string
	ExecuteFunc    SchedulerTaskExecuteFunc
}
