package checker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/node"
	"github.com/bestruirui/bestsub/internal/core/task"
	checkModel "github.com/bestruirui/bestsub/internal/models/check"
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

type Alive struct {
	URL         string `json:"url" name:"测试链接" value:"https://www.gstatic.com/generate_204"`
	ExptectCode int    `json:"exptect_code" name:"期望状态码" value:"204"`
	Thread      int    `json:"thread" name:"线程数" value:"100"`
	Timeout     int    `json:"timeout" name:"超时时间" value:"10" desc:"单个节点检测的超时时间(s)"`
}
type Result struct {
	AliveCount uint16 `json:"alive_count" desc:"存活节点数量"`
	DeadCount  uint16 `json:"dead_count" desc:"死亡节点数量"`
	Delay      uint16 `json:"delay" desc:"平均延迟"`
}

func (e *Alive) Init() error {
	return nil
}

func (e *Alive) Run(ctx context.Context, log *log.Logger, subID []uint16) checkModel.Result {
	startTime := time.Now()
	var nodes []nodeModel.Data
	var aliveCount, deadCount, totalDelay int64
	if len(subID) == 0 {
		nodes = node.GetAll()
	} else {
		nodes = *node.GetBySubId(subID)
	}
	threads := e.Thread
	if threads <= 0 || threads > len(nodes) {
		threads = len(nodes)
	}
	if threads > task.MaxThread() {
		threads = task.MaxThread()
	}
	if threads == 0 || len(nodes) == 0 {
		log.Warnf("alive check task failed, no nodes")
		return checkModel.Result{
			Msg:      "no nodes",
			LastRun:  time.Now(),
			Duration: time.Since(startTime).Milliseconds(),
		}
	}
	sem := make(chan struct{}, threads)
	defer close(sem)

	var wg sync.WaitGroup
	for _, nd := range nodes {
		sem <- struct{}{}
		wg.Add(1)
		n := nd
		task.Submit(func() {
			defer func() {
				<-sem
				wg.Done()
			}()
			var raw map[string]any
			if err := yaml.Unmarshal(n.Raw, &raw); err != nil {
				log.Warnf("yaml.Unmarshal failed: %v", err)
				return
			}
			start := time.Now()
			alive := e.detect(ctx, raw)
			if alive {
				log.Debugf("Node %s is alive ✔", raw["name"].(string))
				atomic.AddInt64(&aliveCount, 1)
				n.Info.SetAliveStatus(nodeModel.Alive, true)
				n.Info.Delay.Update(uint16(time.Since(start).Milliseconds()))
				log.Debugf("Node %s delay: %dms", raw["name"].(string), n.Info.Delay.Average())
				atomic.AddInt64(&totalDelay, int64(n.Info.Delay.Average()))
			} else {
				log.Debugf("Node %s is dead ✘", raw["name"].(string))
				atomic.AddInt64(&deadCount, 1)
				n.Info.SetAliveStatus(nodeModel.Alive, false)
				n.Info.Delay.Update(uint16(65535))
			}

		})
	}
	wg.Wait()
	avgDelay := int64(0)
	if aliveCount > 0 {
		avgDelay = totalDelay / aliveCount
	}
	log.Debugf("alive check task end, alive: %d, dead: %d, average delay: %dms", aliveCount, deadCount, avgDelay)
	return checkModel.Result{
		Msg:      fmt.Sprintf("success, alive: %d, dead: %d, average delay: %dms", aliveCount, deadCount, avgDelay),
		LastRun:  time.Now(),
		Duration: time.Since(startTime).Milliseconds(),
		Extra: map[string]any{
			"alive": aliveCount,
			"dead":  deadCount,
			"delay": avgDelay,
		},
	}
}

func (e *Alive) detect(ctx context.Context, raw map[string]any) bool {
	client := mihomo.Proxy(raw)
	if client == nil {
		return false
	}
	client.Timeout = time.Duration(e.Timeout) * time.Second
	defer client.Release()
	request, err := http.NewRequestWithContext(ctx, "GET", e.URL, nil)
	if err != nil {
		return false
	}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == e.ExptectCode
}

func init() {
	register.Check(&Alive{})
}
