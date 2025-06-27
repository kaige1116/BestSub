package defaults

import (
	"github.com/bestruirui/bestsub/internal/database/models"
)

// Configs 获取默认的系统配置
func Configs() []models.SystemConfig {
	return []models.SystemConfig{
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
	}
}
