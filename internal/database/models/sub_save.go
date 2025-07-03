package models

import (
	"time"
)

// 订阅输出系统的保存功能配置
type SubSaveConfig struct {
	ID          int64     `db:"id" json:"id"`
	Enable      bool      `db:"enable" json:"enable"`           // 是否启用
	Name        string    `db:"name" json:"name"`               // 保存配置名称
	Rename      string    `db:"rename" json:"rename"`           // 保存文件名
	Type        string    `db:"type" json:"type"`               // 保存类型：file, dir
	FileName    string    `db:"file_name" json:"file_name"`     // 保存文件名
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
