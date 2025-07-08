package task

import (
	"context"
	"fmt"
	"time"

	_ "github.com/bestruirui/bestsub/internal/core/task/handlers"
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

// Shutdown 停止任务
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
	return globalScheduler.refreshScheduler()
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
	if globalScheduler.isStarted() {
		if err := globalScheduler.removeTaskFromScheduler(id); err != nil {
			log.Errorf("Failed to remove task %d from scheduler: %v", id, err)
		}
	}

	// 停止正在运行的任务
	globalScheduler.runningTasks.Delete(id)

	log.Infof("从调度器删除任务: %d", id)
	return nil
}
