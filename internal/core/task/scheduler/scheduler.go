package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/core/task/exec"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler interface {
	Start()
	AddTask(task *task.Data) error
	UpdateTask(task *task.Data) error
	RunTask(taskID int64) error
	RemoveTask(taskID int64) error
	StopTask(taskID int64) error
	Stop() error
}

type scheduler struct {
	Cron           gocron.Scheduler
	Ctx            context.Context
	Cancel         context.CancelFunc
	RunningTasks   map[int64]context.CancelFunc
	ScheduledTasks map[int64]gocron.Job
}

func New() Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Error("failed to create cron scheduler: %w", err)
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &scheduler{
		Cron:           s,
		Ctx:            ctx,
		Cancel:         cancel,
		RunningTasks:   make(map[int64]context.CancelFunc),
		ScheduledTasks: make(map[int64]gocron.Job),
	}
}

func (s *scheduler) Start() {
	s.Cron.Start()
	log.Debug("任务调度器启动成功")
}

func (s *scheduler) Stop() error {
	s.Cancel()
	if err := s.Cron.Shutdown(); err != nil {
		log.Error("Failed to shutdown cron scheduler: %v", err)
	}
	log.Debug("任务调度器停止成功")
	return nil
}

func (s *scheduler) AddTask(taskData *task.Data) error {
	job, err := s.Cron.NewJob(
		gocron.CronJob(taskData.Cron, false),
		gocron.NewTask(s.execTask, taskData),
	)
	if err != nil {
		log.Errorf("failed to create cron job for task %d: %v", taskData.ID, err)
		return fmt.Errorf("failed to create cron job for task %d: %w", taskData.ID, err)
	}
	s.ScheduledTasks[taskData.ID] = job
	nextRun, err := job.NextRun()
	if err != nil {
		log.Errorf("failed to get next run time for task %d: %v", taskData.ID, err)
	}
	log.Debugf("task %d added next run at %s", taskData.ID, nextRun.Format(time.RFC3339))
	return nil
}

func (s *scheduler) UpdateTask(taskData *task.Data) error {
	job, ok := s.ScheduledTasks[taskData.ID]
	if !ok {
		log.Errorf("task %d not found", taskData.ID)
		return nil
	}
	job, err := s.Cron.Update(
		job.ID(),
		gocron.CronJob(taskData.Cron, false),
		gocron.NewTask(s.execTask, *taskData),
	)
	if err != nil {
		log.Errorf("failed to update cron job for task %d: %w", taskData.ID, err)
		return err
	}
	s.ScheduledTasks[taskData.ID] = job
	nextRun, err := job.NextRun()
	if err != nil {
		log.Errorf("failed to get next run time for task %d: %v", taskData.ID, err)
	}
	log.Debugf("task %d updated next run at %s", taskData.ID, nextRun.Format(time.RFC3339))
	return nil
}

func (s *scheduler) RunTask(taskID int64) error {
	job, ok := s.ScheduledTasks[taskID]
	if !ok {
		return fmt.Errorf("task %d not found", taskID)
	}
	err := job.RunNow()
	if err != nil {
		log.Errorf("failed to run task %d: %v", taskID, err)
		return err
	}
	log.Debugf("task %d run now", taskID)
	return nil
}

func (s *scheduler) RemoveTask(taskID int64) error {
	job, ok := s.ScheduledTasks[taskID]
	delete(s.ScheduledTasks, taskID)
	if !ok {
		log.Errorf("task %d not found", taskID)
		return nil
	}
	cancel, ok := s.RunningTasks[taskID]
	if ok {
		cancel()
		delete(s.RunningTasks, taskID)
	}
	err := s.Cron.RemoveJob(job.ID())
	if err != nil {
		log.Errorf("remove task %d failed: %v", taskID, err)
		return err
	}
	log.Debugf("task %d removed", taskID)
	return nil
}

func (s *scheduler) StopTask(taskID int64) error {
	cancel, ok := s.RunningTasks[taskID]
	defer delete(s.RunningTasks, taskID)
	if !ok {
		log.Errorf("task %d not found", taskID)
		return nil
	}
	cancel()
	return nil
}

func (s *scheduler) execTask(taskData *task.Data) {
	taskctx, cancel := context.WithCancel(s.Ctx)
	defer cancel()
	s.RunningTasks[taskData.ID] = cancel
	defer delete(s.RunningTasks, taskData.ID)
	taskInfo := &exec.TaskInfo{
		ID:     taskData.ID,
		Type:   taskData.Type,
		Config: []byte(taskData.Config),
	}

	for i := 0; i < taskData.Retry; i++ {
		select {
		case <-taskctx.Done():
			log.Infof("任务 %d 取消", taskData.ID)
			return
		default:
			log.Debugf("task %d running %d/%d", taskData.ID, i+1, taskData.Retry)
			err := exec.Run(taskctx, taskInfo)
			if err != nil {
				log.Errorf("任务 %d 执行失败: %v (重试 %d/%d)", taskData.ID, err, i+1, taskData.Retry)
				continue
			}
			log.Infof("任务 %d 执行成功", taskData.ID)
			return
		}
	}
}
