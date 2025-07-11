package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/local"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Name    string    `json:"-"`
}

var (
	wsChannel chan LogEntry

	encoderConfig zapcore.EncoderConfig

	basePath string

	useConsole bool

	useFile bool

	logger *Logger
)

// Logger 日志记录器
type Logger struct {
	*zap.SugaredLogger
	file *os.File
}
type config struct {
	level      string
	path       string
	useConsole bool
	useFile    bool
	name       string
}

// webSocketHook 发送日志到WebSocket通道
func webSocketHook(entry zapcore.Entry) error {
	if wsChannel == nil {
		return nil
	}

	logEntry := LogEntry{
		Time:    entry.Time,
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

	encoderConfig = zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:  zapcore.RFC3339TimeEncoder,
	}
}

func Initialize(level, path, method string) error {
	basePath = path
	mainPath := filepath.Join(basePath, "main", local.Time().Format("20060102150405")+".log")
	switch method {
	case "console":
		useConsole = true
		useFile = false
	case "file":
		useConsole = false
		useFile = true
	case "all":
		useConsole = true
		useFile = true
	default:
		useConsole = true
		useFile = false
	}
	var err error
	logger, err = NewLogger(config{
		level:      level,
		path:       mainPath,
		useConsole: useConsole,
		useFile:    useFile,
		name:       "main",
	})
	if err != nil {
		return err
	}
	return nil
}

func NewTaskLogger(taskid int64, level string) (*Logger, error) {
	taskidstr := strconv.FormatInt(taskid, 10)
	name := "task_" + taskidstr
	path := filepath.Join(basePath, "task", taskidstr, local.Time().Format("20060102150405")+".log")
	return NewLogger(config{
		level:      level,
		path:       path,
		useConsole: useConsole,
		useFile:    useFile,
		name:       name,
	})
}

// GetWSChannel 获取全局WebSocket通道
func GetWSChannel() <-chan LogEntry {
	return wsChannel
}

// NewLogger 创建日志记录器
func NewLogger(config config) (*Logger, error) {
	parsedLevel, err := zapcore.ParseLevel(config.level)
	if err != nil {
		parsedLevel = zapcore.InfoLevel
	}

	writers, file, err := setupWriters(config.path)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writers...),
		parsedLevel,
	)

	logger := zap.New(core, zap.Hooks(webSocketHook))
	logger = logger.Named(config.name)

	return &Logger{
		SugaredLogger: logger.Sugar(),
		file:          file,
	}, nil
}

// setupWriters 设置输出writers
func setupWriters(path string) ([]zapcore.WriteSyncer, *os.File, error) {
	var writers []zapcore.WriteSyncer
	var file *os.File

	if useConsole {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	if useFile && path != "" {
		var err error
		file, err = createLogFile(path)
		if err != nil {
			return nil, nil, err
		}
		writers = append(writers, zapcore.AddSync(file))
	}
	if len(writers) == 0 {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	return writers, file, nil
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
	logger.Debug(args...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
func Close() error {
	return logger.Close()
}
