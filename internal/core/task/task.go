package task

import (
	"context"
	"fmt"
	"time"

	_ "github.com/bestruirui/bestsub/internal/core/task/handlers"
	"github.com/bestruirui/bestsub/internal/models/task"
)

// CreateTask 创建任务
func CreateTask(ctx context.Context, req *task.CreateRequest) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.createTask(ctx, req)
}

// GetTask 获取任务
func GetTask(ctx context.Context, id int64) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.repo.GetByID(ctx, id)
}

// UpdateTask 更新任务
func UpdateTask(ctx context.Context, id int64, req *task.UpdateRequest) (*task.Data, error) {
	if globalScheduler == nil {
		return nil, fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.updateTask(ctx, id, req)
}

// DeleteTask 删除任务
func DeleteTask(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.deleteTask(ctx, id)
}

// ListTasks 列出任务
func ListTasks(ctx context.Context, offset, limit int) (*[]task.Data, int64, error) {
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

// RunTask 手动运行任务
func RunTask(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.runTask(ctx, id)
}

// StopTask 停止任务
func StopTask(ctx context.Context, id int64) error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return globalScheduler.stopTask(ctx, id)
}

// GetTaskLogs 获取任务日志摘要
func GetTaskLogs(ctx context.Context, taskID int64, offset, limit int) (*[]LogFileInfo, int64, error) {
	return ReadTaskLogSummaries(taskID, offset, limit)
}

// GetTaskLogDetail 获取任务详细日志
func GetTaskLogDetail(ctx context.Context, taskID int64, logTime time.Time) (*[]TaskLog, error) {
	return ReadTaskLogDetail(taskID, logTime)
}
