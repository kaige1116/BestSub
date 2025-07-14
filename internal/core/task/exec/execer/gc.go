package execer

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/bestruirui/bestsub/internal/core/task/exec"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

const (
	bytesToMB = 1024 * 1024
)

// init 自动注册GC处理器
func init() {
	exec.Register(&exec.RegisterInfo{
		Type:    task.TypeGC,
		Handler: &GCExec{},
		Config:  &GCConfig{},
	})
}

// GCExec 垃圾回收任务处理器
type GCExec struct{}

// GCConfig 垃圾回收任务配置
type GCConfig struct {
	ForceGC bool `json:"force" default:"false" description:"是否强制执行垃圾回收，true时会调用debug.FreeOSMemory()"`
}

// Do 执行垃圾回收任务
func (h *GCExec) Do(ctx context.Context, logger *log.Logger, task *exec.TaskInfo) error {
	startTime := time.Now()

	var gcConfig GCConfig
	if len(task.Config) > 0 {
		if err := json.Unmarshal(task.Config, &gcConfig); err != nil {
			return fmt.Errorf("配置解析失败: %w", err)
		}
	}

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	if gcConfig.ForceGC {
		debug.FreeOSMemory()
	} else {
		runtime.GC()
	}

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	duration := time.Since(startTime)
	memFreed := float64(memBefore.Alloc-memAfter.Alloc) / bytesToMB

	logger.Infof("垃圾回收任务执行完成，耗时: %d ms, 释放内存: %.2f MB", duration.Milliseconds(), memFreed)

	return nil
}
