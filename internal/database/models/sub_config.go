package models

import (
	"time"
)

// 全局管理保存方式配置，供订阅输出系统引用
type SubStorageConfig struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`               // 存储配置名称
	Type       string    `db:"type" json:"type"`               // 存储类型：webdav, local, ftp, sftp, s3, oss
	Config     string    `db:"config" json:"config"`           // 存储配置（JSON格式）
	IsActive   bool      `db:"is_active" json:"is_active"`     // 是否启用
	TestResult string    `db:"test_result" json:"test_result"` // 测试连接结果
	LastTest   time.Time `db:"last_test" json:"last_test"`     // 最后测试时间
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// 全局输出模板配置
type SubOutputTemplate struct {
	ID          int64     `db:"id" json:"id"`
	Format      string    `db:"format" json:"format"`           // 格式名称：mihomo, singbox, v2ray, clash
	Version     string    `db:"version" json:"version"`         // 模板版本
	Template    string    `db:"template" json:"template"`       // 模板内容
	Description string    `db:"description" json:"description"` // 模板描述
	IsDefault   bool      `db:"is_default" json:"is_default"`   // 是否为默认模板
	IsActive    bool      `db:"is_active" json:"is_active"`     // 是否启用
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// 全局节点输出筛选规则
type SubNodeFilterRule struct {
	ID        int64     `db:"id" json:"id"`
	RuleType  string    `db:"rule_type" json:"rule_type"`   // 规则类型：link, detector, speed, country, protocol
	Operator  string    `db:"operator" json:"operator"`     // 操作符：in, not_in, gt, lt, eq, ne, contains
	Value     string    `db:"value" json:"value"`           // 规则值（JSON格式）
	IsEnabled bool      `db:"is_enabled" json:"is_enabled"` // 是否启用
	Priority  int       `db:"priority" json:"priority"`     // 优先级
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
