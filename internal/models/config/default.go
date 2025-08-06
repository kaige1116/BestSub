package config

var defaultBase = Base{
	Server: ServerConfig{
		Port: 8080,
		Host: "0.0.0.0",
	},
	Database: DatabaseConfig{
		Type: "sqlite",
	},
	Log: LogConfig{
		Level:  "debug",
		Output: "console",
	},
}

var system = GroupAdvance{
	GroupName:   "system",
	Description: "系统配置",
	Data: []Advance{
		{
			Name:    "代理",
			Type:    "bool",
			Key:     "proxy.enable",
			Default: "false",
			Desc:    "是否启用代理",
		},
		{
			Name:    "代理地址",
			Type:    "string",
			Key:     "proxy.url",
			Default: "socks5://user:pass@127.0.0.1:1080",
		},
		{
			Name:    "节点池大小",
			Type:    "number",
			Key:     "node.pool_size",
			Default: "1000",
		},
		{
			Name:    "任务最大超时时间（秒）",
			Type:    "number",
			Key:     "task.max_timeout",
			Default: "60",
		},
		{
			Name:    "任务最大重试次数",
			Type:    "number",
			Key:     "task.max_retry",
			Default: "3",
		},
		{
			Name:    "日志保留天数",
			Type:    "number",
			Key:     "log.retention_days",
			Default: "7",
		},
		{
			Name:    "需要通知的操作类型",
			Type:    "number",
			Key:     "notify.operation",
			Default: "0",
		},
		{
			Name:    "系统默认通知渠道",
			Type:    "number",
			Key:     "notify.id",
			Default: "0",
		},
	},
}

var defaultAdvance = []GroupAdvance{
	system,
}

func DefaultAdvance() []GroupAdvance {
	return defaultAdvance
}

func DefaultBase() Base {
	return defaultBase
}
