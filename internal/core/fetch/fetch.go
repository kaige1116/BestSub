package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/nodepool"
	parserModel "github.com/bestruirui/bestsub/internal/models/parser"
	subModel "github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/modules/parser"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func createFailureResult(msg string, startTime time.Time) subModel.Result {
	return subModel.Result{
		Success:  0,
		Fail:     1,
		Msg:      msg,
		LastRun:  time.Now(),
		Duration: uint16(time.Since(startTime).Milliseconds()),
	}
}

func createSuccessResult(rawCount, count uint32, startTime time.Time) subModel.Result {
	return subModel.Result{
		Success:  1,
		Fail:     0,
		Msg:      "fetch 任务执行完成",
		RawCount: rawCount,
		Count:    count,
		LastRun:  time.Now(),
		Duration: uint16(time.Since(startTime).Milliseconds()),
	}
}

func Do(ctx context.Context, subID uint16, config string) subModel.Result {
	startTime := time.Now()

	var subConfig subModel.Config
	if err := json.Unmarshal([]byte(config), &subConfig); err != nil {
		log.Warnf("fetch 任务执行失败 %d: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	log.Infof("fetch 任务执行中 %d", subID)

	client := mihomo.Default(subConfig.Proxy)
	if client == nil {
		log.Warnf("fetch 任务执行失败 %d: 代理配置错误", subID)
		return createFailureResult("代理配置错误", startTime)
	}
	defer client.Release()
	client.Timeout = time.Duration(subConfig.Timeout) * time.Second

	req, err := http.NewRequestWithContext(ctx, "GET", subConfig.Url, nil)
	if err != nil {
		log.Warnf("fetch 任务执行失败 %d: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("fetch 任务执行失败 %d: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("fetch 任务执行失败 %d: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	nodes, err := parser.Parse(&content, parserModel.ParserTypeMihomo, subID)
	if err != nil {
		log.Warnf("fetch 任务执行失败 %d: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	rawCount := uint32(len(*nodes))
	subStorage := nodepool.GetPoolBySubID(subID, len(*nodes))
	count := subStorage.AddNode(nodes)

	log.Infof("fetch 任务执行完成 %d, 原始节点数: %d, 有效节点数: %d, 耗时: %dms",
		subID, rawCount, count, uint16(time.Since(startTime).Milliseconds()))

	return createSuccessResult(rawCount, count, startTime)
}
