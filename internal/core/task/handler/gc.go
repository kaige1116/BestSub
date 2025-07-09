package handler

import (
	"encoding/json"
	"runtime"
	"runtime/debug"

	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// init 自动注册GC处理器
func init() {
	Register(&RegistryInfo{
		Type:    task.TypeGC,
		Handler: &GCHandler{},
		Config:  &GCConfig{},
	})
}

// GCHandler 垃圾回收任务处理器
type GCHandler struct{}

// GCConfig 垃圾回收任务配置
type GCConfig struct {
	ForceGC bool `json:"force" default:"false" description:"是否强制执行垃圾回收，true时会调用debug.FreeOSMemory()"`
}

// Execute 执行垃圾回收任务
func (h *GCHandler) Execute(taskInfo *TaskInfo) error {
	var gcConfig GCConfig
	if taskInfo.Config != "" {
		if err := json.Unmarshal([]byte(taskInfo.Config), &gcConfig); err != nil {
			return err
		}
	}

	log.Infof("开始执行任务 %d (%s) 的 GC 操作", taskInfo.ID, taskInfo.Name)

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	if gcConfig.ForceGC {
		debug.FreeOSMemory()
		log.Infof("任务 %d (%s) 强制GC完成", taskInfo.ID, taskInfo.Name)
	} else {
		runtime.GC()
		log.Infof("任务 %d (%s) GC完成", taskInfo.ID, taskInfo.Name)
	}

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	log.Infof("任务 %d (%s) GC完成，释放内存: %.2f MB", taskInfo.ID, taskInfo.Name,
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
