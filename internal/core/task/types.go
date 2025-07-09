package task

import (
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/go-co-op/gocron/v2"
)

// globalScheduler 全局调度器实例
var globalScheduler *Scheduler

// Scheduler 任务调度器
// 负责管理和执行定时任务，支持cron表达式调度
type Scheduler struct {
	cron           gocron.Scheduler          // gocron调度器实例
	repo           interfaces.TaskRepository // 任务数据仓库
	runningTasks   sync.Map                  // 任务ID -> 运行状态，记录正在运行的任务
	scheduledTasks sync.Map                  // 任务ID -> gocron.Job，记录已调度的任务
	mu             sync.RWMutex              // 读写锁，保护调度器状态
	started        bool                      // 调度器是否已启动
}

// TaskLog 任务日志结构
type TaskLog struct {
	TaskID      int64          `json:"task_id"`
	ExecutionID string         `json:"execution_id"`
	Timestamp   time.Time      `json:"timestamp"`
	Level       string         `json:"level"` // INFO, WARN, ERROR, STATUS
	Message     string         `json:"message"`
	Error       string         `json:"error,omitempty"`
	Progress    int            `json:"progress"` // 执行进度 0-100
	Status      string         `json:"status,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
}

// LogFileInfo 日志文件信息结构
type LogFileInfo struct {
	Time      time.Time `json:"time"`       // 执行时间
	IsSuccess bool      `json:"is_success"` // 是否成功
	TaskID    int64     `json:"task_id"`    // 任务ID
}
