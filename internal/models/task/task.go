package task

import (
	"time"

	"github.com/bestruirui/bestsub/internal/models/common"
)

// 任务状态枚举
const (
	StatusPending   = "pending"   // 等待执行
	StatusRunning   = "running"   // 正在执行
	StatusCompleted = "completed" // 执行完成
	StatusFailed    = "failed"    // 执行失败
	StatusCancelled = "cancelled" // 已取消
)

// Task 任务基础模型（数据库模型）
type Data struct {
	common.BaseDbModel
	IsSysTask       bool       `db:"is_sys_task" json:"is_sys_task"`             // 是否系统任务
	Cron            string     `db:"cron" json:"cron" example:"0 */6 * * *"`     // Cron表达式
	Type            string     `db:"type" json:"type" example:"sub_fetch"`       // 任务类型
	Timeout         int        `db:"timeout" json:"timeout" example:"60"`        // 任务超时时间（秒）
	Retry           int        `db:"retry" json:"retry" example:"3"`             // 任务重试次数
	Config          string     `db:"config" json:"config"`                       // 任务配置（JSON格式）
	Status          string     `db:"-" json:"status" example:"pending"`          // 任务状态
	SuccessCount    int        `db:"success_count" json:"success_count"`         // 成功次数
	FailedCount     int        `db:"failed_count" json:"failed_count"`           // 失败次数
	LastRunResult   string     `db:"last_run_result" json:"last_run_result"`     // 上次执行结果
	LastRunTime     *time.Time `db:"last_run_time" json:"last_run_time"`         // 上次执行时间
	LastRunDuration *int       `db:"last_run_duration" json:"last_run_duration"` // 上次执行耗时（毫秒）
}

type CreateRequest struct {
	common.BaseRequestModel
	Cron    string `json:"cron" example:"0 */6 * * *"`                                                                                                                       // Cron表达式
	Type    string `json:"type" example:"sub_fetch"`                                                                                                                         // 任务类型
	Config  string `json:"config" example:"{\"proxy_enable\":false,\"retries\":3,\"sub_id\":1,\"timeout\":30,\"type\":\"auto\",\"url\":\"\",\"user_agent\":\"clash.meta\"}"` // 任务配置（JSON格式）
	Timeout int    `db:"timeout" json:"timeout" example:"60"`                                                                                                                // 任务超时时间（秒）
	Retry   int    `db:"retry" json:"retry" example:"3"`                                                                                                                     // 任务重试次数
}
type UpdateRequest struct {
	common.BaseUpdateRequestModel
	Cron    string `json:"cron" example:"0 */6 * * *"`                                                                                                                       // Cron表达式
	Config  string `json:"config" example:"{\"proxy_enable\":false,\"retries\":3,\"sub_id\":1,\"timeout\":30,\"type\":\"auto\",\"url\":\"\",\"user_agent\":\"clash.meta\"}"` // 任务配置（JSON格式）
	Timeout int    `db:"timeout" json:"timeout" example:"60"`                                                                                                                // 任务超时时间（秒）
	Retry   int    `db:"retry" json:"retry" example:"3"`                                                                                                                     // 任务重试次数
}
