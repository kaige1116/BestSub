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

// NotifyRepository 通知渠道数据访问实现
type NotifyRepository struct {
	db *database.Database
}

// newNotificationChannelRepository 创建通知渠道仓库
func newNotifyRepository(db *database.Database) interfaces.NotifyRepository {
	return &NotifyRepository{db: db}
}

// Create 创建通知渠道
func (r *NotifyRepository) Create(ctx context.Context, channel *models.Notify) error {
	query := `INSERT INTO notify (name, type, config, is_active, test_result, last_test, created_at, updated_at) 
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
func (r *NotifyRepository) GetByID(ctx context.Context, id int64) (*models.Notify, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notify WHERE id = ?`

	var channel models.Notify
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
func (r *NotifyRepository) Update(ctx context.Context, channel *models.Notify) error {
	query := `UPDATE notify SET name = ?, type = ?, config = ?, is_active = ?, 
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
func (r *NotifyRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notify WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

// List 获取通知渠道列表
func (r *NotifyRepository) List(ctx context.Context, offset, limit int) ([]*models.Notify, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notify ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.Notify
	for rows.Next() {
		var channel models.Notify
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
func (r *NotifyRepository) ListActive(ctx context.Context) ([]*models.Notify, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notify WHERE is_active = true ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active notification channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.Notify
	for rows.Next() {
		var channel models.Notify
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
func (r *NotifyRepository) ListByType(ctx context.Context, channelType string) ([]*models.Notify, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM notify WHERE type = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, channelType)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels by type: %w", err)
	}
	defer rows.Close()

	var channels []*models.Notify
	for rows.Next() {
		var channel models.Notify
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
func (r *NotifyRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM notify`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notification channels: %w", err)
	}

	return count, nil
}

// UpdateTestResult 更新测试结果
func (r *NotifyRepository) UpdateTestResult(ctx context.Context, id int64, testResult string) error {
	query := `UPDATE notify SET test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	now := timeutils.Now()
	_, err := r.db.ExecContext(ctx, query, testResult, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update notification channel test result: %w", err)
	}

	return nil
}
