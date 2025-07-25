package task

import (
	"context"
	"errors"
	"time"

	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/exec"
	"github.com/bestruirui/bestsub/internal/utils/generic"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/robfig/cron/v3"
)

type Info struct {
	EntryID  cron.EntryID
	Name     string
	CronExpr string
}

type CronTask struct {
	Name     string
	taskFunc *generic.MapOf[uint16, func()]
	running  *generic.MapOf[uint16, context.CancelFunc]
	cron     *generic.MapOf[uint16, Info]
	instance *cron.Cron
	afterDo  func(taskId uint16, result taskModel.ReturnResult)
}

func New(name string, afterDo func(taskId uint16, result taskModel.ReturnResult)) *CronTask {
	return &CronTask{
		Name:     name,
		taskFunc: &generic.MapOf[uint16, func()]{},
		running:  &generic.MapOf[uint16, context.CancelFunc]{},
		cron:     &generic.MapOf[uint16, Info]{},
		instance: cron.New(cron.WithLocation(time.Local)),
		afterDo:  afterDo,
	}
}

func (ct *CronTask) Start() {
	ct.instance.Start()
}

func (ct *CronTask) Stop() error {
	ct.running.Range(func(key uint16, cancel context.CancelFunc) bool {
		cancel()
		return true
	})
	ct.instance.Stop()
	ct.running = &generic.MapOf[uint16, context.CancelFunc]{}
	ct.cron = &generic.MapOf[uint16, Info]{}
	ct.taskFunc = &generic.MapOf[uint16, func()]{}
	return nil
}

func (ct *CronTask) Add(taskConfig *taskModel.Config, extra string) {
	ct.taskFunc.Store(taskConfig.ID, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskConfig.Timeout)*time.Second)
		defer cancel()
		ct.running.Store(taskConfig.ID, cancel)
		defer ct.running.Delete(taskConfig.ID)
		logger, err := log.NewTaskLogger(ct.Name, taskConfig.ID, taskConfig.LogLevel, taskConfig.LogWriteFile)
		if err != nil {
			log.Errorf("failed to create logger: %v", err)
			return
		}
		go func() {
			<-ctx.Done()
			logger.Close()
		}()
		execer, err := exec.Get(taskConfig.Type, extra)
		if err != nil {
			log.Errorf("failed to get execer: %v", err)
			return
		}
		result := execer.Run(ctx, logger)
		ct.afterDo(taskConfig.ID, result)
	})
	ct.cron.Store(taskConfig.ID, Info{
		EntryID:  0,
		Name:     taskConfig.Name,
		CronExpr: taskConfig.CronExpr,
	})
	log.Infof("task %d %s added", taskConfig.ID, taskConfig.Name)
}

func (ct *CronTask) Enable(taskId uint16) error {
	if info, ok := ct.cron.Load(taskId); ok {
		if taskFunc, ok := ct.taskFunc.Load(taskId); ok {
			entryID, err := ct.instance.AddFunc(info.CronExpr, taskFunc)
			if err != nil {
				log.Errorf("failed to add task: %v", err)
				return err
			}
			info.EntryID = entryID
			ct.cron.Store(taskId, info)
			return nil
		}
	}
	return errors.New("task not found")
}

func (ct *CronTask) Disable(taskId uint16) error {
	if info, ok := ct.cron.Load(taskId); ok {
		if cancel, ok := ct.running.Load(taskId); ok {
			cancel()
		}

		ct.instance.Remove(info.EntryID)

		log.Infof("task %d %s disabled", taskId, info.Name)
		ct.running.Delete(taskId)
		return nil
	}
	return errors.New("task not found")
}

func (ct *CronTask) NextRunTime(taskId uint16) time.Time {
	if info, ok := ct.cron.Load(taskId); ok {
		if info.EntryID != 0 {
			return ct.instance.Entry(info.EntryID).Next
		}
	}
	return time.Time{}
}

func (ct *CronTask) Run(taskId uint16) error {
	if taskFunc, ok := ct.taskFunc.Load(taskId); ok {
		go taskFunc()
		return nil
	}
	return errors.New("task not found")
}

func (ct *CronTask) StopTask(taskId uint16) error {
	if _, ok := ct.cron.Load(taskId); ok {
		if cancel, ok := ct.running.Load(taskId); ok {
			cancel()
			return nil
		} else {
			return errors.New("task not running")
		}
	} else {
		return errors.New("task not found")
	}
}

func (ct *CronTask) Remove(taskId uint16) error {
	if info, ok := ct.cron.Load(taskId); ok {
		if cancel, ok := ct.running.Load(taskId); ok {
			cancel()
		}
		ct.instance.Remove(info.EntryID)
		log.Infof("task %d %s removed", taskId, info.Name)
		ct.cron.Delete(taskId)
		ct.running.Delete(taskId)
		ct.taskFunc.Delete(taskId)
		return nil
	}
	return errors.New("task not found")
}

func (ct *CronTask) Update(taskConfig *taskModel.Config, extra string) error {
	if info, ok := ct.cron.Load(taskConfig.ID); ok {
		if cancel, ok := ct.running.Load(taskConfig.ID); ok {
			cancel()
		}
		ct.instance.Remove(info.EntryID)
		ct.Add(taskConfig, extra)
		return nil
	}
	return errors.New("task not found")
}
func (ct *CronTask) Status(taskId uint16) string {
	if _, ok := ct.running.Load(taskId); ok {
		return "running"
	}
	if _, ok := ct.cron.Load(taskId); ok {
		return "pending"
	}
	return "stopped"
}
