package interfaces

import (
	"context"
)

// Migration 代表一个应用程序升级迁移
type Migration struct {
	Version     string // 迁移版本号
	Description string // 迁移描述
}

// Migrator 数据库迁移器接口
type Migrator interface {
	// Apply 验证并应用所有待执行的迁移到最新版本
	Apply(ctx context.Context) error
}
