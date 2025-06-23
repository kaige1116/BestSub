package models

import (
	"time"
)

// 订阅输出系统的保存功能配置
type SubSaveConfig struct {
	ID               int64     `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`                             // 保存配置名称
	Description      string    `db:"description" json:"description"`               // 配置描述
	IsEnabled        bool      `db:"is_enabled" json:"is_enabled"`                 // 是否启用
	StorageID        int64     `db:"storage_id" json:"storage_id"`                 // 关联的存储配置ID
	OutputTemplateID int64     `db:"output_template_id" json:"output_template_id"` // 关联的输出模板ID
	NodeFilterID     int64     `db:"node_filter_id" json:"node_filter_id"`         // 关联的节点筛选配置ID
	FileName         string    `db:"file_name" json:"file_name"`                   // 保存文件名
	SaveInterval     int       `db:"save_interval" json:"save_interval"`           // 保存间隔（秒）
	LastSave         time.Time `db:"last_save" json:"last_save"`                   // 最后保存时间
	LastStatus       string    `db:"last_status" json:"last_status"`               // 最后保存状态
	ErrorMsg         string    `db:"error_msg" json:"error_msg"`                   // 错误信息
	SaveCount        int       `db:"save_count" json:"save_count"`                 // 保存次数统计
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
