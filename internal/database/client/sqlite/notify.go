package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// NotifyRepository 通知渠道数据访问实现
type NotifyRepository struct {
	db *DB
}

// newNotificationChannelRepository 创建通知渠道仓库
func (db *DB) Notify() interfaces.NotifyRepository {
	return &NotifyRepository{db: db}
}

// Create 创建通知渠道
func (r *NotifyRepository) Create(ctx context.Context, channel *notify.Data) error {
	log.Debugf("Create notify config")
	query := `INSERT INTO notify_config (type, config )
	          VALUES (?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		channel.Type,
		channel.Config,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification channel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get notification channel id: %w", err)
	}

	channel.ID = uint16(id)

	return nil
}

// GetByID 根据ID获取通知渠道
func (r *NotifyRepository) GetByID(ctx context.Context, id uint16) (*notify.Data, error) {
	log.Debugf("Get notify config by id")
	query := `SELECT id, type, config
	          FROM notify_config WHERE id = ?`

	var channel notify.Data
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&channel.ID,
		&channel.Type,
		&channel.Config,
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
func (r *NotifyRepository) Update(ctx context.Context, channel *notify.Data) error {
	log.Debugf("Update notify config")
	query := `UPDATE notify_config SET type = ?, config = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		channel.Type,
		channel.Config,
		channel.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification channel: %w", err)
	}

	return nil
}

// Delete 删除通知渠道
func (r *NotifyRepository) Delete(ctx context.Context, id uint16) error {
	log.Debugf("Delete notify config")
	query := `DELETE FROM notify_config WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

// List 获取通知渠道列表
func (r *NotifyRepository) List(ctx context.Context) (*[]notify.Data, error) {
	log.Debugf("List notify config")
	query := `SELECT id, type, config
	          FROM notify_config ORDER BY id DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels: %w", err)
	}
	defer rows.Close()

	var channels []notify.Data
	for rows.Next() {
		var channel notify.Data
		err := rows.Scan(
			&channel.ID,
			&channel.Type,
			&channel.Config,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return &channels, nil
}

// Count 获取通知渠道总数
func (r *NotifyRepository) Count(ctx context.Context) (uint16, error) {
	query := `SELECT COUNT(*) FROM notify_config`

	var count uint16
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notification channels: %w", err)
	}

	return count, nil
}

// GetByTaskID 根据任务ID获取通知渠道列表
func (r *NotifyRepository) GetByTaskID(ctx context.Context, taskID uint16) (*[]notify.Data, error) {
	query := `SELECT n.id, n.type, n.config
	          FROM notify_config n
	          INNER JOIN notify_task_relations ntr ON n.id = ntr.notify_id
	          WHERE ntr.task_id = ? ORDER BY n.created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification channels by task id: %w", err)
	}
	defer rows.Close()

	var channels []notify.Data
	for rows.Next() {
		var channel notify.Data
		err := rows.Scan(
			&channel.ID,
			&channel.Type,
			&channel.Config,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return &channels, nil
}

// AddTaskRelation 添加通知渠道与任务的关联
func (r *NotifyRepository) AddTaskRelation(ctx context.Context, notifyID, taskID uint16) error {
	query := `INSERT OR IGNORE INTO notify_task_relations (task_id, notify_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, taskID, notifyID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}

// NotifyTemplateRepository 通知模板数据访问实现
type NotifyTemplateRepository struct {
	db *DB
}

// newNotifyTemplateRepository 创建通知模板仓库
func (db *DB) NotifyTemplate() interfaces.NotifyTemplateRepository {
	return &NotifyTemplateRepository{db: db}
}

// Create 创建通知模板
func (r *NotifyTemplateRepository) Create(ctx context.Context, template *notify.Template) error {
	query := `INSERT INTO notify_templates (name, templates)
	          VALUES (?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		template.Name,
		template.Template,
	)

	if err != nil {
		return fmt.Errorf("failed to create notify template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get notify template id: %w", err)
	}

	template.ID = uint16(id)

	return nil
}

// GetByID 根据ID获取通知模板
func (r *NotifyTemplateRepository) GetByID(ctx context.Context, id uint16) (*notify.Template, error) {
	query := `SELECT id, name, templates
	          FROM notify_templates WHERE id = ?`

	var template notify.Template
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Template,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notify template by id: %w", err)
	}

	return &template, nil
}

// Update 更新通知模板
func (r *NotifyTemplateRepository) Update(ctx context.Context, template *notify.Template) error {
	query := `UPDATE notify_templates SET name = ?, templates = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		template.Name,
		template.Template,
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notify template: %w", err)
	}

	return nil
}

// Delete 删除通知模板
func (r *NotifyTemplateRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM notify_templates WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notify template: %w", err)
	}

	return nil
}

// List 获取通知模板列表
func (r *NotifyTemplateRepository) List(ctx context.Context) (*[]notify.Template, error) {
	query := `SELECT id, name, templates
	          FROM notify_templates ORDER BY id DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notify templates: %w", err)
	}
	defer rows.Close()

	var templates []notify.Template
	for rows.Next() {
		var template notify.Template
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Template,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notify template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notify templates: %w", err)
	}

	return &templates, nil
}

// Count 获取通知模板总数
func (r *NotifyTemplateRepository) Count(ctx context.Context) (uint16, error) {
	query := `SELECT COUNT(*) FROM notify_templates`

	var count uint16
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notify templates: %w", err)
	}

	return count, nil
}

// GetByTaskID 根据任务ID获取通知模板
func (r *NotifyTemplateRepository) GetByTaskID(ctx context.Context, taskID uint16) (*notify.Template, error) {
	query := `SELECT nt.id, nt.name, nt.templates
	          FROM notify_templates nt
	          INNER JOIN task_notify_template_relations tntr ON nt.id = tntr.notify_template_id
	          WHERE tntr.task_id = ?`

	var template notify.Template
	err := r.db.db.QueryRowContext(ctx, query, taskID).Scan(
		&template.ID,
		&template.Name,
		&template.Template,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notify template by task id: %w", err)
	}

	return &template, nil
}

// AddTaskRelation 添加通知模板与任务的关联
func (r *NotifyTemplateRepository) AddTaskRelation(ctx context.Context, templateID, taskID uint16) error {
	query := `INSERT OR IGNORE INTO task_notify_template_relations (task_id, notify_template_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, taskID, templateID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}
