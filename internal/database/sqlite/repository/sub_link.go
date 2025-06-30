package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/parser"
	"github.com/bestruirui/bestsub/internal/models/sublink"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// SubLinkRepository 链接数据访问实现
type SubLinkRepository struct {
	db *database.Database
}

// newSubLinkRepository 创建链接仓库
func newSubLinkRepository(db *database.Database) interfaces.SubLinkRepository {
	return &SubLinkRepository{db: db}
}

// Create 创建链接
func (r *SubLinkRepository) Create(ctx context.Context, link *sublink.Data) error {
	// 序列化 detector 配置
	detectorJSON, err := json.Marshal(link.Detector)
	if err != nil {
		return fmt.Errorf("failed to marshal detector config: %w", err)
	}

	// 序列化 notify 配置
	notifyJSON, err := json.Marshal(link.Notify)
	if err != nil {
		return fmt.Errorf("failed to marshal notify config: %w", err)
	}

	query := `INSERT INTO sub_links (name, url, type, user_agent, proxy_enable, timeout, retries, is_enabled, detector, notify, cron_expr,
	          last_status, error_msg, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		link.Name,
		link.FetchConfig.URL,
		string(link.FetchConfig.Type),
		link.FetchConfig.UserAgent,
		link.FetchConfig.ProxyEnable,
		link.FetchConfig.Timeout,
		link.FetchConfig.Retries,
		link.IsEnabled,
		string(detectorJSON),
		string(notifyJSON),
		link.CronExpr,
		link.LastStatus,
		link.ErrorMsg,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create sub link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get sub link id: %w", err)
	}

	link.ID = id
	link.CreatedAt = now
	link.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取链接
func (r *SubLinkRepository) GetByID(ctx context.Context, id int64) (*sublink.Data, error) {
	query := `SELECT id, name, url, type, user_agent, proxy_enable, timeout, retries, is_enabled, detector, notify, cron_expr,
	          last_status, error_msg, created_at, updated_at
	          FROM sub_links WHERE id = ?`

	var link sublink.Data
	var typeStr string
	var detectorJSON string
	var notifyJSON string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.Name,
		&link.FetchConfig.URL,
		&typeStr,
		&link.FetchConfig.UserAgent,
		&link.FetchConfig.ProxyEnable,
		&link.FetchConfig.Timeout,
		&link.FetchConfig.Retries,
		&link.IsEnabled,
		&detectorJSON,
		&notifyJSON,
		&link.CronExpr,
		&link.LastStatus,
		&link.ErrorMsg,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub link by id: %w", err)
	}

	// 反序列化类型和配置
	link.FetchConfig.Type = parser.ParserType(typeStr)
	if detectorJSON != "" {
		if err := json.Unmarshal([]byte(detectorJSON), &link.Detector); err != nil {
			return nil, fmt.Errorf("failed to unmarshal detector config: %w", err)
		}
	}
	if notifyJSON != "" {
		if err := json.Unmarshal([]byte(notifyJSON), &link.Notify); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notify config: %w", err)
		}
	}

	return &link, nil
}

// GetByURL 根据URL获取链接
func (r *SubLinkRepository) GetByURL(ctx context.Context, url string) (*sublink.Data, error) {
	query := `SELECT id, name, url, type, user_agent, proxy_enable, timeout, retries, is_enabled, detector, notify, cron_expr,
	          last_status, error_msg, created_at, updated_at
	          FROM sub_links WHERE url = ?`

	var link sublink.Data
	var typeStr string
	var detectorJSON string
	var notifyJSON string
	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&link.ID,
		&link.Name,
		&link.FetchConfig.URL,
		&typeStr,
		&link.FetchConfig.UserAgent,
		&link.FetchConfig.ProxyEnable,
		&link.FetchConfig.Timeout,
		&link.FetchConfig.Retries,
		&link.IsEnabled,
		&detectorJSON,
		&notifyJSON,
		&link.CronExpr,
		&link.LastStatus,
		&link.ErrorMsg,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub link by url: %w", err)
	}

	// 反序列化类型和配置
	link.FetchConfig.Type = parser.ParserType(typeStr)
	if detectorJSON != "" {
		if err := json.Unmarshal([]byte(detectorJSON), &link.Detector); err != nil {
			return nil, fmt.Errorf("failed to unmarshal detector config: %w", err)
		}
	}
	if notifyJSON != "" {
		if err := json.Unmarshal([]byte(notifyJSON), &link.Notify); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notify config: %w", err)
		}
	}

	return &link, nil
}

// Update 更新链接
func (r *SubLinkRepository) Update(ctx context.Context, link *sublink.Data) error {
	// 序列化 detector 配置
	detectorJSON, err := json.Marshal(link.Detector)
	if err != nil {
		return fmt.Errorf("failed to marshal detector config: %w", err)
	}

	// 序列化 notify 配置
	notifyJSON, err := json.Marshal(link.Notify)
	if err != nil {
		return fmt.Errorf("failed to marshal notify config: %w", err)
	}

	query := `UPDATE sub_links SET name = ?, url = ?, type = ?, user_agent = ?, proxy_enable = ?, timeout = ?, retries = ?, is_enabled = ?,
	          detector = ?, notify = ?, cron_expr = ?, last_status = ?, error_msg = ?, updated_at = ? WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query,
		link.Name,
		link.FetchConfig.URL,
		string(link.FetchConfig.Type),
		link.FetchConfig.UserAgent,
		link.FetchConfig.ProxyEnable,
		link.FetchConfig.Timeout,
		link.FetchConfig.Retries,
		link.IsEnabled,
		string(detectorJSON),
		string(notifyJSON),
		link.CronExpr,
		link.LastStatus,
		link.ErrorMsg,
		timeutils.Now(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub link: %w", err)
	}

	return nil
}

