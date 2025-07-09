package task

import (
	"context"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/core/task/handler"
	_ "github.com/bestruirui/bestsub/internal/core/task/handler"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// Create 创建任务
func Create(ctx context.Context, req *task.CreateRequest) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.createTask(ctx, req)
}

// Get 获取任务
func Get(ctx context.Context, id int64) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.repo.GetByID(ctx, id)
}

// Update 更新任务
func Update(ctx context.Context, req *task.UpdateRequest) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.updateTask(ctx, req)
}

// DeleteWithDb 删除任务
func DeleteWithDb(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.deleteTask(ctx, id)
}

// List 列出任务
func List(ctx context.Context, offset, limit int) (*[]task.Data, int64, error) {
	if globalScheduler == nil {
		return nil, 0, fmt.Errorf("scheduler not initialized")
	}

	tasks, err := globalScheduler.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := globalScheduler.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

// Run 手动运行任务
func Run(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.runTask(ctx, id)
}

// Stop 停止任务
func Stop(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.stopTask(ctx, id)
}

// GetLogs 获取任务日志摘要
func GetLogs(ctx context.Context, taskID int64, offset, limit int) (*[]LogFileInfo, int64, error) {
	return ReadTaskLogSummaries(taskID, offset, limit)
}

// GetLogDetail 获取任务详细日志
func GetLogDetail(ctx context.Context, taskID int64, logTime time.Time) (*[]TaskLog, error) {
	return ReadTaskLogDetail(taskID, logTime)
}

// RefreshScheduler 刷新调度器
// 重新同步数据库中的任务到调度器，移除已被删除的任务，添加新启用的任务
// 主要用于处理数据库触发器删除任务后的调度器同步问题
func RefreshScheduler(ctx context.Context) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.refreshScheduler(ctx)
}

// Delete 仅从调度器中删除任务，不操作数据库
// 用于配合数据库触发器使用，当触发器已经删除数据库记录时调用
func Delete(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	globalScheduler.mu.Lock()
	defer globalScheduler.mu.Unlock()

	// 从调度器移除
	if globalScheduler.started {
		if err := globalScheduler.removeTaskFromScheduler(id); err != nil {
			log.Errorf("Failed to remove task %d from scheduler: %v", id, err)
		}
	}

	// 停止正在运行的任务
	globalScheduler.runningTasks.Delete(id)

	log.Infof("从调度器删除任务: %d", id)
	return nil
}

// GetAllHandlers 返回所有handler的类型和对应的配置项结构，便于前端处理
func GetAllHandlers() []handler.Response {
	return handler.GetAll()
}

// GetHandlerTypes 获取所有已注册的任务类型
func GetHandlerTypes() []string {
	return handler.GetTypes()
}

// GetHandlerConfig 获取指定任务类型的配置信息
func GetHandlerConfig(taskType string) *handler.Response {
	config := handler.GetConfig(taskType)
	if config == nil {
		return nil
	}

	return &handler.Response{
		Type:   taskType,
		Config: config,
	}
}

// GetSchedulerStats 获取调度器统计信息
func GetSchedulerStats() map[string]any {
	if globalScheduler == nil {
		return map[string]any{
			"initialized": false,
		}
	}

	globalScheduler.mu.RLock()
	defer globalScheduler.mu.RUnlock()

	// 统计调度的任务数量
	scheduledCount := 0
	globalScheduler.scheduledTasks.Range(func(key, value any) bool {
		scheduledCount++
		return true
	})

	// 统计运行中的任务数量
	runningCount := 0
	globalScheduler.runningTasks.Range(func(key, value any) bool {
		runningCount++
		return true
	})

	return map[string]any{
		"initialized":     true,
		"started":         globalScheduler.started,
		"scheduled_tasks": scheduledCount,
		"running_tasks":   runningCount,
	}
}

// GetRunningTasks 获取正在运行的任务列表
func GetRunningTasks() []int64 {
	if globalScheduler == nil {
		return nil
	}

	var runningTasks []int64
	globalScheduler.runningTasks.Range(func(key, value any) bool {
		if taskID, ok := key.(int64); ok {
			runningTasks = append(runningTasks, taskID)
		}
		return true
	})

	return runningTasks
}

// ReloadTask 重新加载单个任务到调度器
func ReloadTask(ctx context.Context, taskID int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	globalScheduler.mu.Lock()
	defer globalScheduler.mu.Unlock()

	// 先移除现有的任务
	if err := globalScheduler.removeTaskFromScheduler(taskID); err != nil {
		log.Warnf("Failed to remove task %d from scheduler: %v", taskID, err)
	}

	// 重新加载任务
	taskData, err := globalScheduler.repo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task %d: %w", taskID, err)
	}

	if taskData != nil && taskData.Enable && globalScheduler.started {
		if err := globalScheduler.addTaskToScheduler(taskData); err != nil {
			return fmt.Errorf("failed to add task %d to scheduler: %w", taskID, err)
		}
	}

	log.Infof("重新加载任务: %d", taskID)
	return nil
}

// ReloadAllTasks 重新加载所有任务到调度器
func ReloadAllTasks(ctx context.Context) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	globalScheduler.mu.Lock()
	defer globalScheduler.mu.Unlock()

	// 清除所有现有的调度任务
	globalScheduler.scheduledTasks.Range(func(key, value any) bool {
		if taskID, ok := key.(int64); ok {
			globalScheduler.removeTaskFromScheduler(taskID)
		}
		return true
	})

	// 重新加载所有启用的任务
	if err := globalScheduler.loadEnabledTasks(ctx); err != nil {
		return fmt.Errorf("failed to reload tasks: %w", err)
	}

	log.Info("重新加载所有任务完成")
	return nil
}
