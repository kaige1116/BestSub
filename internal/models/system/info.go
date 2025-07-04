package system

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`                     // 服务状态
	Timestamp string `json:"timestamp" example:"2024-01-01T12:00:00"` // 检查时间
	Version   string `json:"version" example:"1.0.0"`                 // 版本信息
	Database  string `json:"database" example:"connected"`            // 数据库状态
}

// 系统信息结构
type Info struct {
	MemoryUsed    uint64  `json:"memory_used"`    // 已使用内存 (bytes)
	MemoryTotal   uint64  `json:"memory_total"`   // 总内存 (bytes)
	MemoryPercent float64 `json:"memory_percent"` // 内存使用百分比
	HeapUsed      uint64  `json:"heap_used"`      // 堆内存使用 (bytes)
	HeapTotal     uint64  `json:"heap_total"`     // 堆内存总量 (bytes)
	UptimeSeconds int64   `json:"uptime_seconds"` // 运行时长(秒)
	StartTime     string  `json:"start_time"`     // 启动时间
	CPUCores      int     `json:"cpu_cores"`      // CPU核心数
	Goroutines    int     `json:"goroutines"`     // 协程数量
	UploadBytes   uint64  `json:"upload_bytes"`   // 上传流量 (bytes)
	DownloadBytes uint64  `json:"download_bytes"` // 下载流量 (bytes)
	GCCount       uint32  `json:"gc_count"`       // GC次数
	LastGCTime    string  `json:"last_gc_time"`   // 最后GC时间
}
