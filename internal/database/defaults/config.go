package defaults

import (
	"github.com/bestruirui/bestsub/internal/models/system"
)

// Configs 获取默认的系统配置
func Configs() []system.Data {
	return []system.Data{
		{
			GroupName:   "proxy",
			Key:         "proxy.enable",
			Type:        "bool",
			Value:       "true",
			Description: "是否启用代理",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.type",
			Type:        "string",
			Value:       "",
			Description: "代理类型",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.host",
			Type:        "string",
			Value:       "",
			Description: "代理地址",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.port",
			Type:        "int",
			Value:       "",
			Description: "代理端口",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.username",
			Type:        "string",
			Value:       "",
			Description: "代理用户名",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.password",
			Type:        "string",
			Value:       "",
			Description: "代理密码",
		},
		{
			GroupName:   "task",
			Key:         "task.max_timeout",
			Type:        "int",
			Value:       "600",
			Description: "任务最大超时时间（秒）",
		},
		{
			GroupName:   "task",
			Key:         "task.max_retry",
			Type:        "int",
			Value:       "10",
			Description: "任务最大重试次数",
		},
		{
			GroupName:   "log",
			Key:         "log.max_days",
			Type:        "int",
			Value:       "7",
			Description: "日志保留天数",
		},
	}
}
