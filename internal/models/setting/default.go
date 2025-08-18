package setting

var system = GroupSettingAdvance{
	GroupName:   "system",
	Description: "系统配置",
	Data: []SettingAdvance{
		{
			Name:  "代理",
			Type:  "bool",
			Key:   "proxy.enable",
			Value: "false",
			Desc:  "是否启用代理",
		},
		{
			Name:  "代理地址",
			Type:  "string",
			Key:   "proxy.url",
			Value: "socks5://user:pass@127.0.0.1:1080",
		},
		{
			Name:  "日志保留天数",
			Type:  "number",
			Key:   "log.retention_days",
			Value: "7",
		},
	},
}
var node = GroupSettingAdvance{
	GroupName:   "node",
	Description: "节点配置",
	Data: []SettingAdvance{
		{
			Name:  "节点池大小",
			Type:  "number",
			Key:   "node.pool_size",
			Value: "1000",
		},
		{
			Name:  "默认测试地址",
			Type:  "string",
			Key:   "node.test_url",
			Value: "https://www.gstatic.com/generate_204",
		},
		{
			Name:  "默认测试超时时间（秒）",
			Type:  "number",
			Key:   "node.test_timeout",
			Value: "5",
		},
	},
}
var task = GroupSettingAdvance{
	GroupName:   "task",
	Description: "任务配置",
	Data: []SettingAdvance{
		{
			Name:  "最大线程数",
			Type:  "number",
			Key:   "task.max_thread",
			Value: "200",
		},
		{
			Name:  "任务最大超时时间（秒）",
			Type:  "number",
			Key:   "task.max_timeout",
			Value: "60",
		},
		{
			Name:  "任务最大重试次数",
			Type:  "number",
			Key:   "task.max_retry",
			Value: "3",
		},
	},
}
var notify = GroupSettingAdvance{
	GroupName:   "notify",
	Description: "通知配置",
	Data: []SettingAdvance{
		{
			Name:  "需要通知的操作类型",
			Type:  "number",
			Key:   "notify.operation",
			Value: "0",
		},
		{
			Name:  "系统默认通知渠道",
			Type:  "number",
			Key:   "notify.id",
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
