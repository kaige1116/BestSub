package task

import (
	"context"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/task/handler"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/go-co-op/gocron/v2"
)

// Initialize 初始化任务调度器
func Initialize() error {
	if globalScheduler != nil {
		return fmt.Errorf("scheduler already initialized")
	}

	var err error
	globalScheduler, err = newScheduler()
	return err
}

// newScheduler 创建新的调度器实例
func newScheduler() (*Scheduler, error) {
	// 创建gocron调度器
	cronScheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create cron scheduler: %w", err)
	}

	scheduler := &Scheduler{
		cron: cronScheduler,
		repo: database.Task(),
	}

	log.Info("任务调度器初始化成功")
	return scheduler, nil
}

// Start 启动任务调度器
func Start() error {
	if globalScheduler == nil {
		return fmt.Errorf("scheduler not initialized, call Initialize() first")
	}
	return globalScheduler.start()
}

// Shutdown 停止任务调度器
func Shutdown() error {
	if globalScheduler == nil {
		return nil
	}
	return globalScheduler.stop()
}

// start 启动调度器
func (s *Scheduler) start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("scheduler already started")
	}

	// 启动gocron调度器
	s.cron.Start()
	s.started = true

	// 加载所有启用的任务
	if err := s.loadEnabledTasks(context.Background()); err != nil {
		s.started = false
		s.cron.Shutdown()
		return fmt.Errorf("failed to load enabled tasks: %w", err)
	}

	log.Info("任务调度器启动成功")
	return nil
}

// stop 停止调度器
func (s *Scheduler) stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	// 停止gocron调度器
	if err := s.cron.Shutdown(); err != nil {
		log.Errorf("Failed to shutdown cron scheduler: %v", err)
	}

	// 清除所有调度任务
	s.scheduledTasks.Range(func(key, value any) bool {
		s.scheduledTasks.Delete(key)
		return true
	})

	s.started = false
	log.Debugf("任务调度器停止成功")
	return nil
}

// loadEnabledTasks 加载所有启用的任务
func (s *Scheduler) loadEnabledTasks(ctx context.Context) error {
	tasks, err := s.repo.List(ctx, 0, 1000) // 加载前1000个任务
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	if tasks == nil {
		return nil
	}

	for _, taskData := range *tasks {
		if taskData.Enable {
			if err := s.addTaskToScheduler(&taskData); err != nil {
				log.Errorf("Failed to add task %d to scheduler: %v", taskData.ID, err)
				continue
			}
		}
	}

	log.Debugf("加载了 %d 个启用的任务", len(*tasks))
	return nil
}

// addTaskToScheduler 添加任务到调度器
func (s *Scheduler) addTaskToScheduler(taskData *task.Data) error {
	// 创建任务执行函数
	taskFunc := func() {
		s.executeTask(taskData)
	}

	// 创建cron作业
	job, err := s.cron.NewJob(
		gocron.CronJob(taskData.Cron, false),
		gocron.NewTask(taskFunc),
	)
	if err != nil {
		return fmt.Errorf("failed to create cron job for task %d: %w", taskData.ID, err)
	}

	// 存储作业引用
	s.scheduledTasks.Store(taskData.ID, job)
	log.Debugf("任务 %d (%s) 已添加到调度器", taskData.ID, taskData.Name)
	return nil
}

// removeTaskFromScheduler 从调度器移除任务
func (s *Scheduler) removeTaskFromScheduler(taskID int64) error {
	if job, exists := s.scheduledTasks.Load(taskID); exists {
		if cronJob, ok := job.(gocron.Job); ok {
			if err := s.cron.RemoveJob(cronJob.ID()); err != nil {
				return fmt.Errorf("failed to remove cron job: %w", err)
			}
		}
		s.scheduledTasks.Delete(taskID)

		// 同时清理运行任务列表
		s.runningTasks.Delete(taskID)

		log.Debugf("任务 %d 已从调度器移除", taskID)
	}
	return nil
}

// executeTask 执行任务（用于cron调度，使用默认上下文）
func (s *Scheduler) executeTask(taskData *task.Data) {
	// 为cron调度的任务创建带有任务信息的上下文
	ctx := context.Background()

	s.executeTaskWithContext(ctx, taskData)
}

