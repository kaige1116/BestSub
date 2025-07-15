package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// 默认配置值常量
const (
	DefaultServerPort = 8080
	DefaultServerHost = "0.0.0.0"

	DefaultDatabasePath = "./data/data/bestsub.db"
	DefaultDatabaseType = "sqlite"

	DefaultLogLevel  = "info"
	DefaultLogOutput = "both"
	DefaultLogDir    = "./data/log"
)

// 创建默认配置文件
func CreateDefaultConfig(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	defaultConfig := getDefaultConfig()

	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// 获取默认配置
func getDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port: DefaultServerPort,
			Host: DefaultServerHost,
		},
		Database: DatabaseConfig{
			Type: DefaultDatabaseType,
			Path: DefaultDatabasePath,
		},
		Log: LogConfig{
			Level:  DefaultLogLevel,
			Output: DefaultLogOutput,
			Dir:    DefaultLogDir,
		},
		JWT: JWTConfig{
			Secret: generateSecureJWTSecret(),
		},
	}
}

// 生成安全的JWT密钥
func generateSecureJWTSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "bestsub-secure-jwt-secret-key-32-chars"
	}
	return hex.EncodeToString(bytes)
}
