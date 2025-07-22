package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration002Indexes 添加数据库索引
func Migration002Indexes() string {
	return `
-- 索引：用于通过 sub_id 查询对应的 task
CREATE INDEX IF NOT EXISTS idx_sub_task_sub_id ON sub_task_relations(sub_id);

-- 索引：用于通过 share_id 查询对应的 template
CREATE INDEX IF NOT EXISTS idx_share_template_share_id ON share_template_relations(share_id);

-- 索引：用于通过 share_id 查询对应的 filter
CREATE INDEX IF NOT EXISTS idx_share_filter_share_id ON share_filter_relations(share_id);

-- 索引：用于通过 share_id 查询对应的 sub
CREATE INDEX IF NOT EXISTS idx_share_sub_share_id ON share_sub_relations(share_id);

-- 索引：用于通过 save_id 查询对应的 storage
CREATE INDEX IF NOT EXISTS idx_save_storage_save_id ON save_storage_relations(save_id);

-- 索引：用于通过 save_id 查询对应的 template
CREATE INDEX IF NOT EXISTS idx_save_template_save_id ON save_template_relations(save_id);

-- 索引：用于通过 save_id 查询对应的 filter
CREATE INDEX IF NOT EXISTS idx_save_filter_save_id ON save_filter_relations(save_id);

-- 索引：用于通过 save_id 查询对应的 sub
CREATE INDEX IF NOT EXISTS idx_save_sub_save_id ON save_sub_relations(save_id);

-- 索引：用于通过 task_id 查询对应的 save
CREATE INDEX IF NOT EXISTS idx_save_task_task_id ON save_task_relations(task_id);

-- 索引：用于通过 task_id 查询对应的 notify
CREATE INDEX IF NOT EXISTS idx_notify_task_task_id ON notify_task_relations(task_id);

`
}

// init 自动注册迁移
func init() {
	migration.Register(ClientName, 202507171101, "dev", "Indexes", Migration002Indexes)
}
