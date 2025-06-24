package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/utils"
)

// SubStorageConfigRepository 存储配置数据访问实现
type SubStorageConfigRepository struct {
	db *database.Database
}

// newSubStorageConfigRepository 创建存储配置仓库
func newSubStorageConfigRepository(db *database.Database) interfaces.SubStorageConfigRepository {
	return &SubStorageConfigRepository{db: db}
}

// Create 创建存储配置
func (r *SubStorageConfigRepository) Create(ctx context.Context, config *models.SubStorageConfig) error {
	query := `INSERT INTO sub_storage_configs (name, type, config, is_active, test_result, last_test, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := utils.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.Name,
		config.Type,
		config.Config,
		config.IsActive,
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

	config.ID = id
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取存储配置
func (r *SubStorageConfigRepository) GetByID(ctx context.Context, id int64) (*models.SubStorageConfig, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM sub_storage_configs WHERE id = ?`

	var config models.SubStorageConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.Name,
		&config.Type,
		&config.Config,
		&config.IsActive,
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
func (r *SubStorageConfigRepository) Update(ctx context.Context, config *models.SubStorageConfig) error {
	query := `UPDATE sub_storage_configs SET name = ?, type = ?, config = ?, is_active = ?, 
	          test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		config.Name,
		config.Type,
		config.Config,
		config.IsActive,
		config.TestResult,
		config.LastTest,
		utils.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update storage config: %w", err)
	}

	return nil
}

// Delete 删除存储配置
func (r *SubStorageConfigRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_storage_configs WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete storage config: %w", err)
	}

	return nil
}

// List 获取存储配置列表
func (r *SubStorageConfigRepository) List(ctx context.Context, offset, limit int) ([]*models.SubStorageConfig, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM sub_storage_configs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.queryStorageConfigs(ctx, query, limit, offset)
}

// ListActive 获取活跃的存储配置列表
func (r *SubStorageConfigRepository) ListActive(ctx context.Context) ([]*models.SubStorageConfig, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM sub_storage_configs WHERE is_active = true ORDER BY created_at DESC`

	return r.queryStorageConfigs(ctx, query)
}

// ListByType 根据类型获取存储配置列表
func (r *SubStorageConfigRepository) ListByType(ctx context.Context, storageType string) ([]*models.SubStorageConfig, error) {
	query := `SELECT id, name, type, config, is_active, test_result, last_test, created_at, updated_at 
	          FROM sub_storage_configs WHERE type = ? ORDER BY created_at DESC`

	return r.queryStorageConfigs(ctx, query, storageType)
}

// Count 获取存储配置总数
func (r *SubStorageConfigRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_storage_configs`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count storage configs: %w", err)
	}

	return count, nil
}

// UpdateTestResult 更新测试结果
func (r *SubStorageConfigRepository) UpdateTestResult(ctx context.Context, id int64, testResult string) error {
	query := `UPDATE sub_storage_configs SET test_result = ?, last_test = ?, updated_at = ? WHERE id = ?`

	now := utils.Now()
	_, err := r.db.ExecContext(ctx, query, testResult, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update storage config test result: %w", err)
	}

	return nil
}

// queryStorageConfigs 通用存储配置查询方法
func (r *SubStorageConfigRepository) queryStorageConfigs(ctx context.Context, query string, args ...interface{}) ([]*models.SubStorageConfig, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query storage configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.SubStorageConfig
	for rows.Next() {
		var config models.SubStorageConfig
		err := rows.Scan(
			&config.ID,
			&config.Name,
			&config.Type,
			&config.Config,
			&config.IsActive,
			&config.TestResult,
			&config.LastTest,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan storage config: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate storage configs: %w", err)
	}

	return configs, nil
}

// SubOutputTemplateRepository 输出模板数据访问实现
type SubOutputTemplateRepository struct {
	db *database.Database
}

// newSubOutputTemplateRepository 创建输出模板仓库
func newSubOutputTemplateRepository(db *database.Database) interfaces.SubOutputTemplateRepository {
	return &SubOutputTemplateRepository{db: db}
}

// Create 创建输出模板
func (r *SubOutputTemplateRepository) Create(ctx context.Context, template *models.SubOutputTemplate) error {
	query := `INSERT INTO sub_output_templates (format, version, template, description, is_default, is_active, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := utils.Now()
	result, err := r.db.ExecContext(ctx, query,
		template.Format,
		template.Version,
		template.Template,
		template.Description,
		template.IsDefault,
		template.IsActive,
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

	template.ID = id
	template.CreatedAt = now
	template.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取输出模板
func (r *SubOutputTemplateRepository) GetByID(ctx context.Context, id int64) (*models.SubOutputTemplate, error) {
	query := `SELECT id, format, version, template, description, is_default, is_active, created_at, updated_at 
	          FROM sub_output_templates WHERE id = ?`

	var template models.SubOutputTemplate
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Format,
		&template.Version,
		&template.Template,
		&template.Description,
		&template.IsDefault,
		&template.IsActive,
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
func (r *SubOutputTemplateRepository) Update(ctx context.Context, template *models.SubOutputTemplate) error {
	query := `UPDATE sub_output_templates SET format = ?, version = ?, template = ?, description = ?, 
	          is_default = ?, is_active = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		template.Format,
		template.Version,
		template.Template,
		template.Description,
		template.IsDefault,
		template.IsActive,
		utils.Now(),
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update output template: %w", err)
	}

	return nil
}

// Delete 删除输出模板
func (r *SubOutputTemplateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_output_templates WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete output template: %w", err)
	}

	return nil
}

// List 获取输出模板列表
func (r *SubOutputTemplateRepository) List(ctx context.Context, offset, limit int) ([]*models.SubOutputTemplate, error) {
	query := `SELECT id, format, version, template, description, is_default, is_active, created_at, updated_at 
	          FROM sub_output_templates ORDER BY format, version LIMIT ? OFFSET ?`

	return r.queryOutputTemplates(ctx, query, limit, offset)
}

// ListActive 获取活跃的输出模板列表
func (r *SubOutputTemplateRepository) ListActive(ctx context.Context) ([]*models.SubOutputTemplate, error) {
	query := `SELECT id, format, version, template, description, is_default, is_active, created_at, updated_at 
	          FROM sub_output_templates WHERE is_active = true ORDER BY format, version`

	return r.queryOutputTemplates(ctx, query)
}

// ListByFormat 根据格式获取输出模板列表
func (r *SubOutputTemplateRepository) ListByFormat(ctx context.Context, format string) ([]*models.SubOutputTemplate, error) {
	query := `SELECT id, format, version, template, description, is_default, is_active, created_at, updated_at 
	          FROM sub_output_templates WHERE format = ? ORDER BY version`

	return r.queryOutputTemplates(ctx, query, format)
}

// GetDefault 获取默认模板
func (r *SubOutputTemplateRepository) GetDefault(ctx context.Context, format string) (*models.SubOutputTemplate, error) {
	query := `SELECT id, format, version, template, description, is_default, is_active, created_at, updated_at 
	          FROM sub_output_templates WHERE format = ? AND is_default = true LIMIT 1`

	var template models.SubOutputTemplate
	err := r.db.QueryRowContext(ctx, query, format).Scan(
		&template.ID,
		&template.Format,
		&template.Version,
		&template.Template,
		&template.Description,
		&template.IsDefault,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get default template: %w", err)
	}

	return &template, nil
}

// SetDefault 设置默认模板
func (r *SubOutputTemplateRepository) SetDefault(ctx context.Context, id int64, format string) error {
	tx, err := r.db.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 清除同格式的其他默认模板
	clearQuery := `UPDATE sub_output_templates SET is_default = false WHERE format = ?`
	_, err = tx.Exec(clearQuery, format)
	if err != nil {
		return fmt.Errorf("failed to clear default templates: %w", err)
	}

	// 设置新的默认模板
	setQuery := `UPDATE sub_output_templates SET is_default = true WHERE id = ?`
	_, err = tx.Exec(setQuery, id)
	if err != nil {
		return fmt.Errorf("failed to set default template: %w", err)
	}

	return tx.Commit()
}

// Count 获取输出模板总数
func (r *SubOutputTemplateRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_output_templates`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count output templates: %w", err)
	}

	return count, nil
}

// queryOutputTemplates 通用输出模板查询方法
func (r *SubOutputTemplateRepository) queryOutputTemplates(ctx context.Context, query string, args ...interface{}) ([]*models.SubOutputTemplate, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query output templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.SubOutputTemplate
	for rows.Next() {
		var template models.SubOutputTemplate
		err := rows.Scan(
			&template.ID,
			&template.Format,
			&template.Version,
			&template.Template,
			&template.Description,
			&template.IsDefault,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan output template: %w", err)
		}
		templates = append(templates, &template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate output templates: %w", err)
	}

	return templates, nil
}

// SubNodeFilterRuleRepository 节点筛选规则数据访问实现
type SubNodeFilterRuleRepository struct {
	db *database.Database
}

// newSubNodeFilterRuleRepository 创建节点筛选规则仓库
func newSubNodeFilterRuleRepository(db *database.Database) interfaces.SubNodeFilterRuleRepository {
	return &SubNodeFilterRuleRepository{db: db}
}

// Create 创建筛选规则
func (r *SubNodeFilterRuleRepository) Create(ctx context.Context, rule *models.SubNodeFilterRule) error {
	query := `INSERT INTO sub_node_filter_rules (rule_type, operator, value, is_enabled, priority, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := utils.Now()
	result, err := r.db.ExecContext(ctx, query,
		rule.RuleType,
		rule.Operator,
		rule.Value,
		rule.IsEnabled,
		rule.Priority,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create filter rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get filter rule id: %w", err)
	}

	rule.ID = id
	rule.CreatedAt = now
	rule.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取筛选规则
func (r *SubNodeFilterRuleRepository) GetByID(ctx context.Context, id int64) (*models.SubNodeFilterRule, error) {
	query := `SELECT id, rule_type, operator, value, is_enabled, priority, created_at, updated_at 
	          FROM sub_node_filter_rules WHERE id = ?`

	var rule models.SubNodeFilterRule
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.RuleType,
		&rule.Operator,
		&rule.Value,
		&rule.IsEnabled,
		&rule.Priority,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get filter rule by id: %w", err)
	}

	return &rule, nil
}

// Update 更新筛选规则
func (r *SubNodeFilterRuleRepository) Update(ctx context.Context, rule *models.SubNodeFilterRule) error {
	query := `UPDATE sub_node_filter_rules SET rule_type = ?, operator = ?, value = ?, is_enabled = ?, 
	          priority = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		rule.RuleType,
		rule.Operator,
		rule.Value,
		rule.IsEnabled,
		rule.Priority,
		utils.Now(),
		rule.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update filter rule: %w", err)
	}

	return nil
}

// Delete 删除筛选规则
func (r *SubNodeFilterRuleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_node_filter_rules WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete filter rule: %w", err)
	}

	return nil
}

// List 获取筛选规则列表
func (r *SubNodeFilterRuleRepository) List(ctx context.Context, offset, limit int) ([]*models.SubNodeFilterRule, error) {
	query := `SELECT id, rule_type, operator, value, is_enabled, priority, created_at, updated_at 
	          FROM sub_node_filter_rules ORDER BY priority ASC, created_at ASC LIMIT ? OFFSET ?`

	return r.queryFilterRules(ctx, query, limit, offset)
}

// ListEnabled 获取启用的筛选规则列表
func (r *SubNodeFilterRuleRepository) ListEnabled(ctx context.Context) ([]*models.SubNodeFilterRule, error) {
	query := `SELECT id, rule_type, operator, value, is_enabled, priority, created_at, updated_at 
	          FROM sub_node_filter_rules WHERE is_enabled = true ORDER BY priority ASC, created_at ASC`

	return r.queryFilterRules(ctx, query)
}

// ListByType 根据类型获取筛选规则列表
func (r *SubNodeFilterRuleRepository) ListByType(ctx context.Context, ruleType string) ([]*models.SubNodeFilterRule, error) {
	query := `SELECT id, rule_type, operator, value, is_enabled, priority, created_at, updated_at 
	          FROM sub_node_filter_rules WHERE rule_type = ? ORDER BY priority ASC, created_at ASC`

	return r.queryFilterRules(ctx, query, ruleType)
}

// Count 获取筛选规则总数
func (r *SubNodeFilterRuleRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_node_filter_rules`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count filter rules: %w", err)
	}

	return count, nil
}

// UpdatePriority 更新优先级
func (r *SubNodeFilterRuleRepository) UpdatePriority(ctx context.Context, id int64, priority int) error {
	query := `UPDATE sub_node_filter_rules SET priority = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, priority, utils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update filter rule priority: %w", err)
	}

	return nil
}

// queryFilterRules 通用筛选规则查询方法
func (r *SubNodeFilterRuleRepository) queryFilterRules(ctx context.Context, query string, args ...interface{}) ([]*models.SubNodeFilterRule, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query filter rules: %w", err)
	}
	defer rows.Close()

	var rules []*models.SubNodeFilterRule
	for rows.Next() {
		var rule models.SubNodeFilterRule
		err := rows.Scan(
			&rule.ID,
			&rule.RuleType,
			&rule.Operator,
			&rule.Value,
			&rule.IsEnabled,
			&rule.Priority,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan filter rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate filter rules: %w", err)
	}

	return rules, nil
}
