package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/node"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
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

func createSuccessResult(count uint32, startTime time.Time) subModel.Result {
	return subModel.Result{
		Success:  1,
		Fail:     0,
		Msg:      "sub updated successfully",
		RawCount: count,
		LastRun:  time.Now(),
		Duration: uint16(time.Since(startTime).Milliseconds()),
	}
}

func Do(ctx context.Context, subID uint16, config string) subModel.Result {
	startTime := time.Now()

	var subConfig subModel.Config
	if err := json.Unmarshal([]byte(config), &subConfig); err != nil {
		log.Warnf("fetch task %d failed: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	log.Debugf("fetch task %d started", subID)

	client := mihomo.Default(false)
	if client == nil {
		log.Warnf("fetch task %d failed: proxy config error", subID)
		return createFailureResult("proxy config error", startTime)
	}
	defer client.Release()
	client.Timeout = time.Duration(subConfig.Timeout) * time.Second

	subUrl := genSubConverterUrl(subConfig.Url, subConfig.Proxy)

	req, err := http.NewRequestWithContext(ctx, "GET", subUrl, nil)
	if err != nil {
		log.Warnf("fetch task %d failed: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("fetch task %d failed: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("fetch task %d failed: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}

	nodes, err := parser.Parse(&content, subID)
	if err != nil {
		log.Warnf("fetch task %d failed: %v", subID, err)
		return createFailureResult(err.Error(), startTime)
	}
	count := len(*nodes)

	node.Add(nodes)

	log.Debugf("fetch task %d completed, node count: %d,  duration: %dms",
		subID, count, uint16(time.Since(startTime).Milliseconds()))

	return createSuccessResult(uint32(count), startTime)
}

func genSubConverterUrl(subUrl string, enableProxy bool) string {
	subUrl = url.QueryEscape(subUrl)
	cfg := config.Base()
	scHost := cfg.SubConverter.Host
	scPort := cfg.SubConverter.Port
	if enableProxy {
		proxy := op.GetSettingStr(setting.PROXY_URL)
		proxy = url.QueryEscape(proxy)
		return fmt.Sprintf("http://%s:%d/sub?target=clash&list=true&url=%s&proxy=%s", scHost, scPort, subUrl, proxy)
	}
	return fmt.Sprintf("http://%s:%d/sub?target=clash&list=true&url=%s", scHost, scPort, subUrl)
}
