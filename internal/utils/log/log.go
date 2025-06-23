package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/sirupsen/logrus"
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level  string `json:"level"`  // 日志等级: debug, info, warn, error
	Output string `json:"output"` // 输出方式: console, file, both
	File   string `json:"file"`   // 日志文件路径
}

// Logger 全局日志实例
var Logger *logrus.Logger

// StartupTime 程序启动时间
var StartupTime time.Time

// 自定义格式化器，支持debug模式显示文件位置
type CustomFormatter struct {
	TimestampFormat string
	ShowCaller      bool
}

// Format 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	var caller string
	if f.ShowCaller && entry.HasCaller() {
		// 获取调用者信息
		file := filepath.Base(entry.Caller.File)
		caller = fmt.Sprintf(" [%s:%d]", file, entry.Caller.Line)
	}

	// 格式: [时间] [等级] [文件:行号] 消息
	msg := fmt.Sprintf("[%s] [%s]%s %s\n", timestamp, level, caller, entry.Message)

	// 如果有字段，添加字段信息
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			msg += fmt.Sprintf("  %s=%v\n", key, value)
		}
	}

	return []byte(msg), nil
}

// Init 初始化日志系统
func Init(config config.Config) error {
	// 记录程序启动时间
	StartupTime = time.Now()

	Logger = logrus.New()

	// 设置日志等级
	level, err := logrus.ParseLevel(config.Log.Level)
	if err != nil {
		level = logrus.InfoLevel // 默认为info级别
	}
	Logger.SetLevel(level)

	// 启用调用者信息（用于显示文件和行号）
	Logger.SetReportCaller(true)

	// 设置自定义格式化器
	showCaller := level == logrus.DebugLevel // debug模式下显示文件位置
	Logger.SetFormatter(&CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ShowCaller:      showCaller,
	})

	// 设置输出
	switch strings.ToLower(config.Log.Output) {
	case "console":
		Logger.SetOutput(os.Stdout)
	case "file":
		if err := setupFileOutput(config.Log.Dir); err != nil {
			return fmt.Errorf("设置文件输出失败: %v", err)
		}
	case "both":
		if err := setupBothOutput(config.Log.Dir); err != nil {
			return fmt.Errorf("设置双重输出失败: %v", err)
		}
	default:
		Logger.SetOutput(os.Stdout) // 默认输出到控制台
	}

	// 记录程序启动信息
	if Logger != nil {
		Logger.Info("程序启动")
	}

	return nil
}

// setupFileOutput 设置文件输出
func setupFileOutput(logDir string) error {
	// 生成基于时间戳的日志文件名
	logFileName := generateTimestampLogFileName()
	logFilePath := filepath.Join(logDir, logFileName)

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	Logger.SetOutput(file)
	return nil
}

// setupBothOutput 设置同时输出到控制台和文件
func setupBothOutput(logDir string) error {
	// 生成基于时间戳的日志文件名
	logFileName := generateTimestampLogFileName()
	logFilePath := filepath.Join(logDir, logFileName)

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 创建多重写入器，同时写入控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, file)
	Logger.SetOutput(multiWriter)
	return nil
}

// generateTimestampLogFileName 生成基于时间戳的日志文件名
func generateTimestampLogFileName() string {
	// 使用启动时间生成文件名，格式：bestsub_20060102_150405.log
	timestamp := StartupTime.Format("20060102_150405")
	return fmt.Sprintf("bestsub_%s.log", timestamp)
}

// GetCurrentLogFileName 获取当前日志文件名（用于外部查询）
func GetCurrentLogFileName() string {
	return generateTimestampLogFileName()
}

// GetStartupTime 获取程序启动时间
func GetStartupTime() time.Time {
	return StartupTime
}

// 便捷的日志函数

// Debug 输出debug级别日志
func Debug(args ...interface{}) {
	if Logger != nil {
		Logger.Debug(args...)
	}
}

// Debugf 输出格式化的debug级别日志
func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, args...)
	}
}

// Info 输出info级别日志
func Info(args ...interface{}) {
	if Logger != nil {
		Logger.Info(args...)
	}
}

// Infof 输出格式化的info级别日志
func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, args...)
	}
}

// Warn 输出warn级别日志
func Warn(args ...interface{}) {
	if Logger != nil {
		Logger.Warn(args...)
	}
}

// Warnf 输出格式化的warn级别日志
func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, args...)
	}
}

// Error 输出error级别日志
func Error(args ...interface{}) {
	if Logger != nil {
		Logger.Error(args...)
	}
}

// Errorf 输出格式化的error级别日志
func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}

// Fatal 输出fatal级别日志并退出程序
func Fatal(args ...interface{}) {
	if Logger != nil {
		Logger.Fatal(args...)
	}
}

// Fatalf 输出格式化的fatal级别日志并退出程序
func Fatalf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Fatalf(format, args...)
	}
}

// WithFields 创建带字段的日志条目
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithFields(fields)
	}
	return nil
}

// WithField 创建带单个字段的日志条目
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger != nil {
		return Logger.WithField(key, value)
	}
	return nil
}

// SetLevel 动态设置日志等级
func SetLevel(level string) error {
	if Logger == nil {
		return fmt.Errorf("日志系统未初始化")
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("无效的日志等级: %s", level)
	}

	Logger.SetLevel(logLevel)

	// 如果是debug级别，启用调用者信息显示
	showCaller := logLevel == logrus.DebugLevel
	Logger.SetFormatter(&CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ShowCaller:      showCaller,
	})

	return nil
}

// GetLevel 获取当前日志等级
func GetLevel() string {
	if Logger != nil {
		return Logger.GetLevel().String()
	}
	return "unknown"
}

// IsDebugEnabled 检查是否启用了debug级别
func IsDebugEnabled() bool {
	if Logger != nil {
		return Logger.IsLevelEnabled(logrus.DebugLevel)
	}
	return false
}

// LogWithCaller 手动指定调用者信息的日志函数
func LogWithCaller(level logrus.Level, skip int, msg string) {
	if Logger == nil {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		Logger.Log(level, msg)
		return
	}

	entry := Logger.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	})

	entry.Log(level, msg)
}