// executeTaskWithContext 执行任务（支持传入上下文）
func (s *Scheduler) executeTaskWithContext(ctx context.Context, taskData *task.Data) {
	taskID := taskData.ID

	// 检查任务是否已在运行
	if _, running := s.runningTasks.Load(taskID); running {
		log.Warnf("任务 %d 已在运行中，跳过本次执行", taskID)
		return
	}

	// 标记任务开始运行
	s.runningTasks.Store(taskID, true)
	defer s.runningTasks.Delete(taskID)

	// 更新任务状态为运行中
	taskData.Status = task.StatusRunning
	if err := s.updateTaskStatus(ctx, taskData); err != nil {
		log.Errorf("Failed to update task %d status to running: %v", taskID, err)
	}

	startTime := time.Now()
	var success bool
	var resultMsg string

	// 执行任务逻辑
	success, resultMsg = s.executeTaskLogic(ctx, taskData)

	// 根据执行结果更新任务状态
	if success {
		taskData.Status = task.StatusCompleted
		taskData.SuccessCount++
		log.Infof("任务 %d (%s) 执行成功", taskID, taskData.Name)
	} else {
		taskData.Status = task.StatusFailed
		taskData.FailedCount++
		log.Errorf("任务 %d (%s) 执行失败: %s", taskID, taskData.Name, resultMsg)
	}

	// 更新任务执行结果
	duration := int(time.Since(startTime).Milliseconds())
	taskData.LastRunTime = &startTime
	taskData.LastRunDuration = &duration
	taskData.LastRunResult = resultMsg

	if err := s.updateTaskStatus(ctx, taskData); err != nil {
		log.Errorf("Failed to update task %d execution result: %v", taskID, err)
	}
}

// executeTaskLogic 执行任务的核心逻辑
func (s *Scheduler) executeTaskLogic(ctx context.Context, taskData *task.Data) (bool, string) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(taskData.Timeout)*time.Second)
	defer cancel()
	// 从注册表获取处理器
	taskHandler, exists := handler.Get(taskData.Type)
	if !exists {
		return false, fmt.Sprintf("未找到任务类型处理器: %s", taskData.Type)
	}

	// 验证配置
	if err := taskHandler.Validate(taskData.Config); err != nil {
		return false, fmt.Sprintf("任务配置无效: %v", err)
	}

	taskInfo := &handler.TaskInfo{
		Context: ctx,
		ID:      taskData.ID,
		Name:    taskData.Name,
		Type:    taskData.Type,
		Config:  taskData.Config,
	}

	if err := taskHandler.Execute(taskInfo); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, fmt.Sprintf("任务执行超时 (%d 秒): %v", taskData.Timeout, err)
		}
		return false, fmt.Sprintf("执行失败: %v", err)
	}

	return true, "执行成功"
}

// updateTaskStatus 更新任务状态到数据库
func (s *Scheduler) updateTaskStatus(ctx context.Context, taskData *task.Data) error {
	return s.repo.Update(ctx, taskData)
}

