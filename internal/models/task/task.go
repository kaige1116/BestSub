package task

import (
	"time"
)

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name" description:"任务名称"`
	Enable bool   `db:"enable" json:"enable" description:"是否启用"`
	Config string `db:"config" json:"config" description:"任务的配置"`
	Extra  string `db:"extra" json:"extra" description:"任务的额外参数"`
	Result string `db:"result" json:"result" description:"结果"`
}

type Config struct {
	ID            uint16 `json:"id" description:"任务ID"`
	Name          string `json:"name" description:"任务名称"`
	CronExpr      string `json:"cron_expr" description:"cron表达式"`
	Notify        bool   `json:"notify" description:"是否通知"`
	NotifyChannel int    `json:"notify_channel" description:"通知渠道"`
	LogWriteFile  bool   `json:"log_write_file" description:"是否写入日志文件"`
	LogLevel      string `json:"log_level" description:"日志级别"`
	Timeout       int    `json:"timeout" description:"超时时间"`
	Type          string `json:"type" description:"任务类型"`
	Extra         string `json:"extra" description:"任务的额外参数"`
}

type Result struct {
	Success         int       `json:"success"`
	Failed          int       `json:"failed"`
	LastRunResult   string    `json:"last_run_result"`
	LastRunTime     time.Time `json:"last_run_time"`
	LastRunDuration int       `json:"last_run_duration"`
}

type Response struct {
	Data
	Status string `json:"status"`
}

type CreateRequest struct {
	Name   string `json:"name" example:"test"`
	Enable bool   `json:"enable" example:"true"`
	Config string `json:"config" example:"{\"cron_expr\":\"0 0 * * *\",\"notify\":true,\"notify_channel\":1,\"log_write_file\":true,\"log_level\":\"info\",\"timeout\":60,\"type\":\"test\"}"`
	Extra  string `json:"extra" example:"{\"id\":1}"`
}
type UpdateRequest struct {
	ID     uint16 `json:"id"`
	Name   string `json:"name" example:"test"`
	Enable bool   `json:"enable" example:"true"`
	Config string `json:"config" example:"{\"cron_expr\":\"0 0 * * *\",\"notify\":true,\"notify_channel\":1,\"log_write_file\":true,\"log_level\":\"info\",\"timeout\":60,\"type\":\"test\"}"`
	Extra  string `json:"extra" example:"{\"id\":1}"`
}
