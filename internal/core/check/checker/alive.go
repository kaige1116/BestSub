package checker

import (
	"context"
	"net/http"
	"sync"
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
	URL         string `json:"url" name:"测试链接" default:"https://www.gstatic.com/generate_204" description:"测试链接"`
	ExptectCode int    `json:"exptect_code" name:"期望状态码" default:"204" description:"期望状态码"`
	Thread      int    `json:"thread" name:"线程数" default:"100"`
	Timeout     int    `json:"timeout" name:"超时时间" default:"10"`
}
type Result struct {
	AliveCount uint16 `json:"alive_count" description:"存活节点数量"`
	DeadCount  uint16 `json:"dead_count" description:"死亡节点数量"`
	Delay      uint16 `json:"delay" description:"平均延迟"`
}

func (e *Alive) Init() error {
	return nil
}

func (e *Alive) Run(ctx context.Context, log *log.Logger, subID []uint16) checkModel.Result {
	log.Infof("alive check task start, alive url: %s, thread: %d", e.URL, e.Thread)
	var nodes []nodeModel.Data
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
	sem := make(chan struct{}, threads)

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
			log.Debugf("Start Name: %v", raw["name"])
			start := time.Now()
			alive := e.detect(ctx, raw)
			if alive {
				log.Debugf("Node %s is alive ✔", raw["name"].(string))
			} else {
				log.Debugf("Node %s is dead ✘", raw["name"].(string))
			}
			n.Info.SetAliveStatus(nodeModel.Alive, alive)
			n.Info.Delay.Update(uint16(time.Since(start).Milliseconds()))
		})
	}
	wg.Wait()
	log.Debugf("alive check task end, alive url: %s", e.URL)
	return checkModel.Result{}
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
