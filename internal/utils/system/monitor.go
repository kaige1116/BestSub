package system

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/log"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

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

// 系统监控器
type Monitor struct {
	startTime     time.Time
	uploadBytes   uint64
	downloadBytes uint64
	mu            sync.RWMutex
}

var monitor Monitor

func init() {
	monitor = Monitor{
		startTime: timeutils.Now(),
	}
	log.Debug("System monitor initialized")
}

// 增加上传流量
func AddUploadBytes(bytes uint64) {
	atomic.AddUint64(&monitor.uploadBytes, bytes)
}

// 增加下载流量
func AddDownloadBytes(bytes uint64) {
	atomic.AddUint64(&monitor.downloadBytes, bytes)
}

// 获取系统信息
func GetSystemInfo() *Info {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	now := timeutils.Now()
	uptime := now.Sub(monitor.startTime)

	uploadBytes := atomic.LoadUint64(&monitor.uploadBytes)
	downloadBytes := atomic.LoadUint64(&monitor.downloadBytes)

	var memoryPercent float64
	if memStats.Sys > 0 {
		memoryPercent = float64(memStats.Alloc) / float64(memStats.Sys) * 100
	}

	var lastGCTime string
	if memStats.LastGC > 0 {
		lastGCTime = time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339)
	} else {
		lastGCTime = "never"
	}

	return &Info{
		MemoryUsed:    memStats.Alloc,
		MemoryTotal:   memStats.Sys,
		MemoryPercent: memoryPercent,
		HeapUsed:      memStats.HeapAlloc,
		HeapTotal:     memStats.HeapSys,

		UptimeSeconds: int64(uptime.Seconds()),
		StartTime:     monitor.startTime.Format(time.RFC3339),

		CPUCores:   runtime.NumCPU(),
		Goroutines: runtime.NumGoroutine(),

		UploadBytes:   uploadBytes,
		DownloadBytes: downloadBytes,

		GCCount:    memStats.NumGC,
		LastGCTime: lastGCTime,
	}
}

// 重置监控数据
func Reset() {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()

	atomic.StoreUint64(&monitor.uploadBytes, 0)
	atomic.StoreUint64(&monitor.downloadBytes, 0)
	monitor.startTime = timeutils.Now()

	log.Debug("System monitor data reset")
}
