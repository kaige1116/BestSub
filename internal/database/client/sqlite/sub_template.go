package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
)

type SubTemplateRepository struct {
	db *DB
}

func (db *DB) SubTemplate() interfaces.SubTemplateRepository {
	return &SubTemplateRepository{db: db}
}

func (r *SubTemplateRepository) Create(ctx context.Context, template *sub.Template) error {
	query := `INSERT INTO sub_template (name, type, template)
	          VALUES (?, ?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		template.Name,
		template.Type,
		template.Template,
	)

	if err != nil {
		return fmt.Errorf("failed to create sub template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get sub template id: %w", err)
	}

	template.ID = uint16(id)

	return nil
}

func (r *SubTemplateRepository) GetByID(ctx context.Context, id uint16) (*sub.Template, error) {
	query := `SELECT id, name, type, template
	          FROM sub_template WHERE id = ?`

	var template sub.Template
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Type,
		&template.Template,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub template by id: %w", err)
	}

	return &template, nil
}

func (r *SubTemplateRepository) Update(ctx context.Context, template *sub.Template) error {
	query := `UPDATE sub_template SET name = ?, type = ?, template = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		template.Name,
		template.Type,
		template.Template,
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub template: %w", err)
	}

	return nil
}

func (r *SubTemplateRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_template WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub template: %w", err)
	}

	return nil
}

func (r *SubTemplateRepository) List(ctx context.Context) (*[]sub.Template, error) {
	query := `SELECT id, name, type, template
	          FROM sub_template`

	var templates []sub.Template
	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sub templates: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var template sub.Template
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Type,
			&template.Template,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan sub template: %w", err)
		}
		templates = append(templates, template)
	}

	return &templates, nil
}
