package task

var taskConfigData = []Data{
	{
		Enable: true,
		Name:   "GC垃圾回收",
		System: true,
		Config: `{"cron":"0 */1 * * *","type":"gc","log_level":"info","timeout":10,"notify":false,"notify_channel":""}`,
		Extra:  `{"force":"false"}`,
	},
	{
		Enable: true,
		Name:   "日志清理",
		System: true,
		Config: `{"cron":"0 0 */1 * *","type":"log_clean","log_level":"info","timeout":10,"notify":false,"notify_channel":""}`,
		Extra:  "",
	},
}

func Default() []Data {
	return taskConfigData
}
