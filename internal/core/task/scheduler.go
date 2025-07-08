package task

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
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

// Stop 停止任务调度器
func Stop() error {
	if globalScheduler == nil {
		return nil
	}
	return globalScheduler.stop()
}

// start 启动调度器
func (s *Scheduler) start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted() {
		return fmt.Errorf("scheduler already started")
	}

	// 启动gocron调度器
	s.cron.Start()

	// 加载所有启用的任务
	if err := s.loadEnabledTasks(); err != nil {
		return fmt.Errorf("failed to load enabled tasks: %w", err)
	}

	log.Info("任务调度器启动成功")
	return nil
}

// stop 停止调度器
func (s *Scheduler) stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStarted() {
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

	log.Debugf("任务调度器停止成功")
	return nil
}

// loadEnabledTasks 加载所有启用的任务
func (s *Scheduler) loadEnabledTasks() error {
	ctx := context.Background()
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

	log.Infof("加载了 %d 个启用的任务", len(*tasks))
	return nil
}

// isStarted 检查调度器是否已启动
func (s *Scheduler) isStarted() bool {
	// 通过检查是否有调度的任务来判断调度器状态
	hasScheduledTasks := false
	s.scheduledTasks.Range(func(key, value any) bool {
		hasScheduledTasks = true
		return false // 找到一个就停止遍历
	})
	return hasScheduledTasks
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
	if jobInterface, exists := s.scheduledTasks.Load(taskID); exists {
		if job, ok := jobInterface.(gocron.Job); ok {
			if err := s.cron.RemoveJob(job.ID()); err != nil {
				return fmt.Errorf("failed to remove job for task %d: %w", taskID, err)
			}
			s.scheduledTasks.Delete(taskID)
			log.Debugf("任务 %d 已从调度器移除", taskID)
		}
	}
	return nil
}
