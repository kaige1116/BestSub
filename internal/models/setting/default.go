package setting

var system = GroupSettingAdvance{
	GroupName: "系统配置",
	Data: []SettingAdvance{
		{
			Name:  "代理",
			Type:  "boolean",
			Key:   PROXY_ENABLE,
			Value: "false",
			Desc:  "是否启用代理",
		},
		{
			Name:  "代理地址",
			Type:  "string",
			Key:   PROXY_URL,
			Value: "socks5://user:pass@127.0.0.1:1080",
		},
		{
			Name:  "日志保留天数",
			Type:  "number",
			Key:   LOG_RETENTION_DAYS,
			Value: "7",
		},
	},
}
var node = GroupSettingAdvance{
	GroupName: "节点配置",
	Data: []SettingAdvance{
		{
			Name:  "节点池大小",
			Type:  "number",
			Key:   NODE_POOL_SIZE,
			Value: "1000",
		},
		{
			Name:  "默认测试地址",
			Type:  "string",
			Key:   NODE_TEST_URL,
			Value: "https://www.gstatic.com/generate_204",
		},
		{
			Name:  "默认测试超时时间（秒）",
			Type:  "number",
			Key:   NODE_TEST_TIMEOUT,
			Value: "5",
		},
	},
}
var task = GroupSettingAdvance{
	GroupName: "任务配置",
	Data: []SettingAdvance{
		{
			Name:  "最大线程数",
			Type:  "number",
			Key:   TASK_MAX_THREAD,
			Value: "200",
		},
		{
			Name:  "任务最大超时时间（秒）",
			Type:  "number",
			Key:   TASK_MAX_TIMEOUT,
			Value: "60",
		},
		{
			Name:  "任务最大重试次数",
			Type:  "number",
			Key:   TASK_MAX_RETRY,
			Value: "3",
		},
	},
}
var notify = GroupSettingAdvance{
	GroupName: "通知配置",
	Data: []SettingAdvance{
		{
			Name:  "需要通知的操作类型",
			Type:  "number",
			Key:   NOTIFY_OPERATION,
			Value: "0",
		},
		{
			Name:  "系统默认通知渠道",
			Type:  "number",
			Key:   NOTIFY_ID,
			Value: "0",
		},
	},
}

var defaultSetting = []GroupSettingAdvance{
	system,
	node,
	task,
	notify,
}

func DefaultSetting() []GroupSettingAdvance {
	return defaultSetting
}
