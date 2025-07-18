package migration

import "sort"

type Info struct {
	Date        int64 // 迁移日期 格式 202507171200
	Version     string
	Description string
	Content     func() string
}

var clientMigrations = make(map[string][]*Info)

func Register(client string, date int64, version, description string, contentFunc func() string) {
	info := &Info{
		Date:        date,
		Version:     version,
		Description: description,
		Content:     contentFunc,
	}

	migrations := clientMigrations[client]

	index := sort.Search(len(migrations), func(i int) bool {
		return migrations[i].Date > date
	})

	migrations = append(migrations, nil)
	copy(migrations[index+1:], migrations[index:])
	migrations[index] = info

	clientMigrations[client] = migrations
}

// Get 获取指定客户端的迁移数据
func Get(client string) []*Info {
	if migrations := clientMigrations[client]; migrations != nil {
		return migrations
	}
	return make([]*Info, 0)
}
