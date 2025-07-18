package task

import "github.com/bestruirui/bestsub/internal/models/common"

var taskConfigData = []Data{
	{
		BaseDbModel: common.BaseDbModel{
			Name:        "GC垃圾回收",
			Description: "回收内存",
			Enable:      true,
		},
		Cron:      "0 */1 * * *",
		Timeout:   10,
		Retry:     3,
		IsSysTask: true,
		Type:      TypeGC,
	},
	{
		BaseDbModel: common.BaseDbModel{
			Name:        "日志清理",
			Description: "清理日志",
			Enable:      true,
		},
		Cron:      "0 0 */1 * *",
		Timeout:   10,
		Retry:     3,
		IsSysTask: true,
		Type:      TypeLogClean,
	},
}

func Default() []Data {
	return taskConfigData
}
