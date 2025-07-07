package task

import (
	"time"

	"github.com/bestruirui/bestsub/internal/models/common"
)

// 任务类型枚举
const (
	// 系统任务
	TypeNodePoolCleanup = "nodepool_cleanup" // 节点池清理任务
	TypeGC              = "gc"               // 垃圾回收任务

	// 用户任务
	TypeLinkFetch  = "link_fetch"  // 链接内容获取任务
	TypeNodeDetect = "node_detect" // 节点检测任务
	TypeNodeSave   = "node_save"   // 节点保存任务
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
	Cron            string     `json:"cron" example:"0 */6 * * *"`               // Cron表达式
	Timeout         int        `db:"timeout" json:"timeout"`                     // 任务超时时间（秒）
	Type            string     `json:"type" example:"link_fetch"`                // 任务类型
	Config          string     `json:"config" example:"{\"sub_link_id\": 1}"`    // 任务配置（JSON格式）
	Status          string     `db:"status" json:"status" example:"pending"`     // 任务状态
	SuccessCount    int        `db:"success_count" json:"success_count"`         // 成功次数
	FailedCount     int        `db:"failed_count" json:"failed_count"`           // 失败次数
	LastRunResult   string     `db:"last_run_result" json:"last_run_result"`     // 上次执行结果
	LastRunTime     *time.Time `db:"last_run_time" json:"last_run_time"`         // 上次执行时间
	LastRunDuration *int       `db:"last_run_duration" json:"last_run_duration"` // 上次执行耗时（毫秒）
}
type CreateRequest struct {
	common.BaseRequestModel
	Cron   string `json:"cron" example:"0 */6 * * *"`            // Cron表达式
	Type   string `json:"type" example:"link_fetch"`             // 任务类型
	Config string `json:"config" example:"{\"sub_link_id\": 1}"` // 任务配置（JSON格式）
}
type UpdateRequest struct {
	common.BaseUpdateRequestModel
	Cron   string `json:"cron" example:"0 */6 * * *"`            // Cron表达式
	Config string `json:"config" example:"{\"sub_link_id\": 1}"` // 任务配置（JSON格式）
}
