package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
)

// SubStorageConfigRepository 存储配置数据访问实现
type SubStorageConfigRepository struct {
	db *DB
}

// newSubStorageConfigRepository 创建存储配置仓库
func (db *DB) SubStorage() interfaces.SubStorageConfigRepository {
	return &SubStorageConfigRepository{db: db}
}

// Create 创建存储配置
func (r *SubStorageConfigRepository) Create(ctx context.Context, config *sub.StorageConfig) error {
	query := `INSERT INTO storage_configs (enable, name, description, type, config, test_result, last_test, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		config.Enable,
		config.Name,
		config.Description,
		config.Type,
		config.Config,
		config.TestResult,
		config.LastTest,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create storage config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get storage config id: %w", err)
	}

	config.ID = uint16(id)
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取存储配置
func (r *SubStorageConfigRepository) GetByID(ctx context.Context, id uint16) (*sub.StorageConfig, error) {
	query := `SELECT id, enable, name, description, type, config, test_result, last_test, created_at, updated_at
	          FROM storage_configs WHERE id = ?`

	var config sub.StorageConfig
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.Enable,
		&config.Name,
		&config.Description,
		&config.Type,
		&config.Config,
		&config.TestResult,
		&config.LastTest,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get storage config by id: %w", err)
	}

	return &config, nil
}

// Update 更新存储配置
func (r *SubStorageConfigRepository) Update(ctx context.Context, config *sub.StorageConfig) error {
	query := `UPDATE storage_configs SET enable = ?, name = ?, description = ?, type = ?, config = ?,
	          test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		config.Enable,
		config.Name,
		config.Description,
		config.Type,
		config.Config,
		config.TestResult,
		config.LastTest,
		time.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update storage config: %w", err)
	}

	return nil
}

// Delete 删除存储配置
func (r *SubStorageConfigRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM storage_configs WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete storage config: %w", err)
	}

	return nil
}

// GetBySaveID 根据保存ID获取存储配置列表
func (r *SubStorageConfigRepository) GetBySaveID(ctx context.Context, saveID uint16) (*[]sub.StorageConfig, error) {
	query := `SELECT sc.id, sc.enable, sc.name, sc.description, sc.type, sc.config, sc.test_result, sc.last_test, sc.created_at, sc.updated_at
	          FROM storage_configs sc
	          INNER JOIN save_storage_relations ssr ON sc.id = ssr.storage_id
	          WHERE ssr.save_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, saveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage configs by save id: %w", err)
	}
	defer rows.Close()

	var configs []sub.StorageConfig
	for rows.Next() {
		var config sub.StorageConfig
		err := rows.Scan(
			&config.ID,
			&config.Enable,
			&config.Name,
			&config.Description,
			&config.Type,
			&config.Config,
			&config.TestResult,
			&config.LastTest,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan storage config: %w", err)
		}
		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate storage configs: %w", err)
	}

	return &configs, nil
}

// AddSaveRelation 添加存储配置与保存的关联
func (r *SubStorageConfigRepository) AddSaveRelation(ctx context.Context, configID, saveID uint16) error {
	query := `INSERT OR IGNORE INTO save_storage_relations (storage_id, save_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, configID, saveID)
	if err != nil {
		return fmt.Errorf("failed to add save relation: %w", err)
	}

	return nil
}

// SubOutputTemplateRepository 输出模板数据访问实现
type SubOutputTemplateRepository struct {
	db *DB
}

// newSubOutputTemplateRepository 创建输出模板仓库
func (db *DB) SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	return &SubOutputTemplateRepository{db: db}
}

