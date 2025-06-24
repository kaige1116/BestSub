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
	"github.com/bestruirui/bestsub/internal/utils/color"
	"github.com/sirupsen/logrus"
)

var (
	Logger          *logrus.Logger
	StartupTime     time.Time
	colorStripRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	logFile         *os.File
	logFileMutex    sync.RWMutex
)

// 包装Writer，自动过滤颜色代码
type ColorStripWriter struct {
	writer io.Writer
}

// 创建一个新的颜色过滤Writer
func NewColorStripWriter(w io.Writer) *ColorStripWriter {
	return &ColorStripWriter{writer: w}
}

// 实现io.Writer接口，写入时自动过滤颜色代码
func (csw *ColorStripWriter) Write(p []byte) (n int, err error) {
	stripped := colorStripRegex.ReplaceAll(p, []byte(""))
	_, err = csw.writer.Write(stripped)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// 自定义格式化器
type CustomFormatter struct {
	TimestampFormat string
	ShowCaller      bool
	EnableColor     bool
}

// 根据日志等级获取对应颜色
func (f *CustomFormatter) getColorByLevel(level logrus.Level) string {
	if !f.EnableColor {
		return ""
	}

	colorMap := map[logrus.Level]string{
		logrus.DebugLevel: color.Cyan,
		logrus.InfoLevel:  color.Green,
		logrus.WarnLevel:  color.Yellow,
		logrus.ErrorLevel: color.Red,
		logrus.FatalLevel: color.Red + color.Bold,
		logrus.PanicLevel: color.Red + color.Bold,
	}

	return colorMap[level]
}

// 实现logrus.Formatter接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())[:1]

	var caller string
	if f.ShowCaller {
		// 优先使用我们手动设置的调用者信息
		if file, ok := entry.Data["file"]; ok {
			if line, ok := entry.Data["line"]; ok {
				caller = fmt.Sprintf(" [%s:%v]", file, line)
			}
		} else if entry.HasCaller() {
			// 回退到 logrus 自动获取的调用者信息
			file := filepath.Base(entry.Caller.File)
			caller = fmt.Sprintf(" [%s:%d]", file, entry.Caller.Line)
		}
	}

	logColor := f.getColorByLevel(entry.Level)
	reset := ""
	if f.EnableColor {
		reset = color.Reset
	}

	// 构建主消息
	var msg string
	if f.EnableColor {
		msg = fmt.Sprintf("[%s] %s[%s]%s%s %s%s%s\n",
			timestamp, logColor, level, reset, caller, logColor, entry.Message, reset)
	} else {
		msg = fmt.Sprintf("[%s] [%s]%s %s\n", timestamp, level, caller, entry.Message)
	}

	// 添加字段信息（排除我们用于调用者信息的字段）
	for key, value := range entry.Data {
		if key == "file" || key == "line" {
			continue // 跳过调用者信息字段
		}
		if f.EnableColor {
			msg += fmt.Sprintf("  %s%s=%v%s\n", logColor, key, value, reset)
		} else {
			msg += fmt.Sprintf("  %s=%v\n", key, value)
		}
	}

	return []byte(msg), nil
}

// 检查输出是否为终端
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return (f == os.Stdout || f == os.Stderr) && os.Getenv("TERM") != ""
	}
	return false
}

// 初始化日志系统
func Initialize(config config.LogConfig) error {
	StartupTime = time.Now()
	Logger = logrus.New()

	// 设置日志等级
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
	Logger.SetReportCaller(false) // 我们手动处理调用者信息

	showCaller := level == logrus.DebugLevel
	enableColor := false

	// 根据输出类型设置输出和颜色
	switch strings.ToLower(config.Output) {
	case "console":
		closeLogFile()
		Logger.SetOutput(os.Stdout)
		enableColor = isTerminal(os.Stdout)
	case "file":
		if err := setupFileOutput(config.Dir); err != nil {
			return fmt.Errorf("设置文件输出失败: %v", err)
		}
	case "both":
		if err := setupBothOutput(config.Dir); err != nil {
			return fmt.Errorf("设置双重输出失败: %v", err)
		}
		enableColor = isTerminal(os.Stdout)
	default:
		closeLogFile()
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

func openLogFile(logDir string) (*os.File, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}
	logFilePath := filepath.Join(logDir, fmt.Sprintf("bestsub_%s.log", StartupTime.Format("20060102_150405")))
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}
	logFileMutex.Lock()
	logFile = file
	logFileMutex.Unlock()

	return file, nil
}

