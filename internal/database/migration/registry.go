package migration

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
)

// MigrationFunc 迁移函数类型
type MigrationFunc func() string

// MigrationInfo 迁移信息
type MigrationInfo struct {
	Description string
	Func        MigrationFunc
}

// Register 注册迁移到指定的 map
func Register(migrations map[string]MigrationInfo, version, description string, fn MigrationFunc) {
	migrations[version] = MigrationInfo{
		Description: description,
		Func:        fn,
	}
}

// GetMigrations 获取所有注册的迁移
func GetMigrations(migrations map[string]MigrationInfo) []*interfaces.Migration {
	var result []*interfaces.Migration

	for version, info := range migrations {
		result = append(result, &interfaces.Migration{
			Version:     version,
			Description: info.Description,
		})
	}

	// 按版本号排序
	sort.Slice(result, func(i, j int) bool {
		return CompareVersions(result[i].Version, result[j].Version) < 0
	})

	return result
}

// GetMigrationSQL 根据版本号获取迁移SQL
func GetMigrationSQL(migrations map[string]MigrationInfo, version string) string {
	if info, exists := migrations[version]; exists {
		return info.Func()
	}
	return ""
}

// HasMigration 检查是否存在指定版本的迁移
func HasMigration(migrations map[string]MigrationInfo, version string) bool {
	_, exists := migrations[version]
	return exists
}

// GetVersions 获取所有版本号
func GetVersions(migrations map[string]MigrationInfo) []string {
	var versions []string
	for version := range migrations {
		versions = append(versions, version)
	}

	// 按版本号排序
	sort.Slice(versions, func(i, j int) bool {
		return CompareVersions(versions[i], versions[j]) < 0
	})

	return versions
}

// Validate 验证注册的迁移
func Validate(migrations map[string]MigrationInfo) error {
	// 检查版本号格式
	for version := range migrations {
		if !IsValidVersion(version) {
			return fmt.Errorf("invalid version format: %s", version)
		}
	}

	// 检查每个迁移是否有对应的SQL
	for version, info := range migrations {
		sql := info.Func()
		if sql == "" {
			return fmt.Errorf("no SQL found for migration %s", version)
		}
	}

	return nil
}

// CompareVersions 比较版本号
func CompareVersions(v1, v2 string) int {
	n1, _ := strconv.Atoi(v1)
	n2, _ := strconv.Atoi(v2)

	if n1 < n2 {
		return -1
	} else if n1 > n2 {
		return 1
	}
	return 0
}

// IsValidVersion 检查版本号格式是否有效
func IsValidVersion(version string) bool {
	if len(version) != 3 {
		return false
	}
	_, err := strconv.Atoi(version)
	return err == nil
}
