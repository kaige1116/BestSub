package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

// Migration002AddIndexes 添加数据库索引
func Migration002AddIndexes() string {
	return `

`
}

// init 自动注册所有迁移
func init() {
	migration.Register(migrations, "002", "Add database indexes", Migration002AddIndexes)
}
