package task

import (
	"context"
	"errors"
	"time"

	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/exec"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/robfig/cron/v3"
)

var Fetch = New()
var Check = New()

type Info struct {
	EntryID cron.EntryID
	Name    string
}

type CronTask struct {
	running  map[uint16]context.CancelFunc
	cron     map[uint16]Info
	instance *cron.Cron
}

func New() *CronTask {
	return &CronTask{
		running:  make(map[uint16]context.CancelFunc),
		cron:     make(map[uint16]Info),
		instance: cron.New(cron.WithLocation(time.Local)),
	}
}

func (ct *CronTask) Start() {
	ct.instance.Start()
}

func (ct *CronTask) Stop() error {
	for _, cancel := range ct.running {
		cancel()
	}
	ct.instance.Stop()
	ct.running = make(map[uint16]context.CancelFunc)
	ct.cron = make(map[uint16]Info)
	return nil
}

func (ct *CronTask) Add(taskConfig *taskModel.Config) {
	var err error
	entryID, err := ct.instance.AddFunc(taskConfig.CronExpr, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskConfig.Timeout)*time.Second)
		ct.running[taskConfig.ID] = cancel
		logger, err := log.NewTaskLogger(taskConfig.ID, taskConfig.LogLevel, taskConfig.LogWriteFile)
		if err != nil {
			log.Errorf("failed to create logger: %v", err)
			return
		}
		defer logger.Close()
		execer, err := exec.Get(taskConfig.Type, taskConfig.Extra)
		if err != nil {
			log.Errorf("failed to get execer: %v", err)
			return
		}
		execer.Run(ctx, logger)
		delete(ct.running, taskConfig.ID)
	})
	if err != nil {
		log.Errorf("failed to add task: %v", err)
		return
	}
	ct.cron[taskConfig.ID] = Info{
		EntryID: entryID,
		Name:    taskConfig.Name,
	}
	log.Infof("task %d %s added next run time: %s", taskConfig.ID, taskConfig.Name, ct.instance.Entry(ct.cron[taskConfig.ID].EntryID).Next)
}

func (ct *CronTask) NextRunTime(taskId uint16) time.Time {
	if _, ok := ct.cron[taskId]; ok {
		return ct.instance.Entry(ct.cron[taskId].EntryID).Next
	}
	return time.Time{}
}

func (ct *CronTask) Run(taskId uint16) error {
	if _, ok := ct.cron[taskId]; ok {
		go ct.instance.Entry(ct.cron[taskId].EntryID).Job.Run()
		return nil
	}
	return errors.New("task not found")
}

func (ct *CronTask) StopTask(taskId uint16) error {
	if _, ok := ct.cron[taskId]; ok {
		if _, ok := ct.running[taskId]; ok {
			ct.running[taskId]()
			return nil
		} else {
			return errors.New("task not running")
		}
	} else {
		return errors.New("task not found")
	}
}

func (ct *CronTask) Remove(taskId uint16) error {
	if _, ok := ct.cron[taskId]; ok {
		ct.instance.Remove(ct.cron[taskId].EntryID)
		delete(ct.cron, taskId)
		delete(ct.running, taskId)
		log.Infof("task %d %s removed", taskId, ct.cron[taskId].Name)
		return nil
	}
	return errors.New("task not found")
}

func (ct *CronTask) Update(taskConfig *taskModel.Config) {
	if _, ok := ct.running[taskConfig.ID]; ok {
		ct.running[taskConfig.ID]()
	}
	ct.instance.Remove(ct.cron[taskConfig.ID].EntryID)
	ct.Add(taskConfig)
}

func (ct *CronTask) GetRunningTaskID() []uint16 {
	var taskList []uint16
	for taskId := range ct.running {
		taskList = append(taskList, taskId)
	}
	return taskList
}
