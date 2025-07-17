package sqlite

import (
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/client/sqlite/migration"
	migModel "github.com/bestruirui/bestsub/internal/database/migration"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

const MigrationTable = `
CREATE TABLE IF NOT EXISTS "migrations" (
	"date" INTEGER NOT NULL UNIQUE,
	"version" TEXT NOT NULL,
	"description" TEXT NOT NULL,
	"applied_at" DATETIME NOT NULL,
	PRIMARY KEY("date")
);
`

func (db *DB) Migrate() error {
	migrations := migration.Get()
	if len(*migrations) == 0 {
		return nil
	}
	// 检查 migrations 表是否存在
	exists, err := db.migrationsTableExists()
	if err != nil {
		return fmt.Errorf("failed to check migrations table existence: %w", err)
	}

	if !exists {
		if err := db.ensureMigrationsTable(); err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}

		for _, migration := range *migrations {
			if err := db.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Date, err)
			}
		}
	} else {
		appliedDates, err := db.getAppliedMigrations()
		if err != nil {
			return fmt.Errorf("failed to get applied migrations: %w", err)
		}

		for _, migration := range *migrations {
			if !appliedDates[migration.Date] {
				if err := db.applyMigration(migration); err != nil {
					return fmt.Errorf("failed to apply migration %d: %w", migration.Date, err)
				}
			}
		}
	}

	return nil
}

// migrationsTableExists 检查 migrations 表是否存在
func (db *DB) migrationsTableExists() (bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'`
	var name string
	err := db.db.QueryRow(query).Scan(&name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check migrations table: %w", err)
	}
	return true, nil
}

// ensureMigrationsTable 确保 migrations 表存在
func (db *DB) ensureMigrationsTable() error {
	_, err := db.db.Exec(MigrationTable)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// getAppliedMigrations 一次性获取所有已应用的迁移日期
func (db *DB) getAppliedMigrations() (map[int64]bool, error) {
	appliedDates := make(map[int64]bool)

	rows, err := db.db.Query("SELECT date FROM migrations")
	if err != nil {
		return appliedDates, err
	}
	defer rows.Close()

	for rows.Next() {
		var date int64
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("failed to scan migration date: %w", err)
		}
		appliedDates[date] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate migration rows: %w", err)
	}

	return appliedDates, nil
}

// applyMigration 应用迁移
func (db *DB) applyMigration(migration migModel.Info) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(migration.Content)
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	_, err = tx.Exec("INSERT INTO migrations (date, version, description, applied_at) VALUES (?, ?, ?, ?)",
		migration.Date, migration.Version, migration.Description, local.Time())
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
