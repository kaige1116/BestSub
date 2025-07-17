package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/system"
)

func (db *DB) Config() interfaces.ConfigRepository {
	return &SystemConfigRepository{db: db}
}

// SystemConfigRepository 系统配置数据访问实现
type SystemConfigRepository struct {
	db *DB
}

// Create 批量创建配置
func (r *SystemConfigRepository) Create(ctx context.Context, configs *[]system.Data) error {
	if configs == nil || len(*configs) == 0 {
		return nil
	}

	// 开始事务
	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO system_config (key, group_name, value, description)
	          VALUES (?, ?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 批量执行插入
	for i := range *configs {
		config := &(*configs)[i]
		_, err := stmt.ExecContext(ctx,
			config.Key,
			config.GroupName,
			config.Value,
			config.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to create system config key '%s': %w", config.Key, err)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByKey 根据键获取配置
func (r *SystemConfigRepository) GetByKey(ctx context.Context, key string) (*system.Data, error) {
	query := `SELECT key, group_name, value, description
	          FROM system_config WHERE key = ?`

	var config system.Data
	err := r.db.db.QueryRowContext(ctx, query, key).Scan(
		&config.Key,
		&config.GroupName,
		&config.Value,
		&config.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get system config by key: %w", err)
	}

	return &config, nil
}

// Update 批量更新配置（根据key更新value）
func (r *SystemConfigRepository) Update(ctx context.Context, configs *[]system.Data) error {
	if configs == nil || len(*configs) == 0 {
		return nil
	}

	// 开始事务
	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE system_config SET value = ? WHERE key = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 批量执行更新
	for _, config := range *configs {
		result, err := stmt.ExecContext(ctx,
			config.Value,
			config.Key,
		)
		if err != nil {
			return fmt.Errorf("failed to update system config key '%s': %w", config.Key, err)
		}

		// 检查是否有记录被更新
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for key '%s': %w", config.Key, err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("no config found with key '%s'", config.Key)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetConfigsByGroup 获取指定分组下的所有配置
func (r *SystemConfigRepository) GetConfigsByGroup(ctx context.Context, group string) ([]system.Data, error) {
	query := `SELECT key, group_name, value, description
	          FROM system_config WHERE group_name = ? ORDER BY key`

	rows, err := r.db.db.QueryContext(ctx, query, group)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs by group: %w", err)
	}
	defer rows.Close()

	var configs []system.Data
	for rows.Next() {
		var config system.Data
		if err := rows.Scan(
			&config.Key,
			&config.GroupName,
			&config.Value,
			&config.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate configs: %w", err)
	}

	return configs, nil
}