// createTask 创建任务的内部实现
func (s *Scheduler) createTask(ctx context.Context, req *task.CreateRequest) (*task.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建任务数据
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}
	taskConfig := config.GetTaskConfig()
	if req.Timeout > taskConfig.MaxTimeout {
		req.Timeout = taskConfig.MaxTimeout
		log.Warnf("任务 %s 的执行超时时间超过了最大值，已自动设置为 %d 秒", req.Name, taskConfig.MaxTimeout)
	}
	if req.Retry > taskConfig.MaxRetry {
		req.Retry = taskConfig.MaxRetry
		log.Warnf("任务 %s 的重试次数超过了最大值，已自动设置为 %d 次", req.Name, taskConfig.MaxRetry)
	}
	taskConfig = nil

	taskData := &task.Data{
		BaseDbModel: common.BaseDbModel{
			Name:        req.Name,
			Description: req.Description,
			Enable:      enable,
			CreatedAt:   timeutils.Now(),
			UpdatedAt:   timeutils.Now(),
		},
		Cron:    req.Cron,
		Type:    req.Type,
		Config:  req.Config,
		Timeout: req.Timeout,
		Retry:   req.Retry,
		Status:  task.StatusPending,
	}

	// 保存到数据库
	if err := s.repo.Create(ctx, taskData); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// 如果任务启用，添加到调度器
	if taskData.Enable && s.started {
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
	if s.started {
		if err := s.removeTaskFromScheduler(req.ID); err != nil {
			log.Errorf("Failed to remove task %d from scheduler: %v", req.ID, err)
		}
	}

	enable := existingTask.Enable
	if req.Enable != nil {
		enable = *req.Enable
	}
	taskConfig := config.GetTaskConfig()
	if req.Timeout > taskConfig.MaxTimeout {
		req.Timeout = taskConfig.MaxTimeout
		log.Warnf("任务 %s 的执行超时时间超过了最大值，已自动设置为 %d 秒", req.Name, taskConfig.MaxTimeout)
	}
	if req.Retry > taskConfig.MaxRetry {
		req.Retry = taskConfig.MaxRetry
		log.Warnf("任务 %s 的重试次数超过了最大值，已自动设置为 %d 次", req.Name, taskConfig.MaxRetry)
	}
	taskConfig = nil
	taskData := &task.Data{
		BaseDbModel: common.BaseDbModel{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Enable:      enable,
			CreatedAt:   existingTask.CreatedAt,
			UpdatedAt:   timeutils.Now(),
		},
		Cron:            req.Cron,
		Type:            existingTask.Type,
		Config:          req.Config,
		Timeout:         req.Timeout,
		Retry:           req.Retry,
		Status:          existingTask.Status,
		SuccessCount:    existingTask.SuccessCount,
		FailedCount:     existingTask.FailedCount,
		LastRunTime:     existingTask.LastRunTime,
		LastRunDuration: existingTask.LastRunDuration,
		LastRunResult:   existingTask.LastRunResult,
	}

	// 更新数据库
	if err := s.repo.Update(ctx, taskData); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// 如果任务启用，重新添加到调度器
	if taskData.Enable && s.started {
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
	if s.started {
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

	// 在新的goroutine中执行任务，传递上下文
	go s.executeTaskWithContext(ctx, taskData)

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

	log.Infof("停止任务: %d", id)
	return nil
}

// refreshScheduler 刷新调度器，重新加载数据库中的任务
// 用于处理数据库触发器删除任务后的调度器同步问题
func (s *Scheduler) refreshScheduler(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		log.Debug("调度器未启动，跳过刷新")
		return nil
	}

	log.Info("开始刷新调度器...")

	// 获取数据库中所有任务的ID
	tasks, err := s.repo.List(ctx, 0, 10000) // 获取所有任务
	if err != nil {
		return fmt.Errorf("failed to list tasks from database: %w", err)
	}

	// 创建数据库中存在的任务ID集合
	dbTaskIDs := make(map[int64]bool)
	if tasks != nil {
		for _, taskData := range *tasks {
			dbTaskIDs[taskData.ID] = taskData.Enable
		}
	}

	// 检查调度器中的任务，移除数据库中不存在的任务
	var toRemove []int64
	s.scheduledTasks.Range(func(key, value any) bool {
		taskID := key.(int64)
		if _, exists := dbTaskIDs[taskID]; !exists {
			// 数据库中不存在此任务，需要从调度器中移除
			toRemove = append(toRemove, taskID)
		}
		return true
	})

	// 移除不存在的任务（removeTaskFromScheduler会自动清理运行任务列表）
	for _, taskID := range toRemove {
		if err := s.removeTaskFromScheduler(taskID); err != nil {
			log.Errorf("Failed to remove task %d from scheduler during refresh: %v", taskID, err)
		} else {
			log.Infof("已从调度器移除被删除的任务: %d", taskID)
		}
	}

	// 检查数据库中启用的任务，添加到调度器中（如果尚未添加）
	if tasks != nil {
		for _, taskData := range *tasks {
			if taskData.Enable {
				// 检查任务是否已在调度器中
				if _, exists := s.scheduledTasks.Load(taskData.ID); !exists {
					// 任务不在调度器中，需要添加
					if err := s.addTaskToScheduler(&taskData); err != nil {
						log.Errorf("Failed to add task %d to scheduler during refresh: %v", taskData.ID, err)
					} else {
						log.Infof("已向调度器添加新启用的任务: %d (%s)", taskData.ID, taskData.Name)
					}
				}
			}
		}
	}

	log.Infof("调度器刷新完成，移除了 %d 个已删除的任务", len(toRemove))
	return nil
}
