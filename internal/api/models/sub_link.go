package models

import (
	"github.com/bestruirui/bestsub/internal/models/sublink"
)

// SubLinkCreateRequest 创建订阅链接请求模型
type SubLinkCreateRequest struct {
	ID int64 `json:"id" example:"1"`
	sublink.BaseData
}

// SubLinkUpdateRequest 更新订阅链接请求模型
type SubLinkUpdateRequest struct {
	ID int64 `json:"id" example:"1"`
	sublink.BaseData
}

// SubLinkResponse 订阅链接响应模型
type SubLinkResponse = sublink.Data

// SubLinkListResponse 订阅链接列表响应模型
type SubLinkListResponse struct {
	Items []sublink.Data `json:"items"` // 订阅链接列表
	Total int            `json:"total"` // 总数
}