func setupFileOutput(logDir string) error {
	closeLogFile()
	file, err := openLogFile(logDir)
	if err != nil {
		return err
	}
	Logger.SetOutput(file)
	return nil
}

func setupBothOutput(logDir string) error {
	closeLogFile()
	file, err := openLogFile(logDir)
	if err != nil {
		return err
	}

	fileWriter := NewColorStripWriter(file)
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
	Logger.SetOutput(multiWriter)
	return nil
}

func closeLogFile() error {
	logFileMutex.Lock()
	defer logFileMutex.Unlock()
	var err error
	if logFile != nil {
		err = logFile.Close()
		logFile = nil
	}
	return err
}

func Close() error {
	return closeLogFile()
}

// logWithCaller 带调用者信息的日志记录
func logWithCaller(level logrus.Level, skip int, args ...interface{}) {
	if Logger == nil {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		Logger.Log(level, args...)
		return
	}

	// 创建带调用者信息的日志条目
	entry := Logger.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	})

	entry.Log(level, args...)
}

// logfWithCaller 带调用者信息的格式化日志记录
func logfWithCaller(level logrus.Level, skip int, format string, args ...interface{}) {
	if Logger == nil {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		Logger.Logf(level, format, args...)
		return
	}

	// 创建带调用者信息的日志条目
	entry := Logger.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	})

	entry.Logf(level, format, args...)
}

func Debug(args ...interface{}) {
	logWithCaller(logrus.DebugLevel, 2, args...)
}

func Debugf(format string, args ...interface{}) {
	logfWithCaller(logrus.DebugLevel, 2, format, args...)
}

func Info(args ...interface{}) {
	logWithCaller(logrus.InfoLevel, 2, args...)
}

func Infof(format string, args ...interface{}) {
	logfWithCaller(logrus.InfoLevel, 2, format, args...)
}

func Warn(args ...interface{}) {
	logWithCaller(logrus.WarnLevel, 2, args...)
}

func Warnf(format string, args ...interface{}) {
	logfWithCaller(logrus.WarnLevel, 2, format, args...)
}

func Error(args ...interface{}) {
	logWithCaller(logrus.ErrorLevel, 2, args...)
}

func Errorf(format string, args ...interface{}) {
	logfWithCaller(logrus.ErrorLevel, 2, format, args...)
}

func Fatal(args ...interface{}) {
	logWithCaller(logrus.FatalLevel, 2, args...)
}

func Fatalf(format string, args ...interface{}) {
	logfWithCaller(logrus.FatalLevel, 2, format, args...)
}

// WithFields 创建带字段的日志条目（带调用者信息）
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if Logger == nil {
		return nil
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// 合并调用者信息和用户字段
		allFields := make(map[string]interface{})
		allFields["file"] = filepath.Base(file)
		allFields["line"] = line

		// 添加用户字段
		for k, v := range fields {
			allFields[k] = v
		}

		return Logger.WithFields(allFields)
	}

	return Logger.WithFields(fields)
}

// WithField 创建带单个字段的日志条目（带调用者信息）
func WithField(key string, value interface{}) *logrus.Entry {
	if Logger == nil {
		return nil
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return Logger.WithFields(logrus.Fields{
			"file": filepath.Base(file),
			"line": line,
			key:    value,
		})
	}

	return Logger.WithField(key, value)
}
