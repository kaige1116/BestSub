package models

import (
	"time"
)

// 管理系统的全局设置
type SystemConfig struct {
	ID          int64     `db:"id" json:"id"`
	Key         string    `db:"key" json:"key"`                 // 配置键
	Value       string    `db:"value" json:"value"`             // 配置值（JSON格式）
	Type        string    `db:"type" json:"type"`               // 配置类型：string, int, bool, json
	Group       string    `db:"group" json:"group"`             // 配置分组：system, nodepool, gc, log, auth, api, proxy, monitor
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

