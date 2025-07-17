package system

// 管理系统的全局设置
type Data struct {
	Key         string `db:"key" json:"key"`                 // 配置键
	GroupName   string `db:"group_name" json:"group_name"`   // 配置分组：system, nodepool, gc, log, auth, api, proxy, monitor
	Value       string `db:"value" json:"value"`             // 配置值（JSON格式）
	Description string `db:"description" json:"description"` // 配置描述
}

// UpdateConfigItemRequest 更新配置项请求模型
type UpdateConfigItemRequest struct {
	Data []struct {
		Value       string `json:"value" example:"true"`         // 配置值
		Description string `json:"description" example:"是否启用代理"` // 配置描述
	} `json:"data" binding:"required"` // 配置数据
}

// ConfigItemsResponse 配置项列表响应模型
type ConfigItemsResponse struct {
	Data  []Data `json:"data"`  // 配置项列表
	Total int    `json:"total"` // 总数
}
