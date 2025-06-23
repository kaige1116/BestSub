package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/repository/interfaces"
)

// SystemConfigRepository 系统配置数据访问实现
type SystemConfigRepository struct {
	db *Database
}

// NewSystemConfigRepository 创建系统配置仓库
func NewSystemConfigRepository(db *Database) interfaces.SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// Create 创建配置
func (r *SystemConfigRepository) Create(ctx context.Context, config *models.SystemConfig) error {
	query := `INSERT INTO system_configs (key, value, type, group_name, description, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
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

// GetByID 根据ID获取配置
func (r *SystemConfigRepository) GetByID(ctx context.Context, id int64) (*models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at 
	          FROM system_configs WHERE id = ?`

	var config models.SystemConfig
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
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
		return nil, fmt.Errorf("failed to get system config by id: %w", err)
	}

	return &config, nil
}

// GetByKey 根据键获取配置
func (r *SystemConfigRepository) GetByKey(ctx context.Context, key string) (*models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at 
	          FROM system_configs WHERE key = ?`

	var config models.SystemConfig
	err := r.db.db.QueryRowContext(ctx, query, key).Scan(
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

	_, err := r.db.db.ExecContext(ctx, query,
		config.Key,
		config.Value,
		config.Type,
		config.Group,
		config.Description,
		time.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update system config: %w", err)
	}

	return nil
}

// Delete 删除配置
func (r *SystemConfigRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_configs WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete system config: %w", err)
	}

	return nil
}

// DeleteByKey 根据键删除配置
func (r *SystemConfigRepository) DeleteByKey(ctx context.Context, key string) error {
	query := `DELETE FROM system_configs WHERE key = ?`

	_, err := r.db.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("failed to delete system config by key: %w", err)
	}

	return nil
}

// List 获取配置列表
func (r *SystemConfigRepository) List(ctx context.Context, offset, limit int) ([]*models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at 
	          FROM system_configs ORDER BY group_name, key LIMIT ? OFFSET ?`

	rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list system configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.SystemConfig
	for rows.Next() {
		var config models.SystemConfig
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan system config: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate system configs: %w", err)
	}

	return configs, nil
}

// ListByGroup 根据分组获取配置列表
func (r *SystemConfigRepository) ListByGroup(ctx context.Context, group string) ([]*models.SystemConfig, error) {
	query := `SELECT id, key, value, type, group_name, description, created_at, updated_at 
	          FROM system_configs WHERE group_name = ? ORDER BY key`

	rows, err := r.db.db.QueryContext(ctx, query, group)
	if err != nil {
		return nil, fmt.Errorf("failed to list system configs by group: %w", err)
	}
	defer rows.Close()

	var configs []*models.SystemConfig
	for rows.Next() {
		var config models.SystemConfig
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan system config: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate system configs: %w", err)
	}

	return configs, nil
}

// Count 获取配置总数
func (r *SystemConfigRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM system_configs`

	var count int64
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count system configs: %w", err)
	}

	return count, nil
}

// SetValue 设置配置值
func (r *SystemConfigRepository) SetValue(ctx context.Context, key, value, configType, group, description string) error {
	// 首先尝试更新
	updateQuery := `UPDATE system_configs SET value = ?, type = ?, group_name = ?, description = ?, updated_at = ? WHERE key = ?`
	result, err := r.db.db.ExecContext(ctx, updateQuery, value, configType, group, description, time.Now(), key)
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
		now := time.Now()
		_, err = r.db.db.ExecContext(ctx, insertQuery, key, value, configType, group, description, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert system config value: %w", err)
		}
	}

	return nil
}

// GetValue 获取配置值
func (r *SystemConfigRepository) GetValue(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM system_configs WHERE key = ?`

	var value string
	err := r.db.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get system config value: %w", err)
	}

	return value, nil
}

// NotificationChannelRepository 通知渠道数据访问实现
type NotificationChannelRepository struct {
	db *Database
}

// NewNotificationChannelRepository 创建通知渠道仓库
func NewNotificationChannelRepository(db *Database) interfaces.NotificationChannelRepository {
	return &NotificationChannelRepository{db: db}
}

// Create 创建通知渠道
func (r *NotificationChannelRepository) Create(ctx context.Context, channel *models.NotificationChannel) error {
	query := `INSERT INTO notification_channels (name, type, config, is_active, test_result, last_test, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		channel.Name,
		channel.Type,
		channel.Config,
		channel.IsActive,
		channel.TestResult,
		channel.LastTest,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification channel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get notification channel id: %w", err)
	}

	channel.ID = id
	channel.CreatedAt = now
	channel.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取通知渠道
func (r *NotificationChannelRepository) GetByID(ctx context.Context, id int64) (*models.NotificationChannel, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notification_channels WHERE id = ?`

	var channel models.NotificationChannel
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Type,
		&channel.Config,
		&channel.IsActive,
		&channel.TestResult,
		&channel.LastTest,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification channel by id: %w", err)
	}

	return &channel, nil
}

// Update 更新通知渠道
func (r *NotificationChannelRepository) Update(ctx context.Context, channel *models.NotificationChannel) error {
	query := `UPDATE notification_channels SET name = ?, type = ?, config = ?, is_active = ?, 
	          test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		channel.Name,
		channel.Type,
		channel.Config,
		channel.IsActive,
		channel.TestResult,
		channel.LastTest,
		time.Now(),
		channel.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification channel: %w", err)
	}

	return nil
}

// Delete 删除通知渠道
func (r *NotificationChannelRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notification_channels WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

// List 获取通知渠道列表
func (r *NotificationChannelRepository) List(ctx context.Context, offset, limit int) ([]*models.NotificationChannel, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notification_channels ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.NotificationChannel
	for rows.Next() {
		var channel models.NotificationChannel
		err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.Config,
			&channel.IsActive,
			&channel.TestResult,
			&channel.LastTest,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return channels, nil
}

// ListActive 获取活跃的通知渠道列表
func (r *NotificationChannelRepository) ListActive(ctx context.Context) ([]*models.NotificationChannel, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notification_channels WHERE is_active = true ORDER BY created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active notification channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.NotificationChannel
	for rows.Next() {
		var channel models.NotificationChannel
		err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.Config,
			&channel.IsActive,
			&channel.TestResult,
			&channel.LastTest,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return channels, nil
}

// ListByType 根据类型获取通知渠道列表
func (r *NotificationChannelRepository) ListByType(ctx context.Context, channelType string) ([]*models.NotificationChannel, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notification_channels WHERE type = ? ORDER BY created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query, channelType)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels by type: %w", err)
	}
	defer rows.Close()

	var channels []*models.NotificationChannel
	for rows.Next() {
		var channel models.NotificationChannel
		err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.Config,
			&channel.IsActive,
			&channel.TestResult,
			&channel.LastTest,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, &channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return channels, nil
}

// Count 获取通知渠道总数
func (r *NotificationChannelRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM notification_channels`

	var count int64
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notification channels: %w", err)
	}

	return count, nil
}

// UpdateTestResult 更新测试结果
func (r *NotificationChannelRepository) UpdateTestResult(ctx context.Context, id int64, testResult string) error {
	query := `UPDATE notification_channels SET test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	now := time.Now()
	_, err := r.db.db.ExecContext(ctx, query, testResult, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update notification channel test result: %w", err)
	}

	return nil
}
