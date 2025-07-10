package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
)

// parseLogFileName 解析日志文件名
func parseLogFileName(filename string) (*SessionInfo, error) {
	// 移除.log后缀
	name := strings.TrimSuffix(filename, ".log")
	parts := strings.Split(name, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filename format: %s", filename)
	}

	// 解析时间部分
	timeStr := parts[0]
	execTime, err := time.Parse("20060102150405", timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format in filename: %s", timeStr)
	}

	// 解析状态部分
	statusStr := parts[1]
	status := statusStr == "1"

	return &SessionInfo{
		Time:   execTime,
		Status: status,
	}, nil
}

// generateLogFilePath 生成日志文件路径
func generateLogFilePath(taskID int64, time time.Time, status bool) string {
	timeStr := time.Format("20060102150405")
	statusStr := "0"
	if status {
		statusStr = "1"
	}
	filename := fmt.Sprintf("%s_%s.log", timeStr, statusStr)
	return filepath.Join(config.Get().Log.Dir, "tasks", strconv.FormatInt(taskID, 10), filename)
}

// readLogFileLines 读取日志文件指定行范围的内容，返回JSON数组
func readLogFileLines(filePath string, offset, limit int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentLine int
	var result strings.Builder
	var hasData bool

	result.WriteByte('[')

	for scanner.Scan() {
		if currentLine < offset {
			currentLine++
			continue
		}

		if currentLine >= offset+limit {
			break
		}

		if hasData {
			result.WriteByte(',')
		}

		result.WriteString(scanner.Text())
		hasData = true
		currentLine++
	}

	result.WriteByte(']')

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return result.String(), nil
}

// getFileSize 获取文件大小
func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
