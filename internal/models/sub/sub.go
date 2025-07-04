package sub

import (
	"github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/task"
)

// Sub 订阅基础模型（数据库模型）
type Data struct {
	common.BaseDbModel
	URL string `json:"url"`
}

// CreateRequest 创建订阅链接请求模型
type CreateRequest struct {
	common.BaseRequestModel
	URL  string               `json:"url"`
	Task []task.CreateRequest `json:"task"`
}

// UpdateRequest 更新订阅链接请求模型
type UpdateRequest struct {
	common.BaseUpdateRequestModel
	URL  string               `json:"url"`
	Task []task.UpdateRequest `json:"task"`
}

// Response 订阅链接响应模型
type Response struct {
	Data
	Task []task.Data `json:"task"`
}
