package task

import (
	"context"
	"encoding/json"

	"github.com/bestruirui/bestsub/internal/database/op"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var Fetch = New("fetch", func(taskId uint16, result taskModel.ReturnResult) {
	op.UpdateSubResult(context.Background(), taskId, result)
})
var Check = New("check", func(taskId uint16, result taskModel.ReturnResult) {
	op.UpdateCheckResult(context.Background(), taskId, result)
})

func Load() {
	subList, err := op.GetSubList(context.Background())
	if err != nil {
		log.Errorf("failed to get sub list: %v", err)
		return
	}
	for _, sub := range subList {
		var taskConfig taskModel.Config
		if err := json.Unmarshal([]byte(sub.Config), &taskConfig); err != nil {
			log.Errorf("failed to unmarshal sub task: %v", err)
			continue
		}
		taskConfig.ID = sub.ID
		taskConfig.Name = sub.Name
		taskConfig.Type = "fetch"
		taskConfig.Timeout = 60
		Fetch.Add(&taskConfig, sub.Config)
		if sub.Enable {
			Fetch.Enable(sub.ID)
		}
	}

	checkList, err := op.GetCheckList(context.Background())
	if err != nil {
		log.Errorf("failed to get check list: %v", err)
		return
	}
	for _, check := range checkList {
		var taskConfig taskModel.Config
		if err := json.Unmarshal([]byte(check.Task), &taskConfig); err != nil {
			log.Errorf("failed to unmarshal check task: %v", err)
			continue
		}
		taskConfig.ID = check.ID
		taskConfig.Name = check.Name
		Check.Add(&taskConfig, check.Config)
		if check.Enable {
			Check.Enable(check.ID)
		}
	}
}
