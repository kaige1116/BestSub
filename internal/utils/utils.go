package utils

import (
	"os"
	"path/filepath"
	"time"
)

// 统一的时间函数，用于数据库操作中的时间戳
func Now() time.Time {
	return time.Now().Local()
}

// 检查目录是否可写
func IsWritableDir(dir string) bool {
	// 尝试在目录中创建临时文件
	testFile := filepath.Join(dir, ".write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

// 检查字符串切片是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
