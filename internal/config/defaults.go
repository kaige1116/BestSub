package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DefaultValues 默认配置值常量
const (
	// 服务器默认配置
	DefaultServerPort = 8080
	DefaultServerHost = "0.0.0.0"

	// 数据库默认配置
	DefaultDatabasePath = "./bestsub/data/bestsub.db"

	// 日志默认配置
	DefaultLogLevel  = "info"
	DefaultLogOutput = "both"
	DefaultLogDir    = "./bestsub/log"

	// JWT默认配置
	DefaultJWTExpiresIn = 3600 // 1小时
)

// 创建默认配置文件
func CreateDefaultConfig(filePath string) error {
	// 确保配置目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 创建默认配置
	defaultConfig := getDefaultConfig()

	// 序列化为JSON
	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %v", err)
	}

	// 写入文件
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
			Path: DefaultDatabasePath,
		},
		Log: LogConfig{
			Level:  DefaultLogLevel,
			Output: DefaultLogOutput,
			Dir:    DefaultLogDir,
		},
		JWT: JWTConfig{
			Secret:    generateSecureJWTSecret(),
			ExpiresIn: DefaultJWTExpiresIn,
		},
	}
}

// generateSecureJWTSecret 生成安全的JWT密钥
func generateSecureJWTSecret() string {
	// 生成32字节的随机密钥
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机生成失败，使用固定的安全密钥
		return "bestsub-secure-jwt-secret-key-32-chars"
	}
	return hex.EncodeToString(bytes)
}
