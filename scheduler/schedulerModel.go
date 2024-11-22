package scheduler

type CreateSchedulerTaskRequest struct {
	TaskName       string `json:"taskName"`
	Cron           string `json:"cron"`
	ExecuteContent string `json:"executeContent"`
	Active         bool   `json:"active"`
	CreateUser     string `json:"createUser"`
}

type SchedulerTaskExecuteFunc func() error

type SchedulerTask struct {
	ExecuteContent string
	ExecuteFunc    SchedulerTaskExecuteFunc
}
