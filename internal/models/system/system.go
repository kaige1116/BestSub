package system

import (
	"time"
)

// 管理系统的全局设置
type Data struct {
	ID          int64     `db:"id" json:"id"`
	GroupName   string    `db:"group_name" json:"group_name"`   // 配置分组：system, nodepool, gc, log, auth, api, proxy, monitor
	Key         string    `db:"key" json:"key"`                 // 配置键
	Type        string    `db:"type" json:"type"`               // 配置类型：string, int, bool, json
	Value       string    `db:"value" json:"value"`             // 配置值（JSON格式）
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// UpdateConfigItemRequest 更新配置项请求模型
type UpdateConfigItemRequest struct {
	Data []struct {
		ID          int64  `json:"id" example:"1"`               // 配置ID
		Value       string `json:"value" example:"true"`         // 配置值
		Description string `json:"description" example:"是否启用代理"` // 配置描述
	} `json:"data" binding:"required"` // 配置数据
}

// ConfigItemsResponse 配置项列表响应模型
type ConfigItemsResponse struct {
	Data  []Data `json:"data"`  // 配置项列表
	Total int    `json:"total"` // 总数
}
