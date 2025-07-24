package task

import "time"

type Config struct {
	ID            uint16 `json:"-"`
	Name          string `json:"-"`
	CronExpr      string `json:"cron_expr" example:"0 0 * * *" description:"cron表达式"`
	Notify        bool   `json:"notify" example:"true" description:"是否通知"`
	NotifyChannel int    `json:"notify_channel" example:"1" description:"通知渠道"`
	LogWriteFile  bool   `json:"log_write_file" example:"true" description:"是否写入日志文件"`
	LogLevel      string `json:"log_level" example:"info" description:"日志级别"`
	Timeout       int    `json:"timeout" example:"60" description:"超时时间"`
	Type          string `json:"type" example:"test" description:"任务类型"`
}

type DBResult struct {
	Success         uint32    `json:"success"`
	Failed          uint32    `json:"failed"`
	LastRunResult   string    `json:"last_run_result"`
	LastRunTime     time.Time `json:"last_run_time"`
	LastRunDuration uint32    `json:"last_run_duration"`
}

type ReturnResult struct {
	Status          bool      `json:"status"`
	LastRunResult   string    `json:"last_run_result"`
	LastRunTime     time.Time `json:"last_run_time"`
	LastRunDuration uint32    `json:"last_run_duration"`
}
