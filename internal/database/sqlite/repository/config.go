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

// NotificationChannelRepository 通知渠道数据访问实现
type NotificationChannelRepository struct {
	db *database.Database
}

// newNotificationChannelRepository 创建通知渠道仓库
func newNotificationChannelRepository(db *database.Database) interfaces.NotificationChannelRepository {
	return &NotificationChannelRepository{db: db}
}

// Create 创建通知渠道
func (r *NotificationChannelRepository) Create(ctx context.Context, channel *models.NotificationChannel) error {
	query := `INSERT INTO notification_channels (name, type, config, is_active, test_result, last_test, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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

	_, err := r.db.ExecContext(ctx, query,
		channel.Name,
		channel.Type,
		channel.Config,
		channel.IsActive,
		channel.TestResult,
		channel.LastTest,
		timeutils.Now(),
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

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

// List 获取通知渠道列表
func (r *NotificationChannelRepository) List(ctx context.Context, offset, limit int) ([]*models.NotificationChannel, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notification_channels ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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

	rows, err := r.db.QueryContext(ctx, query)
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

	rows, err := r.db.QueryContext(ctx, query, channelType)
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
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notification channels: %w", err)
	}

	return count, nil
}

// UpdateTestResult 更新测试结果
func (r *NotificationChannelRepository) UpdateTestResult(ctx context.Context, id int64, testResult string) error {
	query := `UPDATE notification_channels SET test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	now := timeutils.Now()
	_, err := r.db.ExecContext(ctx, query, testResult, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update notification channel test result: %w", err)
	}

	return nil
}
