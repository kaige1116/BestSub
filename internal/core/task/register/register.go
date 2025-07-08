package register

import (
	"context"
	"sync"
)

// TaskHandler 任务处理器接口
type TaskHandler interface {
	Execute(ctx context.Context, config string) error
	Validate(config string) error
}

// globalHandlers 全局处理器注册表
var (
	globalHandlers = make(map[string]TaskHandler)
	mu             sync.RWMutex
)

// AddHandler 注册处理器
func AddHandler(taskType string, handler TaskHandler) {
	mu.Lock()
	defer mu.Unlock()

	globalHandlers[taskType] = handler
}

// getHandler 获取指定类型的处理器
func GetHandler(taskType string) (TaskHandler, bool) {
	mu.RLock()
	defer mu.RUnlock()

	handler, exists := globalHandlers[taskType]
	return handler, exists
}
