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
	TaskTypeLinkFetch  = "link_fetch"  // 链接内容获取任务
	TaskTypeNodeDetect = "node_detect" // 节点检测任务
	TaskTypeNodeSave   = "node_save"   // 节点保存任务
)

// 任务状态枚举
const (
	TaskStatusPending   = "pending"   // 等待执行
	TaskStatusRunning   = "running"   // 正在执行
	TaskStatusCompleted = "completed" // 执行完成
	TaskStatusFailed    = "failed"    // 执行失败
	TaskStatusCancelled = "cancelled" // 已取消
)

// 记录系统中所有任务的执行情况
type Task struct {
	ID              int64     `db:"id" json:"id"`
	Enable          bool      `db:"enable" json:"enable"`                       // 是否启用
	Name            string    `db:"name" json:"name"`                           // 任务名称
	Cron            string    `db:"cron" json:"cron_expr"`                      // Cron表达式（定时任务）
	Type            string    `db:"type" json:"type"`                           // 任务类型
	Status          string    `db:"status" json:"status"`                       // 任务状态
	Config          string    `db:"config" json:"config"`                       // 任务配置（JSON格式）
	LastRunResult   string    `db:"last_run_result" json:"last_run_result"`     // 上次执行结果
	LastRunTime     time.Time `db:"last_run_time" json:"last_run_time"`         // 上次执行时间
	LastRunDuration int       `db:"last_run_duration" json:"last_run_duration"` // 上次执行耗时（毫秒）
	Description     string    `db:"description" json:"description"`             // 任务描述
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
