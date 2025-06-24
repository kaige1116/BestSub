package sqlite

import (
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/database/sqlite/migration"
	"github.com/bestruirui/bestsub/internal/database/sqlite/repository"
)

// New 创建新的SQLite数据库连接
func New(databasePath string) (*database.Database, error) {
	return database.New(databasePath)
}

// NewRepo 创建新的SQLite仓库
func NewRepo(db *database.Database) interfaces.Repository {
	return repository.NewRepository(db)
}

// NewMigrator 创建新的SQLite迁移器
func NewMigrator(db *database.Database) interfaces.Migrator {
	return migration.NewMigrator(db)
}
