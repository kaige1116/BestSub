package migration

import (
	"sort"
)

type Info struct {
	Date        int64 // 迁移日期 格式 202507171200
	Version     string
	Description string
	Content     string
}
type Migrations struct {
	migrations []Info
}

func NewMigration(size int) *Migrations {
	return &Migrations{
		migrations: make([]Info, 0, size),
	}
}

func (m *Migrations) Register(date int64, version, description, content string) {
	newInfo := Info{
		Date:        date,
		Version:     version,
		Description: description,
		Content:     content,
	}

	index := sort.Search(len(m.migrations), func(i int) bool {
		return m.migrations[i].Date > date
	})

	m.migrations = append(m.migrations, Info{})
	copy(m.migrations[index+1:], m.migrations[index:])
	m.migrations[index] = newInfo
}

func (m *Migrations) Get() *[]Info {
	return &m.migrations
}
