package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// SystemConfigRepository 系统配置数据访问实现
type SystemConfigRepository struct {
	db *database.Database
}

// NewSystemConfigRepository 创建系统配置仓库
func newSystemConfigRepository(db *database.Database) interfaces.SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// Create 创建配置
func (r *SystemConfigRepository) Create(ctx context.Context, config *models.SystemConfig) error {
	query := `INSERT INTO system_config (key, value, type, group_name, description, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.Key,
		config.Value,
		config.Type,
		config.GroupName,
		config.Description,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create system config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get system config id: %w", err)
	}

	config.ID = id
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByKey 根据键获取配置
func (r *SystemConfigRepository) GetByKey(ctx context.Context, key string) (*models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at 
	          FROM system_config WHERE key = ?`

	var config models.SystemConfig
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&config.ID,
		&config.Key,
		&config.Value,
		&config.Type,
		&config.GroupName,
		&config.Description,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get system config by key: %w", err)
	}

	return &config, nil
}

// Update 更新配置
func (r *SystemConfigRepository) Update(ctx context.Context, config *models.SystemConfig) error {
	query := `UPDATE system_config SET key = ?, value = ?, type = ?, group_name = ?, 
	          description = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		config.Key,
		config.Value,
		config.Type,
		config.GroupName,
		config.Description,
		timeutils.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update system config: %w", err)
	}

	return nil
}

// DeleteByKey 根据键删除配置
func (r *SystemConfigRepository) DeleteByKey(ctx context.Context, key string) error {
	query := `DELETE FROM system_config WHERE key = ?`

	_, err := r.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("failed to delete system config by key: %w", err)
	}

	return nil
}

// SetValue 设置配置值
func (r *SystemConfigRepository) SetValue(ctx context.Context, key, value, configType, group, description string) error {
	// 首先尝试更新
	updateQuery := `UPDATE system_config SET value = ?, type = ?, group_name = ?, description = ?, updated_at = ? WHERE key = ?`
	result, err := r.db.ExecContext(ctx, updateQuery, value, configType, group, description, timeutils.Now(), key)
	if err != nil {
		return fmt.Errorf("failed to update system config value: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	// 如果没有更新到记录，则插入新记录
	if rowsAffected == 0 {
		insertQuery := `INSERT INTO system_config (key, value, type, group_name, description, created_at, updated_at) 
		                VALUES (?, ?, ?, ?, ?, ?, ?)`
		now := timeutils.Now()
		_, err = r.db.ExecContext(ctx, insertQuery, key, value, configType, group, description, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert system config value: %w", err)
		}
	}

	return nil
}

// GetAllKeys 获取所有配置键
func (r *SystemConfigRepository) GetAllKeys(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT key FROM system_config ORDER BY key`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all config keys: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan config key: %w", err)
		}
		keys = append(keys, key)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate config keys: %w", err)
	}

	return keys, nil
}

// GetAllGroups 获取所有配置分组
func (r *SystemConfigRepository) GetAllGroups(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT group_name FROM system_config ORDER BY group_name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all config groups: %w", err)
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			return nil, fmt.Errorf("failed to scan config group: %w", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate config groups: %w", err)
	}

	return groups, nil
}

// GetConfigsByGroup 获取指定分组下的所有配置
func (r *SystemConfigRepository) GetConfigsByGroup(ctx context.Context, group string) ([]models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at
	          FROM system_config WHERE group_name = ? ORDER BY key`

	rows, err := r.db.QueryContext(ctx, query, group)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs by group: %w", err)
	}
	defer rows.Close()

	var configs []models.SystemConfig
	for rows.Next() {
		var config models.SystemConfig
		if err := rows.Scan(
			&config.ID,
			&config.Key,
			&config.Value,
			&config.Type,
			&config.GroupName,
			&config.Description,
			&config.CreatedAt,
			&config.UpdatedAt,
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
