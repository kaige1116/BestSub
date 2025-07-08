package task

import (
	"context"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// createTask 创建任务的内部实现
func (s *Scheduler) createTask(ctx context.Context, req *task.CreateRequest) (*task.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建任务数据
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}

	taskData := &task.Data{
		BaseDbModel: common.BaseDbModel{
			Name:        req.Name,
			Description: req.Description,
			Enable:      enable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Cron:    req.Cron,
		Type:    req.Type,
		Config:  req.Config,
		Timeout: 60, // 默认超时60秒
		Retry:   3,  // 默认重试3次
		Status:  task.StatusPending,
	}

	// 保存到数据库
	if err := s.repo.Create(ctx, taskData); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// 如果任务启用，添加到调度器
	if taskData.Enable && s.isStarted() {
		if err := s.addTaskToScheduler(taskData); err != nil {
			log.Errorf("Failed to add task %d to scheduler: %v", taskData.ID, err)
		}
	}

	log.Infof("创建任务成功: %d (%s)", taskData.ID, taskData.Name)
	return taskData, nil
}

// updateTask 更新任务的内部实现
func (s *Scheduler) updateTask(ctx context.Context, req *task.UpdateRequest) (*task.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取现有任务
	existingTask, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// 从调度器移除旧任务
	if existingTask.Enable && s.isStarted() {
		if err := s.removeTaskFromScheduler(req.ID); err != nil {
			log.Errorf("Failed to remove task %d from scheduler: %v", req.ID, err)
		}
	}

	// 更新任务数据
	enable := existingTask.Enable
	if req.Enable != nil {
		enable = *req.Enable
	}

	taskData := &task.Data{
		BaseDbModel: common.BaseDbModel{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Enable:      enable,
			CreatedAt:   existingTask.CreatedAt,
			UpdatedAt:   time.Now(),
		},
		Cron:   req.Cron,
		Config: req.Config,
		// 保留其他字段
		Type:            existingTask.Type,
		IsSysTask:       existingTask.IsSysTask,
		Timeout:         existingTask.Timeout,
		Retry:           existingTask.Retry,
		SuccessCount:    existingTask.SuccessCount,
		FailedCount:     existingTask.FailedCount,
		LastRunResult:   existingTask.LastRunResult,
		LastRunTime:     existingTask.LastRunTime,
		LastRunDuration: existingTask.LastRunDuration,
	}

	// 更新数据库
	if err := s.repo.Update(ctx, taskData); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// 如果任务启用，重新添加到调度器
	if taskData.Enable && s.isStarted() {
		if err := s.addTaskToScheduler(taskData); err != nil {
			log.Errorf("Failed to reschedule task %d: %v", req.ID, err)
		}
	}

	log.Infof("更新任务成功: %d (%s)", taskData.ID, taskData.Name)
	return taskData, nil
}

// deleteTask 删除任务的内部实现
func (s *Scheduler) deleteTask(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 从调度器移除
	if s.isStarted() {
		if err := s.removeTaskFromScheduler(id); err != nil {
			log.Errorf("Failed to remove task %d from scheduler: %v", id, err)
		}
	}

	// 停止正在运行的任务
	s.runningTasks.Delete(id)

	// 从数据库删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	log.Infof("删除任务成功: %d", id)
	return nil
}

// runTask 手动运行任务的内部实现
func (s *Scheduler) runTask(ctx context.Context, id int64) error {
	// 获取任务数据
	taskData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if !taskData.Enable {
		return fmt.Errorf("task %d is disabled", id)
	}

	// 在新的goroutine中执行任务
	go s.executeTask(taskData)

	log.Infof("手动触发任务: %d (%s)", taskData.ID, taskData.Name)
	return nil
}

// stopTask 停止任务的内部实现
func (s *Scheduler) stopTask(ctx context.Context, id int64) error {
	// 检查任务是否在运行
	if _, running := s.runningTasks.Load(id); !running {
		return fmt.Errorf("task %d is not running", id)
	}

	// 标记任务为取消状态
	taskData, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	taskData.Status = task.StatusCancelled
	if err := s.repo.Update(ctx, taskData); err != nil {
		log.Errorf("Failed to update task %d status to cancelled: %v", id, err)
	}

	// 从运行任务列表中移除
	s.runningTasks.Delete(id)

	log.Infof("停止任务: %d (%s)", taskData.ID, taskData.Name)
	return nil
}
