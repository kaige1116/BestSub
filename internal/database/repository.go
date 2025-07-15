package database

import (
	"github.com/bestruirui/bestsub/internal/database/interfaces"
)

// Auth 获取认证仓库
func Auth() interfaces.AuthRepository {
	return GetRepository().Auth()
}

// SystemConfig 获取系统配置仓库
func SystemConfig() interfaces.SystemConfigRepository {
	return GetRepository().SystemConfig()
}

// Notify 获取通知渠道仓库
func Notify() interfaces.NotifyRepository {
	return GetRepository().Notify()
}

// NotifyTemplate 获取通知模板仓库
func NotifyTemplate() interfaces.NotifyTemplateRepository {
	return GetRepository().NotifyTemplate()
}

// Task 获取任务仓库
func Task() interfaces.TaskRepository {
	return GetRepository().Task()
}

// SubStorageConfig 获取存储配置仓库
func SubStorageConfig() interfaces.SubStorageConfigRepository {
	return GetRepository().SubStorageConfig()
}

// SubOutputTemplate 获取输出模板仓库
func SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	return GetRepository().SubOutputTemplate()
}

// SubNodeFilterRule 获取节点筛选规则仓库
func SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	return GetRepository().SubNodeFilterRule()
}

// SubLink 获取订阅链接仓库
func SubLink() interfaces.SubRepository {
	return GetRepository().SubLink()
}

// SubSaveConfig 获取保存配置仓库
func SubSaveConfig() interfaces.SubSaveRepository {
	return GetRepository().SubSaveConfig()
}

// SubShareLink 获取分享链接仓库
func SubShareLink() interfaces.SubShareRepository {
	return GetRepository().SubShareLink()
}
