package models

type ConfigItemData struct {
	ID          int64  `json:"id" example:"1"`               // 配置ID
	Value       string `json:"value" example:"true"`         // 配置值
	Description string `json:"description" example:"是否启用代理"` // 配置描述
}

// ConfigItemResponse 配置项响应模型
type ConfigItemResponse struct {
	ConfigItemData
	GroupName string `json:"group_name" example:"proxy"`                // 配置分组
	Key       string `json:"key" example:"proxy.enable"`                // 配置键
	Type      string `json:"type" example:"bool"`                       // 配置类型
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"` // 创建时间
	UpdatedAt string `json:"updated_at" example:"2024-01-01T12:00:00Z"` // 更新时间
}

// ConfigItemsResponse 配置项列表响应模型
type ConfigItemsResponse struct {
	Items []ConfigItemResponse `json:"items"` // 配置项列表
	Total int                  `json:"total"` // 总数
}

// UpdateConfigItemRequest 更新配置项请求模型
type UpdateConfigItemRequest struct {
	Data []ConfigItemData `json:"data" binding:"required"` // 配置数据
}
