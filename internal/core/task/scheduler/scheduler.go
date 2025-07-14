package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/core/task/exec"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler interface {
	Start()
	AddTask(cron string, function any, taskInfo *exec.TaskInfo) error
	AddRunTaskMap(taskID int64, cancel context.CancelFunc)
	RemoveRunTaskMap(taskID int64)
	UpdateTask(cron string, function any, taskInfo *exec.TaskInfo) error
	RunTask(taskID int64) error
	RemoveTask(taskID int64) error
	StopTask(taskID int64) error
	Stop() error
}

type scheduler struct {
	Cron           gocron.Scheduler
	Ctx            context.Context
	Cancel         context.CancelFunc
	RunningTasks   sync.Map
	ScheduledTasks sync.Map
}

func New() (Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Error("failed to create cron scheduler: %w", err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &scheduler{
		Cron:           s,
		Ctx:            ctx,
		Cancel:         cancel,
		RunningTasks:   sync.Map{},
		ScheduledTasks: sync.Map{},
	}, nil
}

func (s *scheduler) Start() {
	s.Cron.Start()
	log.Debug("任务调度器启动成功")
}

func (s *scheduler) Stop() error {
	s.Cancel()
	if err := s.Cron.Shutdown(); err != nil {
		log.Error("Failed to shutdown cron scheduler: %v", err)
		return err
	}
	log.Debug("任务调度器停止成功")
	return nil
}

func (s *scheduler) AddTask(cron string, function any, taskInfo *exec.TaskInfo) error {
	job, err := s.Cron.NewJob(
		gocron.CronJob(cron, false),
		gocron.NewTask(function, s.Ctx, *taskInfo),
	)
	if err != nil {
		log.Errorf("failed to create cron job for task %d: %v", taskInfo.ID, err)
		return err
	}
	s.ScheduledTasks.Store(taskInfo.ID, job)
	nextRun, err := job.NextRun()
	if err != nil {
		log.Errorf("failed to get next run time for task %d: %v", taskInfo.ID, err)
		return err
	}
	log.Debugf("task %d added next run at %s", taskInfo.ID, nextRun.Format(time.RFC3339))
	return nil
}

func (s *scheduler) AddRunTaskMap(taskID int64, cancel context.CancelFunc) {
	s.RunningTasks.Store(taskID, cancel)
}

func (s *scheduler) RemoveRunTaskMap(taskID int64) {
	s.RunningTasks.Delete(taskID)
}

func (s *scheduler) UpdateTask(cron string, function any, taskInfo *exec.TaskInfo) error {
	value, ok := s.ScheduledTasks.Load(taskInfo.ID)
	if !ok {
		log.Errorf("task %d not found", taskInfo.ID)
		return fmt.Errorf("task %d not found", taskInfo.ID)
	}
	job := value.(gocron.Job)
	newJob, err := s.Cron.Update(
		job.ID(),
		gocron.CronJob(cron, false),
		gocron.NewTask(function, s.Ctx, *taskInfo),
	)
	if err != nil {
		log.Errorf("failed to update cron job for task %d: %w", taskInfo.ID, err)
		return err
	}
	s.ScheduledTasks.Store(taskInfo.ID, newJob)
	nextRun, err := newJob.NextRun()
	if err != nil {
		log.Errorf("failed to get next run time for task %d: %v", taskInfo.ID, err)
		return err
	}
	log.Debugf("task %d updated next run at %s", taskInfo.ID, nextRun.Format(time.RFC3339))
	return nil
}

func (s *scheduler) RunTask(taskID int64) error {
	value, ok := s.ScheduledTasks.Load(taskID)
	if !ok {
		log.Errorf("task %d not found", taskID)
		return fmt.Errorf("task %d not found", taskID)
	}
	job := value.(gocron.Job)
	err := job.RunNow()
	if err != nil {
		log.Errorf("failed to run task %d: %v", taskID, err)
		return err
	}
	log.Debugf("task %d run now", taskID)
	return nil
}

func (s *scheduler) RemoveTask(taskID int64) error {
	value, ok := s.ScheduledTasks.Load(taskID)
	if !ok {
		log.Errorf("task %d not found", taskID)
		return fmt.Errorf("task %d not found", taskID)
	}
	job := value.(gocron.Job)
	cancel, ok := s.RunningTasks.Load(taskID)
	if ok {
		cancel.(context.CancelFunc)()
		s.RunningTasks.Delete(taskID)
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
	value, ok := s.RunningTasks.Load(taskID)
	if !ok {
		log.Errorf("task %d not found", taskID)
		return fmt.Errorf("task %d not found", taskID)
	}
	cancel := value.(context.CancelFunc)
	cancel()
	s.RunningTasks.Delete(taskID)
	return nil
}
