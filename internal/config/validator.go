package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/bestruirui/bestsub/internal/utils"
)

// ValidateConfig 验证完整配置
func ValidateConfig(config *Config) error {
	if err := ValidateServerConfig(&config.Server); err != nil {
		return fmt.Errorf("服务器配置验证失败: %v", err)
	}

	if err := ValidateDatabaseConfig(&config.Database); err != nil {
		return fmt.Errorf("数据库配置验证失败: %v", err)
	}

	if err := ValidateLogConfig(&config.Log); err != nil {
		return fmt.Errorf("日志配置验证失败: %v", err)
	}

	if err := ValidateJWTConfig(&config.JWT); err != nil {
		return fmt.Errorf("JWT配置验证失败: %v", err)
	}

	return nil
}

// ValidateServerConfig 验证服务器配置
func ValidateServerConfig(config *ServerConfig) error {
	// 验证端口
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("端口号必须在1-65535范围内，当前值: %d", config.Port)
	}

	// 验证主机地址
	if config.Host == "" {
		return fmt.Errorf("主机地址不能为空")
	}

	// 验证主机地址格式
	if config.Host != "0.0.0.0" && config.Host != "localhost" {
		if ip := net.ParseIP(config.Host); ip == nil {
			return fmt.Errorf("无效的主机地址格式: %s", config.Host)
		}
	}

	return nil
}

// 验证数据库配置
func ValidateDatabaseConfig(config *DatabaseConfig) error {
	if config.Type == "" {
		return fmt.Errorf("数据库类型不能为空")
	}
	if config.Path == "" {
		return fmt.Errorf("数据库路径不能为空")
	}
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("无法创建数据库目录 %s: %v", dir, err)
	}
	if !utils.IsWritableDir(dir) {
		return fmt.Errorf("数据库目录 %s 不可写", dir)
	}
	return nil
}

// 验证日志配置
func ValidateLogConfig(config *LogConfig) error {
	// 验证日志等级
	validLevels := []string{"debug", "info", "warn", "error"}
	if !utils.Contains(validLevels, strings.ToLower(config.Level)) {
		return fmt.Errorf("无效的日志等级: %s，支持的等级: %v", config.Level, validLevels)
	}

	// 验证输出方式
	validOutputs := []string{"console", "file", "both"}
	if !utils.Contains(validOutputs, strings.ToLower(config.Output)) {
		return fmt.Errorf("无效的日志输出方式: %s，支持的方式: %v", config.Output, validOutputs)
	}

	// 如果输出到文件，验证文件路径
	if config.Output == "file" || config.Output == "both" {
		if config.Dir == "" {
			return fmt.Errorf("日志输出到文件时，文件路径不能为空")
		}

		// 检查日志目录是否可写
		dir := config.Dir
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("无法创建日志目录 %s: %v", dir, err)
		}

		if !utils.IsWritableDir(dir) {
			return fmt.Errorf("日志目录 %s 不可写", dir)
		}
	}

	return nil
}

// ValidateJWTConfig 验证JWT配置
func ValidateJWTConfig(config *JWTConfig) error {
	if config.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	if len(config.Secret) < 16 {
		return fmt.Errorf("JWT密钥长度不能少于16个字符，当前长度: %d", len(config.Secret))
	}

	// 检查是否使用默认密钥
	if strings.Contains(config.Secret, "change-me") || config.Secret == "bestsub-jwt-secret" {
		return fmt.Errorf("请修改默认的JWT密钥以确保安全性")
	}

	return nil
}
