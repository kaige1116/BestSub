package migration

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/utils"
)

// MigrationFunc 迁移函数类型
type MigrationFunc func() string

// registeredMigrations 注册的迁移
var registeredMigrations = make(map[string]MigrationInfo)

// MigrationInfo 迁移信息
type MigrationInfo struct {
	Description string
	Func        MigrationFunc
}

// RegisterMigration 注册迁移
func RegisterMigration(version, description string, fn MigrationFunc) {
	registeredMigrations[version] = MigrationInfo{
		Description: description,
		Func:        fn,
	}
}

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
	if err := m.validate(); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	// 确保迁移表存在
	if err := m.createMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// 获取所有迁移
	migrations := m.getAllMigrations()
	if len(migrations) == 0 {
		return nil // 没有迁移需要执行
	}

	// 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return compareVersions(migrations[i].Version, migrations[j].Version) < 0
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

// validate 验证迁移脚本
func (m *Migrator) validate() error {
	migrations := m.getAllMigrations()

	// 检查版本号是否重复
	versions := make(map[string]bool)
	for _, migration := range migrations {
		if versions[migration.Version] {
			return fmt.Errorf("duplicate migration version: %s", migration.Version)
		}
		versions[migration.Version] = true
	}

	// 检查版本号格式
	for _, migration := range migrations {
		if !isValidVersion(migration.Version) {
			return fmt.Errorf("invalid version format: %s", migration.Version)
		}
	}

	// 检查每个迁移是否有对应的SQL
	for _, migration := range migrations {
		sql := m.getMigrationSQL(migration.Version)
		if sql == "" {
			return fmt.Errorf("no SQL found for migration %s", migration.Version)
		}
	}

	return nil
}

// getAllMigrations 获取所有注册的迁移
func (m *Migrator) getAllMigrations() []*interfaces.Migration {
	var migrations []*interfaces.Migration

	for version, info := range registeredMigrations {
		migrations = append(migrations, &interfaces.Migration{
			Version:     version,
			Description: info.Description,
		})
	}

	return migrations
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
func (m *Migrator) applyMigration(ctx context.Context, migration *interfaces.Migration) error {
	// 开始事务
	tx, err := m.db.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 获取迁移SQL
	sql := m.getMigrationSQL(migration.Version)
	if sql == "" {
		return fmt.Errorf("no SQL found for migration %s", migration.Version)
	}

	// 执行迁移SQL
	_, err = tx.Exec(sql)
	if err != nil {
		// 记录失败的迁移
		m.recordMigration(ctx, migration.Version, false, err.Error())
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// 记录成功的迁移
	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version, applied_at, success) VALUES (?, ?, TRUE)`,
		migration.Version,
		utils.Now(),
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

// recordMigration 记录迁移结果
func (m *Migrator) recordMigration(ctx context.Context, version string, success bool, errorMsg string) {
	query := `INSERT OR REPLACE INTO schema_migrations (version, applied_at, success, error) VALUES (?, ?, ?, ?)`
	m.db.ExecContext(ctx, query, version, utils.Now(), success, errorMsg)
}

// getMigrationSQL 根据版本号获取迁移SQL
func (m *Migrator) getMigrationSQL(version string) string {
	if info, exists := registeredMigrations[version]; exists {
		return info.Func()
	}
	return ""
}

// isValidVersion 检查版本号格式是否有效
func isValidVersion(version string) bool {
	if len(version) != 3 {
		return false
	}
	_, err := strconv.Atoi(version)
	return err == nil
}

// compareVersions 比较版本号
func compareVersions(v1, v2 string) int {
	n1, _ := strconv.Atoi(v1)
	n2, _ := strconv.Atoi(v2)

	if n1 < n2 {
		return -1
	} else if n1 > n2 {
		return 1
	}
	return 0
}
