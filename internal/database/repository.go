package database

import (
	"github.com/bestruirui/bestsub/internal/database/interfaces"
)

// =============================================================================
// 快捷访问方法 - 直接访问仓库实例的便捷函数
// =============================================================================

// Auth 获取认证仓库
func Auth() interfaces.AuthRepository {
	return GetRepository().Auth()
}

// Session 获取会话仓库
func Session() interfaces.SessionRepository {
	return GetRepository().Session()
}

// SystemConfig 获取系统配置仓库
func SystemConfig() interfaces.SystemConfigRepository {
	return GetRepository().SystemConfig()
}

// Notify 获取通知渠道仓库
func Notify() interfaces.NotifyRepository {
	return GetRepository().Notify()
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
func SubLink() interfaces.SubLinkRepository {
	return GetRepository().SubLink()
}

// SubLinkModuleConfig 获取链接模块配置仓库
func SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository {
	return GetRepository().SubLinkModuleConfig()
}

// SubSaveConfig 获取保存配置仓库
func SubSaveConfig() interfaces.SubSaveConfigRepository {
	return GetRepository().SubSaveConfig()
}

// SubShareLink 获取分享链接仓库
func SubShareLink() interfaces.SubShareLinkRepository {
	return GetRepository().SubShareLink()
}
