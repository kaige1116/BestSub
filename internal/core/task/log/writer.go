package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// Add 添加日志到内存缓存，超过限制时自动处理
func (w *LogWriter) Add(level string, progress int, message string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 创建日志条目
	entry := LogEntry{
		Time:     timeutils.Now().Format(time.RFC3339Nano),
		Level:    level,
		Progress: progress,
		Message:  message,
	}

	w.logs = append(w.logs, entry)

	select {
	case w.logChan <- entry:
	default:
		// 通道满了就跳过，避免阻塞
	}

	if len(w.logs) >= w.maxLogs {
		w.appendToLogFile()

		w.logs = w.logs[:0]
	}
}

// GetStream 获取实时日志流通道（用于WebSocket推送）
func (w *LogWriter) GetStream() <-chan LogEntry {
	return w.logChan
}

// Close 关闭写入器并一次性写入所有日志到文件
func (w *LogWriter) Close(finalStatus bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	close(w.logChan)

	if err := w.appendToLogFile(); err != nil {
		return err
	}

	if finalStatus {
		finalFilePath := strings.Replace(w.filePath, "_0.log", "_1.log", 1)

		if err := os.Rename(w.filePath, finalFilePath); err != nil {
			return fmt.Errorf("failed to rename log file: %w", err)
		}
	}

	return nil
}

// appendToLogFile 将内存中的日志追加写入到日志文件
func (w *LogWriter) appendToLogFile() error {
	if len(w.logs) == 0 {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(w.filePath), 0755); err != nil {
		return err
	}

	var builder strings.Builder
	for _, logEntry := range w.logs {
		logLine, _ := json.Marshal(logEntry)
		builder.Write(logLine)
		builder.WriteByte('\n')
	}

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(builder.String())
	return err
}
