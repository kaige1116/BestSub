package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/sirupsen/logrus"
)

// 颜色常量定义
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorPurple = "\033[35m"
	ColorBold   = "\033[1m"
)

// colorStripRegex 用于匹配ANSI颜色代码的正则表达式
var colorStripRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// ColorStripWriter 包装Writer，自动过滤颜色代码
type ColorStripWriter struct {
	writer io.Writer
}

// NewColorStripWriter 创建一个新的颜色过滤Writer
func NewColorStripWriter(w io.Writer) *ColorStripWriter {
	return &ColorStripWriter{writer: w}
}

// Write 实现io.Writer接口，写入时自动过滤颜色代码
func (csw *ColorStripWriter) Write(p []byte) (n int, err error) {
	// 过滤掉ANSI颜色代码
	stripped := colorStripRegex.ReplaceAll(p, []byte(""))

	// 写入过滤后的数据
	_, err = csw.writer.Write(stripped)
	if err != nil {
		return 0, err
	}

	// 返回原始数据的长度，表示我们"处理"了所有输入数据
	return len(p), nil
}

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

// logFile 当前打开的日志文件
var logFile *os.File

// logFileMutex 保护日志文件操作的互斥锁
var logFileMutex sync.RWMutex

// 自定义格式化器，支持debug模式显示文件位置和颜色
type CustomFormatter struct {
	TimestampFormat string
	ShowCaller      bool
	EnableColor     bool // 是否启用颜色
}

// getColorByLevel 根据日志等级获取对应颜色
func (f *CustomFormatter) getColorByLevel(level logrus.Level) string {
	if !f.EnableColor {
		return ""
	}

	switch level {
	case logrus.DebugLevel:
		return ColorCyan
	case logrus.InfoLevel:
		return ColorGreen
	case logrus.WarnLevel:
		return ColorYellow
	case logrus.ErrorLevel:
		return ColorRed
	case logrus.FatalLevel, logrus.PanicLevel:
		return ColorRed + ColorBold
	default:
		return ""
	}
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

	// 获取颜色代码
	color := f.getColorByLevel(entry.Level)
	reset := ""
	if f.EnableColor {
		reset = ColorReset
	}

	// 格式: [时间] [等级] [文件:行号] 消息（带颜色）
	var msg string
	if f.EnableColor {
		msg = fmt.Sprintf("[%s] %s[%s]%s%s %s%s%s\n",
			timestamp, color, level, reset, caller, color, entry.Message, reset)
	} else {
		msg = fmt.Sprintf("[%s] [%s]%s %s\n", timestamp, level, caller, entry.Message)
	}

	// 如果有字段，添加字段信息
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			if f.EnableColor {
				msg += fmt.Sprintf("  %s%s=%v%s\n", color, key, value, reset)
			} else {
				msg += fmt.Sprintf("  %s=%v\n", key, value)
			}
		}
	}

	return []byte(msg), nil
}

// isTerminal 检查输出是否为终端（简单实现）
func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		// 检查是否为标准输出/错误且是终端
		return (v == os.Stdout || v == os.Stderr) && isatty(v)
	default:
		return false
	}
}

// isatty 检查文件描述符是否为终端
func isatty(f *os.File) bool {
	// 简单的终端检测，检查是否有TERM环境变量
	return os.Getenv("TERM") != ""
}

// 初始化日志系统
func Initialize(config config.LogConfig) error {
	StartupTime = time.Now()

	Logger = logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	Logger.SetReportCaller(true)

	showCaller := level == logrus.DebugLevel

	// 根据输出类型决定是否启用颜色
	enableColor := false
	switch strings.ToLower(config.Output) {
	case "console":
		closeLogFile() // 切换到控制台时关闭之前的文件
		Logger.SetOutput(os.Stdout)
		enableColor = isTerminal(os.Stdout)
	case "file":
		if err := setupFileOutput(config.Dir); err != nil {
			return fmt.Errorf("设置文件输出失败: %v", err)
		}
		enableColor = false // 文件输出不使用颜色
	case "both":
		if err := setupBothOutput(config.Dir); err != nil {
			return fmt.Errorf("设置双重输出失败: %v", err)
		}
		enableColor = isTerminal(os.Stdout) // 控制台启用颜色，文件自动过滤颜色
	default:
		closeLogFile() // 默认情况下也关闭之前的文件
		Logger.SetOutput(os.Stdout)
		enableColor = isTerminal(os.Stdout)
	}

	Logger.SetFormatter(&CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ShowCaller:      showCaller,
		EnableColor:     enableColor,
	})

	return nil
}

// setupFileOutput 设置文件输出
func setupFileOutput(logDir string) error {
	// 关闭之前打开的文件
	closeLogFile()

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

	// 保存文件引用
	logFileMutex.Lock()
	logFile = file
	logFileMutex.Unlock()

	Logger.SetOutput(file)
	return nil
}

// setupBothOutput 设置同时输出到控制台和文件
func setupBothOutput(logDir string) error {
	// 关闭之前打开的文件
	closeLogFile()

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

	// 保存文件引用
	logFileMutex.Lock()
	logFile = file
	logFileMutex.Unlock()

	// 创建颜色过滤器包装文件输出，确保文件中不包含颜色代码
	fileWriter := NewColorStripWriter(file)

	// 创建多重写入器，控制台保留颜色，文件过滤颜色
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
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

// closeLogFile 关闭当前打开的日志文件
func closeLogFile() {
	logFileMutex.Lock()
	defer logFileMutex.Unlock()

	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// Close 关闭日志系统，释放资源
func Close() {
	closeLogFile()
}

// GetLogFile 获取当前日志文件路径（只读）
func GetLogFile() *os.File {
	logFileMutex.RLock()
	defer logFileMutex.RUnlock()
	return logFile
}

// StripColors 移除字符串中的ANSI颜色代码（工具函数）
func StripColors(s string) string {
	return colorStripRegex.ReplaceAllString(s, "")
}

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

	// 检查当前格式化器的颜色设置
	currentFormatter, ok := Logger.Formatter.(*CustomFormatter)
	enableColor := false
	if ok {
		enableColor = currentFormatter.EnableColor
	}

	Logger.SetFormatter(&CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ShowCaller:      showCaller,
		EnableColor:     enableColor,
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

// SetColorEnabled 动态设置是否启用颜色
func SetColorEnabled(enabled bool) error {
	if Logger == nil {
		return fmt.Errorf("日志系统未初始化")
	}

	currentFormatter, ok := Logger.Formatter.(*CustomFormatter)
	if !ok {
		return fmt.Errorf("当前使用的不是自定义格式化器")
	}

	Logger.SetFormatter(&CustomFormatter{
		TimestampFormat: currentFormatter.TimestampFormat,
		ShowCaller:      currentFormatter.ShowCaller,
		EnableColor:     enabled,
	})

	return nil
}

// IsColorEnabled 检查是否启用了颜色
func IsColorEnabled() bool {
	if Logger != nil {
		if formatter, ok := Logger.Formatter.(*CustomFormatter); ok {
			return formatter.EnableColor
		}
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
