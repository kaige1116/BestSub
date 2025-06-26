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
	query := `INSERT INTO system_configs (key, value, type, group_name, description, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.Key,
		config.Value,
		config.Type,
		config.Group,
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
	          FROM system_configs WHERE key = ?`

	var config models.SystemConfig
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&config.ID,
		&config.Key,
		&config.Value,
		&config.Type,
		&config.Group,
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
	query := `UPDATE system_configs SET key = ?, value = ?, type = ?, group_name = ?, 
	          description = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		config.Key,
		config.Value,
		config.Type,
		config.Group,
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
	query := `DELETE FROM system_configs WHERE key = ?`

	_, err := r.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("failed to delete system config by key: %w", err)
	}

	return nil
}

// SetValue 设置配置值
func (r *SystemConfigRepository) SetValue(ctx context.Context, key, value, configType, group, description string) error {
	// 首先尝试更新
	updateQuery := `UPDATE system_configs SET value = ?, type = ?, group_name = ?, description = ?, updated_at = ? WHERE key = ?`
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
		insertQuery := `INSERT INTO system_configs (key, value, type, group_name, description, created_at, updated_at) 
		                VALUES (?, ?, ?, ?, ?, ?, ?)`
		now := timeutils.Now()
		_, err = r.db.ExecContext(ctx, insertQuery, key, value, configType, group, description, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert system config value: %w", err)
		}
	}

	return nil
}
