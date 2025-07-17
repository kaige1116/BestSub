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
			Value:       "true",
			Description: "是否启用代理",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.type",
			Value:       "",
			Description: "代理类型",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.host",
			Value:       "",
			Description: "代理地址",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.port",
			Value:       "",
			Description: "代理端口",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.username",
			Value:       "",
			Description: "代理用户名",
		},
		{
			GroupName:   "proxy",
			Key:         "proxy.password",
			Value:       "",
			Description: "代理密码",
		},
		{
			GroupName:   "task",
			Key:         "task.max_timeout",
			Value:       "600",
			Description: "任务最大超时时间（秒）",
		},
		{
			GroupName:   "task",
			Key:         "task.max_retry",
			Value:       "10",
			Description: "任务最大重试次数",
		},
		{
			GroupName:   "log",
			Key:         "log.max_days",
			Value:       "7",
			Description: "日志保留天数",
		},
	}
}
