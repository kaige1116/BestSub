package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

var migrations = migration.NewMigration(3)

func Get() *[]migration.Info {
	return migrations.Get()
}