// Create 创建输出模板
func (r *SubOutputTemplateRepository) Create(ctx context.Context, template *sub.OutputTemplate) error {
	query := `INSERT INTO sub_output_templates (enable, name, description, type, template, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		template.Enable,
		template.Name,
		template.Description,
		template.Type,
		template.Template,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create output template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get output template id: %w", err)
	}

	template.ID = uint16(id)
	template.CreatedAt = now
	template.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取输出模板
func (r *SubOutputTemplateRepository) GetByID(ctx context.Context, id uint16) (*sub.OutputTemplate, error) {
	query := `SELECT id, enable, name, description, type, template, created_at, updated_at
	          FROM sub_output_templates WHERE id = ?`

	var template sub.OutputTemplate
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Enable,
		&template.Name,
		&template.Description,
		&template.Type,
		&template.Template,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get output template by id: %w", err)
	}

	return &template, nil
}

// Update 更新输出模板
func (r *SubOutputTemplateRepository) Update(ctx context.Context, template *sub.OutputTemplate) error {
	query := `UPDATE sub_output_templates SET enable = ?, name = ?, description = ?, type = ?, template = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		template.Enable,
		template.Name,
		template.Description,
		template.Type,
		template.Template,
		time.Now(),
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update output template: %w", err)
	}

	return nil
}

// Delete 删除输出模板
func (r *SubOutputTemplateRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_output_templates WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete output template: %w", err)
	}

	return nil
}

// GetBySaveID 根据保存ID获取输出模板
func (r *SubOutputTemplateRepository) GetBySaveID(ctx context.Context, saveID uint16) (*sub.OutputTemplate, error) {
	query := `SELECT ot.id, ot.enable, ot.name, ot.description, ot.type, ot.template, ot.created_at, ot.updated_at
	          FROM sub_output_templates ot
	          INNER JOIN save_template_relations str ON ot.id = str.template_id
	          WHERE str.save_id = ?`

	var template sub.OutputTemplate
	err := r.db.db.QueryRowContext(ctx, query, saveID).Scan(
		&template.ID,
		&template.Enable,
		&template.Name,
		&template.Description,
		&template.Type,
		&template.Template,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get output template by save id: %w", err)
	}

	return &template, nil
}

// GetByShareID 根据分享ID获取输出模板
func (r *SubOutputTemplateRepository) GetByShareID(ctx context.Context, shareID uint16) (*sub.OutputTemplate, error) {
	query := `SELECT ot.id, ot.enable, ot.name, ot.description, ot.type, ot.template, ot.created_at, ot.updated_at
	          FROM sub_output_templates ot
	          INNER JOIN share_template_relations str ON ot.id = str.template_id
	          WHERE str.share_id = ?`

	var template sub.OutputTemplate
	err := r.db.db.QueryRowContext(ctx, query, shareID).Scan(
		&template.ID,
		&template.Enable,
		&template.Name,
		&template.Description,
		&template.Type,
		&template.Template,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get output template by share id: %w", err)
	}

	return &template, nil
}

// AddShareRelation 添加输出模板与分享的关联
func (r *SubOutputTemplateRepository) AddShareRelation(ctx context.Context, templateID, shareID uint16) error {
	query := `INSERT OR IGNORE INTO share_template_relations (template_id, share_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, templateID, shareID)
	if err != nil {
		return fmt.Errorf("failed to add share relation: %w", err)
	}

	return nil
}

// AddSaveRelation 添加输出模板与保存的关联
func (r *SubOutputTemplateRepository) AddSaveRelation(ctx context.Context, templateID, saveID uint16) error {
	query := `INSERT OR IGNORE INTO save_template_relations (template_id, save_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, templateID, saveID)
	if err != nil {
		return fmt.Errorf("failed to add save relation: %w", err)
	}

	return nil
}

// SubNodeFilterRuleRepository 节点筛选规则数据访问实现
type SubNodeFilterRuleRepository struct {
	db *DB
}

// newSubNodeFilterRuleRepository 创建节点筛选规则仓库
func (db *DB) SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	return &SubNodeFilterRuleRepository{db: db}
}

