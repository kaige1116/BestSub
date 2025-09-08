package setting

func DefaultSetting() []GroupSettingAdvance {
	return []GroupSettingAdvance{
		{
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
				{
					Name:  "前端地址",
					Type:  "string",
					Key:   FRONTEND_URL,
					Value: "https://github.com/BestSubOrg/Front/releases/latest/download/out.zip",
				},
				{
					Name:  "前端代理",
					Desc:  "是否启用代理更新前端UI",
					Type:  "boolean",
					Key:   FRONTEND_URL_PROXY,
					Value: "false",
				},
				{
					Name:  "subconverter地址",
					Type:  "string",
					Key:   SUBCONVERTER_URL,
					Value: "https://github.com/BestSubOrg/subconverter/releases/latest/download/",
				},
				{
					Name:  "subconverter代理",
					Desc:  "是否启用代理更新subconverter",
					Type:  "boolean",
					Key:   SUBCONVERTER_URL_PROXY,
					Value: "false",
				},
				{
					Name:  "自动禁用订阅",
					Type:  "number",
					Key:   SUB_DISABLE_AUTO,
					Value: "0",
					Desc:  "当订阅获取节点数量为0的次数大于该值时,自动禁用订阅,0为不自动禁用",
				},
			},
		},
		{
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
		},
		{
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
		},
		{
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
		},
	}
}
