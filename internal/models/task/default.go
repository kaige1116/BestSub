package task

var taskConfigData = []Data{
	{
		Name:    "GC垃圾回收",
		Enable:  true,
		Cron:    "0 */1 * * *",
		Timeout: 10,
		System:  true,
		Type:    TypeGC,
	},
	{
		Name:    "日志清理",
		Enable:  true,
		Cron:    "0 0 */1 * *",
		Timeout: 10,
		System:  true,
		Type:    TypeLogClean,
	},
}

func Default() []Data {
	return taskConfigData
}
