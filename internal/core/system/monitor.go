package system

import (
	"context"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/log"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

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
func GetSystemInfo() *system.Info {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	now := timeutils.Now()
	uptime := now.Sub(monitor.startTime)

	uploadBytes := atomic.LoadUint64(&monitor.uploadBytes)
	downloadBytes := atomic.LoadUint64(&monitor.downloadBytes)

	systemMemory := getSystemMemoryInfo()

	processMemory := getProcessMemoryInfo()

	var lastGCTime string
	if memStats.LastGC > 0 {
		lastGCTime = time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339)
	} else {
		lastGCTime = "never"
	}

	return &system.Info{
		MemoryUsed:    processMemory.used,
		MemoryTotal:   systemMemory.total,
		MemoryPercent: processMemory.percent,
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

// 内存信息结构
type memoryInfo struct {
	used    uint64
	total   uint64
	percent float64
}

// 获取系统内存信息
func getSystemMemoryInfo() memoryInfo {
	ctx := context.Background()
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		log.Debugf("Failed to get system memory info: %v", err)
		return memoryInfo{}
	}

	return memoryInfo{
		used:    vmStat.Used,
		total:   vmStat.Total,
		percent: vmStat.UsedPercent,
	}
}

// 获取进程内存信息
func getProcessMemoryInfo() memoryInfo {
	ctx := context.Background()
	pid := int32(os.Getpid())

	proc, err := process.NewProcessWithContext(ctx, pid)
	if err != nil {
		log.Debugf("Failed to create process instance: %v", err)
		return memoryInfo{}
	}

	memInfo, err := proc.MemoryInfoWithContext(ctx)
	if err != nil {
		log.Debugf("Failed to get process memory info: %v", err)
		return memoryInfo{}
	}

	systemMem := getSystemMemoryInfo()
	var percent float64
	if systemMem.total > 0 {
		percent = float64(memInfo.RSS) / float64(systemMem.total) * 100
	}

	return memoryInfo{
		used:    memInfo.RSS,
		total:   systemMem.total,
		percent: percent,
	}
}
