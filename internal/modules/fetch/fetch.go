package subscription

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/core/system"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/parser"
	utilshttp "github.com/bestruirui/bestsub/internal/utils/http"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// Fetch 使用配置获取订阅内容
func Fetch(ctx context.Context, config *task.FetchConfig, subId int64) (*sub.FetchResult, error) {

	var lastErr error

	// 重试逻辑
	for attempt := 0; attempt <= config.Retries; attempt++ {
		if attempt > 0 {
			log.Warnf("Retrying fetch attempt %d/%d for %s", attempt, config.Retries, config.URL)

			// 等待重试延迟（递增延迟：1s, 2s, 3s...）
			retryDelay := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
			}
		}

		result, err := fetchOnce(ctx, config, subId)
		if err == nil {
			if attempt > 0 {
				log.Infof("Successfully fetched %s after %d retries", config.URL, attempt)
			}
			return result, nil
		}

		lastErr = err
		log.Warnf("Fetch attempt %d failed for %s: %v", attempt+1, config.URL, err)
	}

	return nil, fmt.Errorf("failed to fetch after %d attempts: %w", config.Retries+1, lastErr)
}

// fetchOnce 执行单次获取操作
func fetchOnce(ctx context.Context, config *task.FetchConfig, subId int64) (*sub.FetchResult, error) {
	startTime := time.Now()

	// 选择HTTP客户端
	var client *utilshttp.Client
	if config.ProxyEnable {
		client = utilshttp.Proxy()
		log.Debugf("Using proxy client to fetch: %s", config.URL)
	} else {
		client = utilshttp.Direct()
		log.Debugf("Using direct client to fetch: %s", config.URL)
	}
	defer client.Release()

	// 创建请求上下文
	ctx, cancel := context.WithTimeout(ctx, time.Duration(config.Timeout)*time.Second)
	defer cancel()

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", config.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "*/*")

	req.Close = true

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	// 限制响应大小
	const maxSize = 10 * 1024 * 1024 // 10MB
	limitedReader := io.LimitReader(resp.Body, maxSize)

	// 读取响应内容
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查是否超过大小限制
	if int64(len(content)) >= maxSize {
		return nil, fmt.Errorf("response size exceeds limit of %d bytes", maxSize)
	}
	subType, nodeCount, err := parser.Parse(&content, config.Type, subId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}
	duration := time.Since(startTime).Milliseconds()
	result := &sub.FetchResult{
		Type:      string(subType),
		NodeCount: nodeCount,
		Size:      int64(len(content)),
		Duration:  fmt.Sprintf("%dms", duration),
	}
	system.AddDownloadBytes(uint64(len(content)))

	log.Infof("Successfully fetched %d bytes from %s in %v", result.Size, config.URL, duration)

	return result, nil
}
