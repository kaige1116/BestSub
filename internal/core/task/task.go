package task

import (
	"context"
	"time"

	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/exec"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/robfig/cron/v3"
)

var runningTasks = make(map[uint16]context.CancelFunc)
var cronTasks = make(map[uint16]cron.EntryID)
var cronInstance = cron.New(cron.WithLocation(time.Local))

func StartCron() {
	cronInstance.Start()
}

func StopCron() error {
	for _, cancel := range runningTasks {
		cancel()
	}
	cronInstance.Stop()
	runningTasks = make(map[uint16]context.CancelFunc)
	cronTasks = make(map[uint16]cron.EntryID)
	return nil
}

func Add(taskData *taskModel.Data) {
	var err error
	cronTasks[taskData.ID], err = cronInstance.AddFunc(taskData.Cron, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskData.Timeout)*time.Second)
		runningTasks[taskData.ID] = cancel
		logger, err := log.NewTaskLogger(taskData.ID, taskData.LogLevel)
		if err != nil {
			log.Errorf("failed to create logger: %v", err)
			return
		}
		defer logger.Close()
		execer, err := exec.Get(taskData.Type, taskData.Config)
		if err != nil {
			log.Errorf("failed to get execer: %v", err)
			return
		}
		execer.Run(ctx, logger)
		delete(runningTasks, taskData.ID)
	})
	if err != nil {
		log.Errorf("failed to add task: %v", err)
		return
	}
	log.Infof("task %d added", taskData.ID)
}
func NextRunTime(taskId uint16) time.Time {
	if _, ok := cronTasks[taskId]; ok {
		return cronInstance.Entry(cronTasks[taskId]).Next
	}
	return time.Time{}
}
func PrevRunTime(taskId uint16) time.Time {
	if _, ok := cronTasks[taskId]; ok {
		return cronInstance.Entry(cronTasks[taskId]).Prev
	}
	return time.Time{}
}
func Run(taskId uint16) {
	if _, ok := cronTasks[taskId]; ok {
		go cronInstance.Entry(cronTasks[taskId]).Job.Run()
	}
}
func Stop(taskId uint16) {
	if _, ok := runningTasks[taskId]; ok {
		runningTasks[taskId]()
	}
}
func Update(taskId uint16, taskData *taskModel.Data) {
	if _, ok := runningTasks[taskId]; ok {
		runningTasks[taskId]()
	}
	cronInstance.Remove(cronTasks[taskId])
	Add(taskData)
}
