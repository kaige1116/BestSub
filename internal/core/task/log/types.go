package log

import (
	"sync"
	"time"
)

// LogEntry 日志条目 - JSON格式，便于前端解析
type LogEntry struct {
	Time     string `json:"time"`     // ISO 8601格式时间戳
	Level    string `json:"level"`    // 日志等级：INFO, WARN, ERROR, DEBUG
	Progress int    `json:"progress"` // 执行进度 0-100
	Message  string `json:"message"`  // 日志信息
}

// LogWriter 日志写入器 - 在内存中累积日志，最后一次性写入
type LogWriter struct {
	taskID   int64
	logs     []LogEntry    // 内存中的日志缓存
	mu       sync.Mutex    // 保护并发访问
	filePath string        // 日志文件完整路径
	maxLogs  int           // 最大日志数量限制（如1000）
	logChan  chan LogEntry // 实时日志流通道
}

// SessionInfo 会话信息 - 从文件名解析得到
type SessionInfo struct {
	TaskID   int64     `json:"task_id"`   // 任务ID
	Time     time.Time `json:"time"`      // 执行时间
	Status   bool      `json:"status"`    // 执行状态：false=失败, true=成功
	FileSize int64     `json:"file_size"` // 文件大小
}
