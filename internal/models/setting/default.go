package setting

import "github.com/bestruirui/bestsub/internal/utils/desc"

func DefaultSetting() []GroupSettingAdvance {
	return []GroupSettingAdvance{
		{
			GroupName: "系统配置",
			Data: []SettingAdvance{
				{
					Name:  "代理",
					Type:  desc.TypeBoolean,
					Key:   PROXY_ENABLE,
					Value: "false",
					Desc:  "是否启用代理",
				},
				{
					Name:  "代理地址",
					Type:  desc.TypeString,
					Key:   PROXY_URL,
					Value: "socks5://user:pass@127.0.0.1:1080",
				},
				{
					Name:  "日志保留天数",
					Type:  desc.TypeNumber,
					Key:   LOG_RETENTION_DAYS,
					Value: "7",
				},
				{
					Name:  "前端地址",
					Type:  desc.TypeString,
					Key:   FRONTEND_URL,
					Value: "https://github.com/BestSubOrg/Front/releases/latest/download/out.zip",
				},
				{
					Name:  "前端代理",
					Desc:  "是否启用代理更新前端UI",
					Type:  desc.TypeBoolean,
					Key:   FRONTEND_URL_PROXY,
					Value: "false",
				},
				{
					Name:  "subconverter地址",
					Type:  desc.TypeString,
					Key:   SUBCONVERTER_URL,
					Value: "https://github.com/BestSubOrg/subconverter/releases/latest/download/",
				},
				{
					Name:  "subconverter代理",
					Desc:  "是否启用代理更新subconverter",
					Type:  desc.TypeBoolean,
					Key:   SUBCONVERTER_URL_PROXY,
					Value: "false",
				},
				{
					Name:  "自动禁用订阅",
					Type:  desc.TypeNumber,
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
					Type:  desc.TypeNumber,
					Key:   NODE_POOL_SIZE,
					Value: "1000",
				},
				{
					Name:  "默认测试地址",
					Type:  desc.TypeString,
					Key:   NODE_TEST_URL,
					Value: "https://www.gstatic.com/generate_204",
				},
				{
					Name:  "默认测试超时时间（秒）",
					Type:  desc.TypeNumber,
					Key:   NODE_TEST_TIMEOUT,
					Value: "5",
				},
				{
					Name:  "全局协议过滤启用",
					Type:  desc.TypeBoolean,
					Key:   NODE_PROTOCOL_FILTER_ENABLE,
					Value: "false",
					Desc:  "是否启用全局协议过滤",
				},
				{
					Name:  "全局协议过滤模式",
					Type:  desc.TypeBoolean,
					Key:   NODE_PROTOCOL_FILTER_MODE,
					Value: "false",
					Desc:  "关闭为排除,打开为包含",
				},
				{
					Name:    "全局协议过滤",
					Type:    desc.TypeMultiSelect,
					Key:     NODE_PROTOCOL_FILTER,
					Options: "vless,vmess,trojan,shadowsocks,socks5,http,https",
				},
			},
		},
		{
			GroupName: "任务配置",
			Data: []SettingAdvance{
				{
					Name:  "最大线程数",
					Type:  desc.TypeNumber,
					Key:   TASK_MAX_THREAD,
					Value: "200",
				},
				{
					Name:  "任务最大超时时间（秒）",
					Type:  desc.TypeNumber,
					Key:   TASK_MAX_TIMEOUT,
					Value: "60",
				},
				{
					Name:  "任务最大重试次数",
					Type:  desc.TypeNumber,
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
					Type:  desc.TypeNumber,
					Key:   NOTIFY_OPERATION,
					Value: "0",
				},
				{
					Name:  "系统默认通知渠道",
					Type:  desc.TypeNumber,
					Key:   NOTIFY_ID,
					Value: "0",
				},
			},
		},
	}
}
