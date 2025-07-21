package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bestruirui/bestsub/internal/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Name    string `json:"-"`
}

var (
	wsChannel chan LogEntry

	basePath string = "build"

	useConsole bool

	useFile bool

	logger *Logger
)
var consoleEncoder = zapcore.EncoderConfig{
	TimeKey:       "time",
	LevelKey:      "level",
	MessageKey:    "msg",
	CallerKey:     "caller",
	StacktraceKey: "stacktrace",
	EncodeLevel:   zapcore.CapitalColorLevelEncoder,
	EncodeTime:    zapcore.RFC3339TimeEncoder,
	EncodeCaller:  zapcore.ShortCallerEncoder,
}
var fileEncoder = zapcore.EncoderConfig{
	TimeKey:       "time",
	LevelKey:      "level",
	MessageKey:    "msg",
	CallerKey:     "caller",
	StacktraceKey: "stacktrace",
	EncodeLevel:   zapcore.LowercaseLevelEncoder,
	EncodeTime:    zapcore.RFC3339TimeEncoder,
	EncodeCaller:  zapcore.ShortCallerEncoder,
}

// Logger 日志记录器
type Logger struct {
	*zap.SugaredLogger
	file *os.File
}
type Config struct {
	Level      string
	Path       string
	UseConsole bool
	UseFile    bool
	Name       string
	CallerSkip int
}

// webSocketHook 发送日志到WebSocket通道
func webSocketHook(entry zapcore.Entry) error {
	if wsChannel == nil {
		return nil
	}

	logEntry := LogEntry{
		Level:   entry.Level.String(),
		Message: entry.Message,
		Name:    entry.LoggerName,
	}

	select {
	case wsChannel <- logEntry:
	default:
	}

	return nil
}

func init() {
	wsChannel = make(chan LogEntry, 1000)

	logger, _ = NewLogger(Config{
		Level:      "debug",
		UseConsole: true,
		CallerSkip: 1,
		UseFile:    false,
		Name:       "main",
	})
}

func Initialize(level, path, method string) error {
	logger.Close()

	basePath = path
	mainPath := filepath.Join(basePath, "main", time.Now().Format("20060102150405")+".log")

	switch method {
	case "console":
		useConsole = true
		useFile = false
	case "file":
		useConsole = false
		useFile = true
	case "both":
		useConsole = true
		useFile = true
	default:
		useConsole = true
		useFile = false
	}

	var err error
	logger, err = NewLogger(Config{
		Level:      level,
		Path:       mainPath,
		UseConsole: useConsole,
		UseFile:    useFile,
		Name:       "main",
		CallerSkip: 1,
	})
	if err != nil {
		return err
	}
	return nil
}
func GetDefaultLogger() *Logger {
	return logger
}
func NewTaskLogger(taskid uint16, level string) (*Logger, error) {
	taskidstr := strconv.FormatUint(uint64(taskid), 10)
	name := "task_" + taskidstr
	path := filepath.Join(basePath, "task", taskidstr, time.Now().Format("20060102150405")+".log")
	return NewLogger(Config{
		Level:      level,
		Path:       path,
		UseConsole: utils.IsDebug(),
		UseFile:    useFile,
		Name:       name,
		CallerSkip: 1,
	})
}

// GetWSChannel 获取全局WebSocket通道
func GetWSChannel() <-chan LogEntry {
	return wsChannel
}

// NewLogger 创建日志记录器
func NewLogger(config Config) (*Logger, error) {
	parsedLevel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		parsedLevel = zapcore.InfoLevel
	}

	var cores []zapcore.Core
	var file *os.File

	if config.UseConsole {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoder),
			zapcore.AddSync(os.Stdout),
			parsedLevel,
		)
		cores = append(cores, consoleCore)
	}

	if config.UseFile && config.Path != "" {
		file, err = createLogFile(config.Path)
		if err != nil {
			return nil, err
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoder),
			zapcore.AddSync(file),
			parsedLevel,
		)
		cores = append(cores, fileCore)
	}

	wsEncoderConfig := zapcore.EncoderConfig{
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}

	wsCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(wsEncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.DebugLevel,
	)
	cores = append(cores, wsCore)

	core := zapcore.NewTee(cores...)
	logger := zap.New(
		core,
		zap.Hooks(webSocketHook),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	logger.Named(config.Name)

	return &Logger{
		SugaredLogger: logger.Sugar(),
		file:          file,
	}, nil
}

// createLogFile 创建日志文件
func createLogFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return file, nil
}

func (l *Logger) Close() error {
	l.SugaredLogger.Sync()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

func Debug(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Debug(args...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Error(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Error(args...)
}
func Fatal(args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Fatalf(template, args...)
}
func Close() error {
	return logger.Close()
}
