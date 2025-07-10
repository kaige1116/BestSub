package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// CreateWriter 创建日志写入器
func CreateWriter(taskID int64) *LogWriter {
	filePath := generateLogFilePath(taskID, timeutils.Now(), false)

	fmt.Println("filePath", filePath)

	return &LogWriter{
		taskID:   taskID,
		logs:     make([]LogEntry, 0, 1000),
		filePath: filePath,
		maxLogs:  1000,
		logChan:  make(chan LogEntry, 100), // 缓冲通道
	}
}

// GetTaskSessions 获取任务的执行会话列表
func GetTaskSessions(taskID int64, offset, limit int) (*[]SessionInfo, int64, error) {
	logDir := filepath.Join(config.Get().Log.Dir, "tasks", strconv.FormatInt(taskID, 10))

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return &[]SessionInfo{}, 0, nil
	}

	// 读取目录中的所有文件
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read log directory: %w", err)
	}

	var sessions []SessionInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		// 解析文件名获取会话信息
		sessionInfo, err := parseLogFileName(file.Name())
		if err != nil {
			continue // 跳过无效格式的文件
		}

		// 设置任务ID
		sessionInfo.TaskID = taskID

		// 获取文件大小
		filePath := filepath.Join(logDir, file.Name())
		if fileSize, err := getFileSize(filePath); err == nil {
			sessionInfo.FileSize = fileSize
		}

		sessions = append(sessions, *sessionInfo)
	}

	// 按时间倒序排序
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Time.After(sessions[j].Time)
	})

	// 应用分页
	total := int64(len(sessions))
	start := offset
	end := offset + limit

	if start >= len(sessions) {
		return &[]SessionInfo{}, total, nil
	}

	if end > len(sessions) {
		end = len(sessions)
	}

	result := sessions[start:end]
	return &result, total, nil
}

// GetSession 获取指定会话的日志内容
func GetSession(taskID int64, sessionTime time.Time, status bool, offset, limit int) (string, error) {
	filePath := generateLogFilePath(taskID, sessionTime, status)

	if _, err := os.Stat(filePath); err == nil {
		// 找到文件，读取内容
		return readLogFileLines(filePath, offset, limit)
	}

	return "", fmt.Errorf("log file not found for task %d at time %s", taskID, sessionTime)
}

// Cleanup 清理过期日志
func Cleanup(beforeTime time.Time) error {
	logsDir := filepath.Join(config.Get().Log.Dir, "tasks")

	// 检查目录是否存在
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return nil
	}

	// 遍历所有任务目录
	taskDirs, err := os.ReadDir(logsDir)
	if err != nil {
		return fmt.Errorf("failed to read logs directory: %w", err)
	}

	var deletedCount int
	for _, taskDir := range taskDirs {
		if !taskDir.IsDir() {
			continue
		}

		taskDirPath := filepath.Join(logsDir, taskDir.Name())

		// 读取任务目录中的日志文件
		logFiles, err := os.ReadDir(taskDirPath)
		if err != nil {
			continue
		}

		for _, logFile := range logFiles {
			if logFile.IsDir() || !strings.HasSuffix(logFile.Name(), ".log") {
				continue
			}

			// 解析文件名获取时间
			sessionInfo, err := parseLogFileName(logFile.Name())
			if err != nil {
				continue
			}

			// 检查是否需要删除
			if sessionInfo.Time.Before(beforeTime) {
				filePath := filepath.Join(taskDirPath, logFile.Name())
				if err := os.Remove(filePath); err != nil {
					// 记录错误但继续处理其他文件
					continue
				}
				deletedCount++
			}
		}
	}

	return nil
}
