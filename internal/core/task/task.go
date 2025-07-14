package task

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/bestruirui/bestsub/internal/core/task/exec"
	"github.com/bestruirui/bestsub/internal/core/task/scheduler"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var taskScheduler scheduler.Scheduler
var mu sync.RWMutex

// Initialize 初始化任务调度器
func Initialize() error {
	taskScheduler = scheduler.New()
	return nil
}

// Start 启动任务调度器
func Start() {
	taskScheduler.Start()
}

func AddTask(req *task.CreateRequest) error {
	mu.Lock()
	defer mu.Unlock()
	repo := database.Task()
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}
	t := &task.Data{
		BaseDbModel: common.BaseDbModel{
			Name:        req.Name,
			Enable:      enable,
			Description: req.Description,
		},
		Cron:     req.Cron,
		Type:     req.Type,
		LogLevel: req.LogLevel,
		Config:   req.Config,
		Timeout:  req.Timeout,
		Retry:    req.Retry,
	}
	err := repo.Create(context.Background(), t)
	if err != nil {
		return err
	}
	return taskScheduler.AddTask(t)
}

func UpdateTask(req *task.UpdateRequest) error {
	mu.Lock()
	defer mu.Unlock()
	repo := database.Task()
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}
	t := &task.Data{
		BaseDbModel: common.BaseDbModel{
			Name:        req.Name,
			Enable:      enable,
			ID:          req.ID,
			Description: req.Description,
		},
		Cron:     req.Cron,
		LogLevel: req.LogLevel,
		Config:   req.Config,
		Timeout:  req.Timeout,
		Retry:    req.Retry,
	}
	err := repo.Update(context.Background(), t)
	if err != nil {
		log.Errorf("failed to update task %d: %w", req.ID, err)
		return err
	}
	return taskScheduler.UpdateTask(t)
}
func RunTask(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	repo := database.Task()
	task, err := repo.GetByID(context.Background(), taskID)
	if err != nil {
		return err
	}
	if task.Enable {
		return taskScheduler.RunTask(taskID)
	} else {
		return fmt.Errorf("task %d is not enabled", taskID)
	}
}

func RemoveTask(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	repo := database.Task()
	err := repo.Delete(context.Background(), taskID)
	if err != nil {
		return err
	}
	return taskScheduler.RemoveTask(taskID)
}

func StopTask(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	return taskScheduler.StopTask(taskID)
}

func Shutdown() error {
	mu.Lock()
	defer mu.Unlock()
	return taskScheduler.Stop()
}
