package checker

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/nodepool"
	checkModel "github.com/bestruirui/bestsub/internal/models/check"
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

type Alive struct {
	URL         string   `json:"url" description:"URL"`
	ExptectCode int      `json:"exptect_code" description:"期望状态码"`
	SubID       []uint16 `json:"sub_id" description:"订阅ID"`
	Thread      int      `json:"thread" description:"线程数"`
	Timeout     int      `json:"timeout" description:"超时时间"`
}
type Result struct {
	AliveCount uint16 `json:"alive_count" description:"存活节点数量"`
	DeadCount  uint16 `json:"dead_count" description:"死亡节点数量"`
	Delay      uint16 `json:"delay" description:"平均延迟"`
}

func (e *Alive) Init() error {
	return nil
}

func (e *Alive) Run(ctx context.Context, log *log.Logger) checkModel.Result {
	log.Infof("alive %s 任务执行开始 %d", e.URL, e.SubID)
	pool, _ := ants.NewPool(e.Thread)
	defer pool.Release()
	var wg sync.WaitGroup
	for _, subID := range e.SubID {
		subStorage := nodepool.GetPoolBySubID(subID, 0)
		for _, node := range subStorage.GetAllNode() {
			wg.Add(1)
			pool.Submit(func() {
				var raw map[string]any
				if err := json.Unmarshal(node.Raw, &raw); err != nil {
					return
				}
				start := time.Now()
				alive := e.detect(ctx, raw)
				if alive {
					log.Infof("节点 %s 存活 ✔", raw["name"].(string))
				} else {
					log.Infof("节点 %s 死亡 ✘", raw["name"].(string))
				}
				node.Info.SetAliveStatus(nodeModel.Alive, alive)
				node.Info.Delay.Update(uint16(time.Since(start).Milliseconds()))
				wg.Done()
			})
		}
	}
	wg.Wait()
	nodepool.UpdateStats()
	log.Infof("alive %s 任务执行结束", e.URL)
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
	if response.StatusCode != e.ExptectCode {
		return false
	}
	return true
}

func init() {
	register.Check(&Alive{})
}
