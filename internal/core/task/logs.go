package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

// WriteTaskLog 写入任务日志到文件
func WriteTaskLog(taskLog TaskLog, isSuccess bool) error {
	// 生成日志文件路径
	filePath := GetLogFilePath(taskLog.TaskID, taskLog.ExecutionID, isSuccess)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 将日志转换为JSON
	logLine, err := json.Marshal(taskLog)
	if err != nil {
		return fmt.Errorf("failed to marshal task log: %w", err)
	}

	// 追加写入文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// 写入日志行
	if _, err := file.WriteString(string(logLine) + "\n"); err != nil {
		return fmt.Errorf("failed to write log line: %w", err)
	}

	return nil
}

// GetLogFilePath 生成日志文件路径
// 格式：{timestamp}_{status}_{task_id}.log
func GetLogFilePath(taskID int64, executionID string, isSuccess bool) string {
	status := "0" // 失败
	if isSuccess {
		status = "1" // 成功
	}
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s_%d.log", timestamp, status, taskID)
	return filepath.Join("logs", "tasks", filename)
}

// ParseLogFileName 解析日志文件名
func ParseLogFileName(filename string) (*LogFileInfo, error) {
	// 解析格式：20240115120000_1_123.log
	name := strings.TrimSuffix(filename, ".log")
	parts := strings.Split(name, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid filename format: %s", filename)
	}

	// 解析时间戳
	timestamp, err := time.Parse("20060102150405", parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp in filename: %s", parts[0])
	}

	// 解析成功状态
	isSuccess := parts[1] == "1"

	// 解析任务ID
	taskID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID in filename: %s", parts[2])
	}

	return &LogFileInfo{
		Time:      timestamp,
		IsSuccess: isSuccess,
		TaskID:    taskID,
	}, nil
}

// ListLogFiles 查询日志文件列表
func ListLogFiles(taskID int64, isSuccess *bool) ([]string, error) {
	logDir := filepath.Join("logs", "tasks")

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// 读取目录中的所有文件
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	var matchingFiles []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		// 解析文件名
		fileInfo, err := ParseLogFileName(file.Name())
		if err != nil {
			log.Warnf("Failed to parse log filename %s: %v", file.Name(), err)
			continue
		}

		// 检查任务ID是否匹配
		if fileInfo.TaskID != taskID {
			continue
		}

		// 检查成功状态是否匹配（如果指定了过滤条件）
		if isSuccess != nil && fileInfo.IsSuccess != *isSuccess {
			continue
		}

		matchingFiles = append(matchingFiles, filepath.Join(logDir, file.Name()))
	}

	return matchingFiles, nil
}

// ReadTaskLogSummaries 读取任务日志摘要
func ReadTaskLogSummaries(taskID int64, offset, limit int) (*[]LogFileInfo, int64, error) {
	// 获取所有匹配的日志文件
	files, err := ListLogFiles(taskID, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list log files: %w", err)
	}

	// 解析文件名获取摘要信息
	var summaries []LogFileInfo
	for _, filePath := range files {
		filename := filepath.Base(filePath)
		fileInfo, err := ParseLogFileName(filename)
		if err != nil {
			log.Warnf("Failed to parse log filename %s: %v", filename, err)
			continue
		}
		summaries = append(summaries, *fileInfo)
	}

	// 按时间倒序排序
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Time.After(summaries[j].Time)
	})

	// 应用分页
	total := int64(len(summaries))
	start := offset
	end := offset + limit

	if start >= len(summaries) {
		return &[]LogFileInfo{}, total, nil
	}

	if end > len(summaries) {
		end = len(summaries)
	}

	result := summaries[start:end]
	return &result, total, nil
}

// ReadTaskLogDetail 读取任务详细日志
func ReadTaskLogDetail(taskID int64, logTime time.Time) (*[]TaskLog, error) {
	// 根据时间戳查找对应的日志文件
	timestamp := logTime.Format("20060102150405")

	// 尝试查找成功和失败的日志文件
	patterns := []string{
		filepath.Join("logs", "tasks", fmt.Sprintf("%s_1_%d.log", timestamp, taskID)),
		filepath.Join("logs", "tasks", fmt.Sprintf("%s_0_%d.log", timestamp, taskID)),
	}

	for _, filePath := range patterns {
		if _, err := os.Stat(filePath); err == nil {
			return readLogFile(filePath)
		}
	}

	return &[]TaskLog{}, nil
}

// readLogFile 读取日志文件内容
func readLogFile(filePath string) (*[]TaskLog, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var logs []TaskLog

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var taskLog TaskLog
		if err := json.Unmarshal([]byte(line), &taskLog); err != nil {
			log.Warnf("Failed to unmarshal log line: %v", err)
			continue
		}

		logs = append(logs, taskLog)
	}

	return &logs, nil
}

// CleanupOldLogFiles 清理过期日志文件
func CleanupOldLogFiles(beforeTime time.Time) error {
	logDir := filepath.Join("logs", "tasks")

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return nil
	}

	// 读取目录中的所有文件
	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var deletedCount int
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		// 解析文件名获取时间戳
		fileInfo, err := ParseLogFileName(file.Name())
		if err != nil {
			log.Warnf("Failed to parse log filename %s: %v", file.Name(), err)
			continue
		}

		// 检查是否需要删除
		if fileInfo.Time.Before(beforeTime) {
			filePath := filepath.Join(logDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Errorf("Failed to remove old log file %s: %v", filePath, err)
			} else {
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		log.Infof("清理了 %d 个过期日志文件", deletedCount)
	}

	return nil
}
