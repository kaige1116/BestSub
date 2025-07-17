package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

const ClientName = "sqlite"

func Get() []*migration.Info {
	return migration.Get(ClientName)
}
