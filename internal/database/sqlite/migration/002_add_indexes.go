package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration002AddIndexes 添加数据库索引
func Migration002AddIndexes() string {
	return `
-- 会话表索引
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);

-- 系统配置表索引
CREATE INDEX IF NOT EXISTS idx_system_config_key ON system_config(key);

-- 通知渠道表索引
CREATE INDEX IF NOT EXISTS idx_notification_channels_type ON notification_channels(type);
CREATE INDEX IF NOT EXISTS idx_notification_channels_enabled ON notification_channels(enabled);

-- 任务表索引
CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks(type);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_scheduled_at ON tasks(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- 订阅存储配置表索引
CREATE INDEX IF NOT EXISTS idx_sub_storage_configs_type ON sub_storage_configs(type);
CREATE INDEX IF NOT EXISTS idx_sub_storage_configs_enabled ON sub_storage_configs(enabled);

-- 订阅输出模板表索引
CREATE INDEX IF NOT EXISTS idx_sub_output_templates_format ON sub_output_templates(format);
CREATE INDEX IF NOT EXISTS idx_sub_output_templates_enabled ON sub_output_templates(enabled);

-- 订阅节点筛选规则表索引
CREATE INDEX IF NOT EXISTS idx_sub_node_filter_rules_type ON sub_node_filter_rules(type);
CREATE INDEX IF NOT EXISTS idx_sub_node_filter_rules_enabled ON sub_node_filter_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_sub_node_filter_rules_priority ON sub_node_filter_rules(priority);

-- 订阅链接表索引
CREATE INDEX IF NOT EXISTS idx_sub_links_enabled ON sub_links(enabled);
CREATE INDEX IF NOT EXISTS idx_sub_links_last_updated ON sub_links(last_updated);
CREATE INDEX IF NOT EXISTS idx_sub_links_created_at ON sub_links(created_at);

-- 订阅链接模块配置表索引
CREATE INDEX IF NOT EXISTS idx_sub_link_module_configs_sub_link_id ON sub_link_module_configs(sub_link_id);
CREATE INDEX IF NOT EXISTS idx_sub_link_module_configs_module_type ON sub_link_module_configs(module_type);
CREATE INDEX IF NOT EXISTS idx_sub_link_module_configs_enabled ON sub_link_module_configs(enabled);
CREATE INDEX IF NOT EXISTS idx_sub_link_module_configs_priority ON sub_link_module_configs(priority);

-- 订阅保存配置表索引
CREATE INDEX IF NOT EXISTS idx_sub_save_configs_enabled ON sub_save_configs(enabled);
CREATE INDEX IF NOT EXISTS idx_sub_save_configs_output_template_id ON sub_save_configs(output_template_id);
CREATE INDEX IF NOT EXISTS idx_sub_save_configs_storage_config_id ON sub_save_configs(storage_config_id);

-- 订阅分享链接表索引
CREATE INDEX IF NOT EXISTS idx_sub_share_links_token ON sub_share_links(token);
CREATE INDEX IF NOT EXISTS idx_sub_share_links_sub_save_config_id ON sub_share_links(sub_save_config_id);
CREATE INDEX IF NOT EXISTS idx_sub_share_links_enabled ON sub_share_links(enabled);
CREATE INDEX IF NOT EXISTS idx_sub_share_links_expires_at ON sub_share_links(expires_at);
CREATE INDEX IF NOT EXISTS idx_sub_share_links_created_at ON sub_share_links(created_at);
`
}

// init 自动注册所有迁移
func init() {
	migration.Register(migrations, "002", "Add database indexes", Migration002AddIndexes)
}
