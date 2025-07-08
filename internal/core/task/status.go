package task

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

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
		"started":         globalScheduler.isStarted(),
		"scheduled_tasks": scheduledCount,
		"running_tasks":   runningCount,
	}
}

// IsRunning 检查调度器是否正在运行
func IsRunning() bool {
	if globalScheduler == nil {
		return false
	}
	globalScheduler.mu.RLock()
	defer globalScheduler.mu.RUnlock()
	return globalScheduler.isStarted()
}

// GetRunningTasks 获取正在运行的任务ID列表
func GetRunningTasks() []int64 {
	if globalScheduler == nil {
		return []int64{}
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
func ReloadTask(taskID int64) error {
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
	taskData, err := globalScheduler.repo.GetByID(context.Background(), taskID)
	if err != nil {
		return fmt.Errorf("failed to get task %d: %w", taskID, err)
	}

	if taskData != nil && taskData.Enable && globalScheduler.isStarted() {
		if err := globalScheduler.addTaskToScheduler(taskData); err != nil {
			return fmt.Errorf("failed to add task %d to scheduler: %w", taskID, err)
		}
	}

	log.Infof("重新加载任务: %d", taskID)
	return nil
}

// ReloadAllTasks 重新加载所有任务到调度器
func ReloadAllTasks() error {
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
	if err := globalScheduler.loadEnabledTasks(); err != nil {
		return fmt.Errorf("failed to reload tasks: %w", err)
	}

	log.Info("重新加载所有任务完成")
	return nil
}
