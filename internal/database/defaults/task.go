package defaults

import (
	"github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/task"
)

func Tasks() []task.Data {
	return []task.Data{
		{
			BaseDbModel: common.BaseDbModel{
				Name:        "GC垃圾回收",
				Description: "回收内存",
				Enable:      true,
			},
			Cron:      "0 */1 * * *",
			IsSysTask: true,
			Type:      task.TypeGC,
		},
		{
			BaseDbModel: common.BaseDbModel{
				Name:        "会话清理",
				Description: "清理会话",
				Enable:      true,
			},
			Cron:      "0 */1 * * *",
			IsSysTask: true,
			Type:      task.TypeSessionClean,
		},
	}
}
