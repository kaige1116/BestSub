package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/bestruirui/bestsub/internal/core/task/exec"
	_ "github.com/bestruirui/bestsub/internal/core/task/exec/execer"
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
	var err error
	taskScheduler, err = scheduler.New()
	if err != nil {
		return err
	}
	return nil
}

// Start 启动任务调度器
func Start() {
	taskScheduler.Start()
}

func AddTask(req *task.CreateRequest) (*task.Data, error) {
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
	id, err := repo.Create(context.Background(), t)
	if err != nil {
		return nil, err
	}
	taskInfo := &exec.TaskInfo{
		ID:     id,
		Type:   req.Type,
		Config: []byte(req.Config),
	}
	err = taskScheduler.AddTask(req.Cron, execTask, taskInfo)
	return t, err
}

func UpdateTask(req *task.UpdateRequest) (*task.Data, error) {
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
		return nil, err
	}
	taskInfo := &exec.TaskInfo{
		ID:     req.ID,
		Type:   t.Type,
		Config: []byte(req.Config),
	}
	err = taskScheduler.UpdateTask(req.Cron, execTask, taskInfo)
	return t, err
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
func RemoveTaskWithDb(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	repo := database.Task()
	err := repo.Delete(context.Background(), taskID)
	if err != nil {
		return err
	}
	return taskScheduler.RemoveTask(taskID)
}

func RemoveTask(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	return taskScheduler.RemoveTask(taskID)
}

func StopTask(taskID int64) error {
	mu.Lock()
	defer mu.Unlock()
	return taskScheduler.StopTask(taskID)
}
func ListTasks(offset, pageSize int) (*[]task.Data, int64, error) {
	mu.RLock()
	defer mu.RUnlock()
	repo := database.Task()
	tasks, err := repo.List(context.Background(), offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return tasks, int64(len(*tasks)), nil
}

func GetTask(taskID int64) (*task.Data, error) {
	mu.RLock()
	defer mu.RUnlock()
	repo := database.Task()
	return repo.GetByID(context.Background(), taskID)
}

func Shutdown() error {
	mu.Lock()
	defer mu.Unlock()
	return taskScheduler.Stop()
}

func execTask(ctx context.Context, taskInfo exec.TaskInfo) {
	taskctx, cancel := context.WithCancel(ctx)
	defer cancel()
	taskScheduler.AddRunTaskMap(taskInfo.ID, cancel)
	defer taskScheduler.RemoveRunTaskMap(taskInfo.ID)

	for i := 0; i < taskInfo.Retry; i++ {
		select {
		case <-taskctx.Done():
			log.Infof("任务 %d 取消", taskInfo.ID)
			return
		default:
			log.Debugf("task %d running %d/%d", taskInfo.ID, i+1, taskInfo.Retry)
			err := exec.Run(taskctx, &taskInfo)
			if err != nil {
				log.Errorf("任务 %d 执行失败: %v (重试 %d/%d)", taskInfo.ID, err, i+1, taskInfo.Retry)
				continue
			}
			log.Infof("任务 %d 执行成功", taskInfo.ID)
			return
		}
	}
}
