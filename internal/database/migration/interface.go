package migration

import (
	"context"
	"database/sql"
	"time"
)

// Migration 迁移接口
type Migration interface {
	// ID 获取迁移ID
	ID() string

	// Description 获取迁移描述
	Description() string

	// Up 执行向上迁移
	Up(ctx context.Context, tx *sql.Tx) error

	// Down 执行向下迁移（回滚）
	Down(ctx context.Context, tx *sql.Tx) error
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID          string    `db:"id" json:"id"`
	Description string    `db:"description" json:"description"`
	AppliedAt   time.Time `db:"applied_at" json:"applied_at"`
	Checksum    string    `db:"checksum" json:"checksum"`
}

// Migrator 迁移器接口
type Migrator interface {
	// Initialize 初始化迁移系统
	Initialize(ctx context.Context) error

	// GetAppliedMigrations 获取已应用的迁移
	GetAppliedMigrations(ctx context.Context) ([]MigrationRecord, error)

	// ApplyMigration 应用单个迁移
	ApplyMigration(ctx context.Context, migration Migration) error

	// RollbackMigration 回滚单个迁移
	RollbackMigration(ctx context.Context, migration Migration) error

	// RecordMigration 记录迁移
	RecordMigration(ctx context.Context, migration Migration) error

	// RemoveMigrationRecord 移除迁移记录
	RemoveMigrationRecord(ctx context.Context, migrationID string) error
}

// Manager 迁移管理器接口
type Manager interface {
	// RegisterMigrations 注册迁移
	RegisterMigrations(migrations ...Migration)

	// GetMigrations 获取所有注册的迁移
	GetMigrations() []Migration

	// GetMigration 根据ID获取迁移
	GetMigration(id string) (Migration, bool)

	// GetPendingMigrations 获取待执行的迁移
	GetPendingMigrations(ctx context.Context) ([]Migration, error)

	// MigrateUp 执行向上迁移
	MigrateUp(ctx context.Context) error

	// MigrateUpTo 迁移到指定版本
	MigrateUpTo(ctx context.Context, targetID string) error

	// MigrateDown 执行向下迁移
	MigrateDown(ctx context.Context, steps int) error

	// MigrateDownTo 回滚到指定版本
	MigrateDownTo(ctx context.Context, targetID string) error

	// GetStatus 获取迁移状态
	GetStatus(ctx context.Context) (*MigrationStatus, error)

	// Validate 验证迁移完整性
	Validate(ctx context.Context) error
}

// MigrationStatus 迁移状态
type MigrationStatus struct {
	TotalMigrations   int               `json:"total_migrations"`
	AppliedMigrations int               `json:"applied_migrations"`
	PendingMigrations int               `json:"pending_migrations"`
	LastMigration     *MigrationRecord  `json:"last_migration,omitempty"`
	AppliedList       []MigrationRecord `json:"applied_list"`
	PendingList       []Migration       `json:"pending_list"`
}

// MigrationFactory 迁移工厂接口
type MigrationFactory interface {
	// CreateMigrator 创建迁移器
	CreateMigrator() Migrator

	// CreateManager 创建迁移管理器
	CreateManager(migrator Migrator) Manager

	// GetAllMigrations 获取所有迁移
	GetAllMigrations() []Migration
}

// DatabaseMigrationSupport 数据库迁移支持接口
type DatabaseMigrationSupport interface {
	// GetMigrationFactory 获取迁移工厂
	GetMigrationFactory() MigrationFactory

	// RunMigrations 运行迁移
	RunMigrations(ctx context.Context) error

	// GetMigrationStatus 获取迁移状态
	GetMigrationStatus(ctx context.Context) (*MigrationStatus, error)
}
