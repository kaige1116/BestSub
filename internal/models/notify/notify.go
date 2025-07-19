package notify

import (
	"github.com/bestruirui/bestsub/internal/models/common"
)

// 全局管理通知渠道配置
type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Type   string `db:"type" json:"type"`     // 通知类型
	Config string `db:"config" json:"config"` // 通知配置
}

// CreateRequest 创建通知渠道请求模型
type CreateRequest struct {
	common.BaseRequestModel
	Type   string `json:"type" binding:"required"`   // 通知类型
	Config string `json:"config" binding:"required"` // 通知配置
}

// UpdateRequest 更新通知渠道请求模型
type UpdateRequest struct {
	common.BaseUpdateRequestModel
	Type   string `json:"type"`   // 通知类型
	Config string `json:"config"` // 通知配置
}

// 通知模板
type Template struct {
	ID          uint16 `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`               // 模板名称
	Description string `db:"description" json:"description"` // 模板描述
	Template    string `db:"templates" json:"templates"`     // 模板内容
}

// TemplateCreateRequest 创建通知模板请求模型
type TemplateCreateRequest struct {
	common.BaseRequestModel
	Template string `json:"templates" binding:"required"` // 模板内容
}

// TemplateUpdateRequest 更新通知模板请求模型
type TemplateUpdateRequest struct {
	common.BaseUpdateRequestModel
	Template string `json:"templates"` // 模板内容
}
