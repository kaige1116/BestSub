package models

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`                     // 服务状态
	Timestamp string `json:"timestamp" example:"2024-01-01T12:00:00"` // 检查时间
	Version   string `json:"version" example:"1.0.0"`                 // 版本信息
	Database  string `json:"database" example:"connected"`            // 数据库状态
}
