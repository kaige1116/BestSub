package migration

import (
	"context"
	"database/sql"
)

// BaseMigration 基础迁移实现
type BaseMigration struct {
	id          string
	description string
	upSQL       string
	downSQL     string
}

// NewBaseMigration 创建基础迁移
func NewBaseMigration(id, description, upSQL, downSQL string) Migration {
	return &BaseMigration{
		id:          id,
		description: description,
		upSQL:       upSQL,
		downSQL:     downSQL,
	}
}

// ID 获取迁移ID
func (m *BaseMigration) ID() string {
	return m.id
}

// Description 获取迁移描述
func (m *BaseMigration) Description() string {
	return m.description
}

// Up 执行向上迁移
func (m *BaseMigration) Up(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, m.upSQL)
	return err
}

// Down 执行向下迁移
func (m *BaseMigration) Down(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, m.downSQL)
	return err
}

// CustomMigration 自定义迁移实现
type CustomMigration struct {
	id          string
	description string
	upFunc      func(ctx context.Context, tx *sql.Tx) error
	downFunc    func(ctx context.Context, tx *sql.Tx) error
}

// NewCustomMigration 创建自定义迁移
func NewCustomMigration(id, description string, upFunc, downFunc func(ctx context.Context, tx *sql.Tx) error) Migration {
	return &CustomMigration{
		id:          id,
		description: description,
		upFunc:      upFunc,
		downFunc:    downFunc,
	}
}

// ID 获取迁移ID
func (m *CustomMigration) ID() string {
	return m.id
}

// Description 获取迁移描述
func (m *CustomMigration) Description() string {
	return m.description
}

// Up 执行向上迁移
func (m *CustomMigration) Up(ctx context.Context, tx *sql.Tx) error {
	if m.upFunc == nil {
		return nil
	}
	return m.upFunc(ctx, tx)
}

// Down 执行向下迁移
func (m *CustomMigration) Down(ctx context.Context, tx *sql.Tx) error {
	if m.downFunc == nil {
		return nil
	}
	return m.downFunc(ctx, tx)
}

// BaseManager 基础迁移管理器实现
type BaseManager struct {
	migrator   Migrator
	migrations map[string]Migration
	orderedIDs []string
}

// NewManager 创建迁移管理器
func NewManager(migrator Migrator) Manager {
	return &BaseManager{
		migrator:   migrator,
		migrations: make(map[string]Migration),
		orderedIDs: make([]string, 0),
	}
}

// RegisterMigrations 注册迁移
func (m *BaseManager) RegisterMigrations(migrations ...Migration) {
	for _, migration := range migrations {
		if _, exists := m.migrations[migration.ID()]; !exists {
			m.migrations[migration.ID()] = migration
			m.orderedIDs = append(m.orderedIDs, migration.ID())
		}
	}
}

// GetMigrations 获取所有注册的迁移
func (m *BaseManager) GetMigrations() []Migration {
	result := make([]Migration, 0, len(m.orderedIDs))
	for _, id := range m.orderedIDs {
		if migration, exists := m.migrations[id]; exists {
			result = append(result, migration)
		}
	}
	return result
}

// GetMigration 根据ID获取迁移
func (m *BaseManager) GetMigration(id string) (Migration, bool) {
	migration, exists := m.migrations[id]
	return migration, exists
}

// GetPendingMigrations 获取待执行的迁移
func (m *BaseManager) GetPendingMigrations(ctx context.Context) ([]Migration, error) {
	appliedRecords, err := m.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, record := range appliedRecords {
		appliedMap[record.ID] = true
	}

	var pending []Migration
	for _, id := range m.orderedIDs {
		if !appliedMap[id] {
			if migration, exists := m.migrations[id]; exists {
				pending = append(pending, migration)
			}
		}
	}

	return pending, nil
}

// MigrateUp 执行向上迁移
func (m *BaseManager) MigrateUp(ctx context.Context) error {
	pending, err := m.GetPendingMigrations(ctx)
	if err != nil {
		return err
	}

	for _, migration := range pending {
		if err := m.migrator.ApplyMigration(ctx, migration); err != nil {
			return err
		}
	}

	return nil
}

// MigrateUpTo 迁移到指定版本
func (m *BaseManager) MigrateUpTo(ctx context.Context, targetID string) error {
	pending, err := m.GetPendingMigrations(ctx)
	if err != nil {
		return err
	}

	for _, migration := range pending {
		if err := m.migrator.ApplyMigration(ctx, migration); err != nil {
			return err
		}
		if migration.ID() == targetID {
			break
		}
	}

	return nil
}

// MigrateDown 执行向下迁移
func (m *BaseManager) MigrateDown(ctx context.Context, steps int) error {
	appliedRecords, err := m.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// 按应用时间倒序排列
	for i := len(appliedRecords) - 1; i >= 0 && steps > 0; i-- {
		record := appliedRecords[i]
		if migration, exists := m.migrations[record.ID]; exists {
			if err := m.migrator.RollbackMigration(ctx, migration); err != nil {
				return err
			}
			steps--
		}
	}

	return nil
}

// MigrateDownTo 回滚到指定版本
func (m *BaseManager) MigrateDownTo(ctx context.Context, targetID string) error {
	appliedRecords, err := m.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// 按应用时间倒序排列，回滚到目标版本之后的所有迁移
	for i := len(appliedRecords) - 1; i >= 0; i-- {
		record := appliedRecords[i]
		if record.ID == targetID {
			break
		}
		if migration, exists := m.migrations[record.ID]; exists {
			if err := m.migrator.RollbackMigration(ctx, migration); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetStatus 获取迁移状态
func (m *BaseManager) GetStatus(ctx context.Context) (*MigrationStatus, error) {
	appliedRecords, err := m.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}

	pending, err := m.GetPendingMigrations(ctx)
	if err != nil {
		return nil, err
	}

	var lastMigration *MigrationRecord
	if len(appliedRecords) > 0 {
		lastMigration = &appliedRecords[len(appliedRecords)-1]
	}

	return &MigrationStatus{
		TotalMigrations:   len(m.migrations),
		AppliedMigrations: len(appliedRecords),
		PendingMigrations: len(pending),
		LastMigration:     lastMigration,
		AppliedList:       appliedRecords,
		PendingList:       pending,
	}, nil
}

// Validate 验证迁移完整性
func (m *BaseManager) Validate(ctx context.Context) error {
	// 这里可以添加迁移完整性检查逻辑
	// 例如检查迁移ID的唯一性、顺序等
	return nil
}
