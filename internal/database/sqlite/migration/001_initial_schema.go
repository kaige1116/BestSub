package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration001InitialSchema 初始数据库架构
func Migration001InitialSchema() string {
	return `
-- 认证表
CREATE TABLE IF NOT EXISTS auth (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 会话表
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    refresh_token TEXT UNIQUE,
    ip_address TEXT,
    user_agent TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 通知渠道表
CREATE TABLE IF NOT EXISTS notification_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL, -- email, webhook, telegram, etc.
    config TEXT NOT NULL, -- JSON配置
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- subscription_update, cleanup, etc.
    status TEXT NOT NULL DEFAULT 'pending', -- pending, running, completed, failed
    config TEXT, -- JSON配置
    result TEXT, -- 执行结果
    error_message TEXT,
    scheduled_at DATETIME,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 订阅存储配置表
CREATE TABLE IF NOT EXISTS sub_storage_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL, -- local, s3, ftp, etc.
    config TEXT NOT NULL, -- JSON配置
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 订阅输出模板表
CREATE TABLE IF NOT EXISTS sub_output_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    format TEXT NOT NULL, -- clash, v2ray, surge, etc.
    template TEXT NOT NULL, -- 模板内容
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 订阅节点筛选规则表
CREATE TABLE IF NOT EXISTS sub_node_filter_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL, -- include, exclude
    pattern TEXT NOT NULL, -- 正则表达式或关键词
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 订阅链接表
CREATE TABLE IF NOT EXISTS sub_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    update_interval INTEGER NOT NULL DEFAULT 3600, -- 更新间隔（秒）
    last_updated DATETIME,
    last_error TEXT,
    node_count INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 订阅链接模块配置表
CREATE TABLE IF NOT EXISTS sub_link_module_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sub_link_id INTEGER NOT NULL,
    module_type TEXT NOT NULL, -- filter, template, storage
    module_id INTEGER NOT NULL,
    config TEXT, -- 额外配置
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sub_link_id) REFERENCES sub_links(id) ON DELETE CASCADE
);

-- 订阅保存配置表
CREATE TABLE IF NOT EXISTS sub_save_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    sub_link_ids TEXT NOT NULL, -- JSON数组，包含订阅链接ID
    output_template_id INTEGER,
    storage_config_id INTEGER,
    filter_rule_ids TEXT, -- JSON数组，包含筛选规则ID
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (output_template_id) REFERENCES sub_output_templates(id),
    FOREIGN KEY (storage_config_id) REFERENCES sub_storage_configs(id)
);

-- 订阅分享链接表
CREATE TABLE IF NOT EXISTS sub_share_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    sub_save_config_id INTEGER NOT NULL,
    access_count INTEGER NOT NULL DEFAULT 0,
    max_access_count INTEGER, -- NULL表示无限制
    expires_at DATETIME, -- NULL表示永不过期
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sub_save_config_id) REFERENCES sub_save_configs(id) ON DELETE CASCADE
);
`
}

// Migration001InitialSchemaDown 回滚初始数据库架构
func Migration001InitialSchemaDown() string {
	return `
-- 删除表（按依赖关系逆序）
DROP TABLE IF EXISTS sub_share_links;
DROP TABLE IF EXISTS sub_save_configs;
DROP TABLE IF EXISTS sub_link_module_configs;
DROP TABLE IF EXISTS sub_links;
DROP TABLE IF EXISTS sub_node_filter_rules;
DROP TABLE IF EXISTS sub_output_templates;
DROP TABLE IF EXISTS sub_storage_configs;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS notification_channels;
DROP TABLE IF EXISTS system_config;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS auth;
`
}

// init 自动注册所有迁移
func init() {
	migration.Register(migrations, "001", "Initial database schema", Migration001InitialSchema)
}
