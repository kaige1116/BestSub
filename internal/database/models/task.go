package models

import (
	"time"
)

// 任务类型枚举
const (
	// 系统任务
	TaskTypeNodePoolCleanup = "nodepool_cleanup" // 节点池清理任务
	TaskTypeGC              = "gc"               // 垃圾回收任务

	// 用户任务
	TaskTypeLinkFetch    = "link_fetch"   // 链接内容获取任务
	TaskTypeNodeDetect   = "node_detect"  // 节点检测任务
	TaskTypeNotification = "notification" // 通知任务
	TaskTypeNodeSave     = "node_save"    // 节点保存任务
)

// 任务状态枚举
const (
	TaskStatusPending   = "pending"   // 等待执行
	TaskStatusRunning   = "running"   // 正在执行
	TaskStatusCompleted = "completed" // 执行完成
	TaskStatusFailed    = "failed"    // 执行失败
	TaskStatusCancelled = "cancelled" // 已取消
)

// 任务优先级枚举
const (
	TaskPriorityHigh   = 1 // 高优先级（系统任务）
	TaskPriorityNormal = 2 // 普通优先级（用户任务）
	TaskPriorityLow    = 3 // 低优先级
)

// 记录系统中所有任务的执行情况
type Task struct {
	ID          int64     `db:"id" json:"id"`
	Type        string    `db:"type" json:"type"`                 // 任务类型
	Name        string    `db:"name" json:"name"`                 // 任务名称
	Description string    `db:"description" json:"description"`   // 任务描述
	Status      string    `db:"status" json:"status"`             // 任务状态
	Priority    int       `db:"priority" json:"priority"`         // 任务优先级
	LinkID      int64     `db:"link_id" json:"link_id,omitempty"` // 关联的链接ID（用户任务）
	Config      string    `db:"config" json:"config"`             // 任务配置（JSON格式）
	Result      string    `db:"result" json:"result"`             // 任务结果（JSON格式）
	ErrorMsg    string    `db:"error_msg" json:"error_msg"`       // 错误信息
	StartTime   time.Time `db:"start_time" json:"start_time"`     // 开始时间
	EndTime     time.Time `db:"end_time" json:"end_time"`         // 结束时间
	Duration    int       `db:"duration" json:"duration"`         // 执行耗时（毫秒）
	RetryCount  int       `db:"retry_count" json:"retry_count"`   // 重试次数
	MaxRetries  int       `db:"max_retries" json:"max_retries"`   // 最大重试次数
	NextRun     time.Time `db:"next_run" json:"next_run"`         // 下次执行时间（定时任务）
	CronExpr    string    `db:"cron_expr" json:"cron_expr"`       // Cron表达式（定时任务）
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
