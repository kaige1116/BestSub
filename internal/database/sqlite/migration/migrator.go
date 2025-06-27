package migration

import (
	"context"
	"fmt"
	"sort"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/migration"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// migrations 迁移注册表
var migrations = make(map[string]migration.MigrationInfo)

// Migrator SQLite迁移器实现
type Migrator struct {
	db *database.Database
}

// NewMigrator 创建新的SQLite迁移器
func NewMigrator(db *database.Database) interfaces.Migrator {
	return &Migrator{db: db}
}

// Apply 验证并应用所有待执行的迁移到最新版本
func (m *Migrator) Apply(ctx context.Context) error {
	// 首先验证迁移脚本
	if err := migration.Validate(migrations); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	// 确保迁移表存在
	if err := m.createMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// 获取所有迁移
	migrations := migration.GetMigrations(migrations)
	if len(migrations) == 0 {
		return nil // 没有迁移需要执行
	}

	// 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return migration.CompareVersions(migrations[i].Version, migrations[j].Version) < 0
	})

	// 执行待执行的迁移
	for _, migration := range migrations {
		applied, err := m.isApplied(ctx, migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check if migration %s is applied: %w", migration.Version, err)
		}

		if !applied {
			if err := m.applyMigration(ctx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
		}
	}

	return nil
}

// createMigrationTable 创建迁移记录表
func (m *Migrator) createMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL,
			success BOOLEAN NOT NULL DEFAULT TRUE,
			error TEXT
		)
	`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

// isApplied 检查迁移是否已应用
func (m *Migrator) isApplied(ctx context.Context, version string) (bool, error) {
	query := `SELECT COUNT(*) FROM schema_migrations WHERE version = ? AND success = TRUE`

	var count int
	err := m.db.QueryRowContext(ctx, query, version).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// applyMigration 应用单个迁移
func (m *Migrator) applyMigration(ctx context.Context, migra *interfaces.Migration) error {
	// 开始事务
	tx, err := m.db.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 获取迁移SQL
	sql := migration.GetMigrationSQL(migrations, migra.Version)
	if sql == "" {
		return fmt.Errorf("no SQL found for migration %s", migra.Version)
	}

	// 执行迁移SQL
	_, err = tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// 记录成功的迁移
	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version, applied_at, success) VALUES (?, ?, TRUE)`,
		migra.Version,
		timeutils.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
