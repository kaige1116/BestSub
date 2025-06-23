package models

import (
	"time"
)

// 订阅输出系统的分享链接功能
type ShareLink struct {
	ID               int64     `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`                             // 分享链接名称
	Description      string    `db:"description" json:"description"`               // 链接描述
	Token            string    `db:"token" json:"token"`                           // 随机复杂字符串（用于URL）
	IsEnabled        bool      `db:"is_enabled" json:"is_enabled"`                 // 是否启用
	OutputTemplateID int64     `db:"output_template_id" json:"output_template_id"` // 关联的输出模板ID
	NodeFilterID     int64     `db:"node_filter_id" json:"node_filter_id"`         // 关联的节点筛选配置ID
	ExpiresAt        time.Time `db:"expires_at" json:"expires_at"`                 // 过期时间
	MaxDownloads     int       `db:"max_downloads" json:"max_downloads"`           // 最大下载次数（0表示无限制）
	DownloadCount    int       `db:"download_count" json:"download_count"`         // 已下载次数
	LastAccess       time.Time `db:"last_access" json:"last_access"`               // 最后访问时间
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
