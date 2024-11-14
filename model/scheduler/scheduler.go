package schedulerModel

type CreateSchedulerTaskRequest struct {
	TaskName       string `json:"task_name"`
	Cron           string `json:"cron"`
	ExecuteContent string `json:"execute_content"`
	Active         bool   `json:"active"`
}
