package handlers

import (
	"context"
	"encoding/json"
	"runtime"
	"runtime/debug"

	"github.com/bestruirui/bestsub/internal/core/task/register"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// init 自动注册GC处理器
func init() {
	register.AddHandler(task.TypeGC, &GCHandler{})
}

// GCHandler 垃圾回收任务处理器
type GCHandler struct{}

// GCConfig 垃圾回收任务配置
type GCConfig struct {
	ForceGC bool `json:"force_gc"`
}

// Execute 执行垃圾回收任务
func (h *GCHandler) Execute(ctx context.Context, config string) error {
	var gcConfig GCConfig
	if config != "" {
		if err := json.Unmarshal([]byte(config), &gcConfig); err != nil {
			return err
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

	log.Infof("GC完成，释放内存: %.2f MB",
		float64(memBefore.Alloc-memAfter.Alloc)/1024/1024)

	return nil
}

// Validate 验证配置
func (h *GCHandler) Validate(config string) error {
	if config == "" {
		return nil // 空配置是有效的
	}

	var gcConfig GCConfig
	return json.Unmarshal([]byte(config), &gcConfig)
}
