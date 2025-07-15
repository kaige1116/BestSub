package utils

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
func RemoveAllControlCharacters(data *[]byte) {
	var cleanedData []byte
	original := *data
	for len(original) > 0 {
		r, size := utf8.DecodeRune(original)
		if r != utf8.RuneError && (r >= 32 && r <= 126) || r == '\n' || r == '\t' || r == '\r' || unicode.Is(unicode.Han, r) {
			cleanedData = append(cleanedData, original[:size]...)
		}
		original = original[size:]
	}
	*data = cleanedData
}
func IsDebug() bool {
	debug := os.Getenv("DEBUG")
	return strings.ToLower(debug) == "true"
}
