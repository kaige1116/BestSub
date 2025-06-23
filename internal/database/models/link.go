package models

import (
	"time"
)

// 管理订阅链接的基本信息和配置
type Link struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`               // 链接名称
	URL        string    `db:"url" json:"url"`                 // 订阅链接地址
	UserAgent  string    `db:"user_agent" json:"user_agent"`   // 自定义 User-Agent
	IsEnabled  bool      `db:"is_enabled" json:"is_enabled"`   // 启用/禁用状态
	UseProxy   bool      `db:"use_proxy" json:"use_proxy"`     // 是否启用代理更新订阅
	CronExpr   string    `db:"cron_expr" json:"cron_expr"`     // 定时更新设置（cron表达式）
	LastUpdate time.Time `db:"last_update" json:"last_update"` // 最后更新时间
	LastStatus string    `db:"last_status" json:"last_status"` // 最后更新状态
	ErrorMsg   string    `db:"error_msg" json:"error_msg"`     // 错误信息
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// 每个链接的每个功能模块都有独立的配置实例
type ModuleConfig struct {
	ID         int64     `db:"id" json:"id"`
	LinkID     int64     `db:"link_id" json:"link_id"`
	ModuleType string    `db:"module_type" json:"module_type"` // parser, detector, notifier, namer
	ModuleName string    `db:"module_name" json:"module_name"` // 具体模块名称
	IsEnabled  bool      `db:"is_enabled" json:"is_enabled"`   // 启用/禁用状态
	Priority   int       `db:"priority" json:"priority"`       // 优先级（检测器使用）
	Config     string    `db:"config" json:"config"`           // JSON格式的配置参数
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