// Create 创建筛选规则
func (r *SubNodeFilterRuleRepository) Create(ctx context.Context, rule *sub.NodeFilterRule) error {
	query := `INSERT INTO sub_node_filter_rules (name, field, operator, value, description, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		rule.Name,
		rule.Field,
		rule.Operator,
		rule.Value,
		rule.Description,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create node filter rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get node filter rule id: %w", err)
	}

	rule.ID = uint16(id)
	rule.CreatedAt = now
	rule.UpdatedAt = now

	return nil
}

// Update 更新筛选规则
func (r *SubNodeFilterRuleRepository) Update(ctx context.Context, rule *sub.NodeFilterRule) error {
	query := `UPDATE sub_node_filter_rules SET name = ?, field = ?, operator = ?, value = ?, description = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		rule.Name,
		rule.Field,
		rule.Operator,
		rule.Value,
		rule.Description,
		time.Now(),
		rule.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update node filter rule: %w", err)
	}

	return nil
}

// GetByID 根据ID获取筛选规则
func (r *SubNodeFilterRuleRepository) GetByID(ctx context.Context, id uint16) (*sub.NodeFilterRule, error) {
	query := `SELECT id, name, field, operator, value, description, created_at, updated_at
	          FROM sub_node_filter_rules WHERE id = ?`

	var rule sub.NodeFilterRule
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.Name,
		&rule.Field,
		&rule.Operator,
		&rule.Value,
		&rule.Description,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get node filter rule by id: %w", err)
	}

	return &rule, nil
}

// Delete 删除筛选规则
func (r *SubNodeFilterRuleRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_node_filter_rules WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete node filter rule: %w", err)
	}

	return nil
}

// GetBySaveID 根据保存ID获取筛选规则
func (r *SubNodeFilterRuleRepository) GetBySaveID(ctx context.Context, saveID uint16) (*[]sub.NodeFilterRule, error) {
	query := `SELECT nfr.id, nfr.name, nfr.field, nfr.operator, nfr.value, nfr.description, nfr.created_at, nfr.updated_at
	          FROM sub_node_filter_rules nfr
	          INNER JOIN save_filter_relations sfr ON nfr.id = sfr.filter_id
	          WHERE sfr.save_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, saveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node filter rules by save id: %w", err)
	}
	defer rows.Close()

	var rules []sub.NodeFilterRule
	for rows.Next() {
		var rule sub.NodeFilterRule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Field,
			&rule.Operator,
			&rule.Value,
			&rule.Description,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node filter rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate node filter rules: %w", err)
	}

	return &rules, nil
}

// GetByShareID 根据分享ID获取筛选规则
func (r *SubNodeFilterRuleRepository) GetByShareID(ctx context.Context, shareID uint16) (*[]sub.NodeFilterRule, error) {
	query := `SELECT nfr.id, nfr.name, nfr.field, nfr.operator, nfr.value, nfr.description, nfr.created_at, nfr.updated_at
	          FROM sub_node_filter_rules nfr
	          INNER JOIN share_filter_relations sfr ON nfr.id = sfr.filter_id
	          WHERE sfr.share_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, shareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node filter rules by share id: %w", err)
	}
	defer rows.Close()

	var rules []sub.NodeFilterRule
	for rows.Next() {
		var rule sub.NodeFilterRule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Field,
			&rule.Operator,
			&rule.Value,
			&rule.Description,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node filter rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate node filter rules: %w", err)
	}

	return &rules, nil
}

// AddShareRelation 添加筛选规则与分享的关联
func (r *SubNodeFilterRuleRepository) AddShareRelation(ctx context.Context, ruleID, shareID uint16) error {
	query := `INSERT OR IGNORE INTO share_filter_relations (filter_id, share_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, ruleID, shareID)
	if err != nil {
		return fmt.Errorf("failed to add share relation: %w", err)
	}

	return nil
}

// AddSaveRelation 添加筛选规则与保存的关联
func (r *SubNodeFilterRuleRepository) AddSaveRelation(ctx context.Context, ruleID, saveID uint16) error {
	query := `INSERT OR IGNORE INTO save_filter_relations (filter_id, save_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, ruleID, saveID)
	if err != nil {
		return fmt.Errorf("failed to add save relation: %w", err)
	}

	return nil
}
