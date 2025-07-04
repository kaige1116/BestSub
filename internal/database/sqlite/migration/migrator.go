package migration

import (
	"context"
	"fmt"
	"sort"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/migration"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

var migrations = make(map[string]migration.MigrationInfo)

type Migrator struct {
	db *database.Database
}

func NewMigrator(db *database.Database) interfaces.Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Apply(ctx context.Context) error {
	if err := migration.Validate(migrations); err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	if err := m.createMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	migrations := migration.GetMigrations(migrations)
	if len(migrations) == 0 {
		return nil
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migration.CompareVersions(migrations[i].Version, migrations[j].Version) < 0
	})

	for _, migration := range migrations {
		applied, err := m.isApplied(ctx, migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check if migration %s is applied: %w", migration.Version, err)
		}

		if !applied {
			if err := m.applyMigration(ctx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
		}
	}

	return nil
}

func (m *Migrator) createMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at DATETIME NOT NULL,
			success BOOLEAN NOT NULL DEFAULT TRUE,
			error TEXT
		)
	`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

func (m *Migrator) isApplied(ctx context.Context, version string) (bool, error) {
	query := `SELECT COUNT(*) FROM schema_migrations WHERE version = ? AND success = TRUE`

	var count int
	err := m.db.QueryRowContext(ctx, query, version).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *Migrator) applyMigration(ctx context.Context, migra *interfaces.Migration) error {
	tx, err := m.db.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	sql := migration.GetMigrationSQL(migrations, migra.Version)
	if sql == "" {
		return fmt.Errorf("no SQL found for migration %s", migra.Version)
	}

	_, err = tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version, description, applied_at, success) VALUES (?, ?, ?, TRUE)`,
		migra.Version,
		migra.Description,
		timeutils.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