// Delete 删除链接
func (r *SubLinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_links WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub link: %w", err)
	}

	return nil
}

// List 获取链接列表
func (r *SubLinkRepository) List(ctx context.Context, offset, limit int) ([]*sublink.Data, error) {
	query := `SELECT id, name, url, type, user_agent, proxy_enable, timeout, retries, is_enabled, detector, notify, cron_expr,
	          last_status, error_msg, created_at, updated_at
	          FROM sub_links ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.querySubLinks(ctx, query, limit, offset)
}

// Count 获取链接总数
func (r *SubLinkRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_links`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sub links: %w", err)
	}

	return count, nil
}

// querySubLinks 通用链接查询方法
func (r *SubLinkRepository) querySubLinks(ctx context.Context, query string, args ...any) ([]*sublink.Data, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sub links: %w", err)
	}
	defer rows.Close()

	var links []*sublink.Data
	for rows.Next() {
		var link sublink.Data
		var typeStr string
		var detectorJSON string
		var notifyJSON string
		err := rows.Scan(
			&link.ID,
			&link.Name,
			&link.FetchConfig.URL,
			&typeStr,
			&link.FetchConfig.UserAgent,
			&link.FetchConfig.ProxyEnable,
			&link.FetchConfig.Timeout,
			&link.FetchConfig.Retries,
			&link.IsEnabled,
			&detectorJSON,
			&notifyJSON,
			&link.CronExpr,
			&link.LastStatus,
			&link.ErrorMsg,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub link: %w", err)
		}

		// 反序列化类型和配置
		link.FetchConfig.Type = parser.ParserType(typeStr)
		if detectorJSON != "" {
			if err := json.Unmarshal([]byte(detectorJSON), &link.Detector); err != nil {
				return nil, fmt.Errorf("failed to unmarshal detector config: %w", err)
			}
		}
		if notifyJSON != "" {
			if err := json.Unmarshal([]byte(notifyJSON), &link.Notify); err != nil {
				return nil, fmt.Errorf("failed to unmarshal notify config: %w", err)
			}
		}

		links = append(links, &link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sub links: %w", err)
	}

	return links, nil
}
