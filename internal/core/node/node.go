package node

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var poolMutex sync.RWMutex
var nodePool []nodeModel.Data
var nodeExist *exist
var nodeProcess *exist

var wgSync sync.WaitGroup
var wgStatus bool
var validNodes []nodeModel.Data
var validMutex sync.Mutex

func InitNodePool(size int) {
	nodePool = make([]nodeModel.Data, 0, size)
	nodeExist = NewExist(size)
	nodeProcess = NewExist(size)
}

func Add(node *[]nodeModel.Base) int {
	var nodesToProcess []nodeModel.Base

	for _, n := range *node {
		if !nodeExist.Exist(n.UniqueKey) && !nodeProcess.Exist(n.UniqueKey) {
			nodeProcess.Add(n.UniqueKey)
			nodesToProcess = append(nodesToProcess, n)
		}
	}

	if len(nodesToProcess) > 0 {
		wgSync.Add(1)
		go func() {
			defer wgSync.Done()
			log.Debugf("nodesToProcess: %d", len(nodesToProcess))
			for _, node := range nodesToProcess {
				task.Submit(func() {
					defer nodeProcess.Remove(node.UniqueKey)
					var raw map[string]any
					if err := json.Unmarshal(node.Raw, &raw); err != nil {
						return
					}
					client := mihomo.Proxy(raw)
					if client == nil {
						return
					}
					defer client.Release()
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.GetConfigInt("node.test_timeout"))*time.Second)
					defer cancel()
					request, err := http.NewRequestWithContext(ctx, "GET", op.GetConfigStr("node.test_url"), nil)
					if err != nil {
						return
					}
					start := time.Now()
					response, err := client.Do(request)
					if err != nil {
						return
					}
					defer response.Body.Close()
					if response.StatusCode != 204 {
						return
					}

					var info nodeModel.Info
					info.Delay.Update(uint16(time.Since(start).Milliseconds()))
					validMutex.Lock()
					validNodes = append(validNodes, nodeModel.Data{
						Base: node,
						Info: &info,
					})
					validMutex.Unlock()
				})
			}
			log.Debugf("nodesToProcess end: %d", len(nodesToProcess))
		}()
		if !wgStatus {
			wgStatus = true
			go func() {
				wgSync.Wait()
				if len(validNodes) > 0 {
					mergeNodesToPool(validNodes)
					log.Debugf("validNodes: %d", len(validNodes))
					validNodes = validNodes[:0]
				}
				wgStatus = false
			}()
		}
	}
	return len(nodesToProcess)
}

func ForEach(fn func(node []byte)) {
	poolMutex.RLock()
	defer poolMutex.RUnlock()
	log.Debugf("nodePool: %d", len(nodePool))
	for _, node := range nodePool {
		fn(node.Raw)
	}
}

func GetAll() []nodeModel.Data {
	poolMutex.RLock()
	defer poolMutex.RUnlock()
	return nodePool
}

func GetBySubId(subId uint16) *[]nodeModel.Data {
	poolMutex.RLock()
	defer poolMutex.RUnlock()
	var result []nodeModel.Data
	for _, node := range nodePool {
		if node.Base.SubId == subId {
			result = append(result, node)
		}
	}
	return &result
}

func GetByFilter(filter nodeModel.Filter) *[]nodeModel.Data {
	poolMutex.RLock()
	defer poolMutex.RUnlock()
	var result []nodeModel.Data
	for _, node := range nodePool {
		if filter.SubId != 0 && node.Base.SubId != filter.SubId {
			continue
		}
		if filter.AliveStatus != 0 && node.Info.AliveStatus&filter.AliveStatus == 0 {
			continue
		}
		if len(filter.Country) > 0 && !slices.Contains(filter.Country, node.Info.Country) {
			continue
		}
		if filter.SpeedUpMore != 0 && node.Info.SpeedUp.Average() < filter.SpeedUpMore {
			continue
		}
		if filter.SpeedDownMore != 0 && node.Info.SpeedDown.Average() < filter.SpeedDownMore {
			continue
		}
		if filter.DelayLessThan != 0 && node.Info.Delay.Average() > filter.DelayLessThan {
			continue
		}
		if filter.RiskLessThan != 0 && node.Info.Risk > filter.RiskLessThan {
			continue
		}
		result = append(result, node)
	}
	return &result
}

func mergeNodesToPool(newNodes []nodeModel.Data) {
	sort.Slice(newNodes, func(i, j int) bool {
		return newNodes[i].Info.Delay.Average() < newNodes[j].Info.Delay.Average()
	})

	poolMutex.Lock()
	defer poolMutex.Unlock()

	poolLen := len(nodePool)
	poolCap := cap(nodePool)

	if poolLen < poolCap {
		log.Debugf("poolLen < poolCap")
		remainingCap := poolCap - poolLen
		if len(newNodes) < remainingCap {
			log.Debugf("newNodes < remainingCap")
			nodePool = append(nodePool, newNodes...)
			for _, node := range newNodes {
				nodeExist.Add(node.Base.UniqueKey)
			}
			return
		} else {
			log.Debugf("newNodes > remainingCap")
			nodePool = append(nodePool, newNodes[:remainingCap]...)
			for _, node := range newNodes[:remainingCap] {
				nodeExist.Add(node.Base.UniqueKey)
			}
			newNodes = newNodes[remainingCap:]
		}
	}

	sort.Slice(nodePool, func(i, j int) bool {
		return nodePool[i].Info.Delay.Average() < nodePool[j].Info.Delay.Average()
	})

	newNodeIndex := 0
	for i := len(nodePool) - 1; i >= 0 && newNodeIndex < len(newNodes); i-- {
		if newNodes[newNodeIndex].Info.Delay.Average() < nodePool[i].Info.Delay.Average() {
			nodeExist.Remove(nodePool[i].Base.UniqueKey)
			nodePool[i] = newNodes[newNodeIndex]
			nodeExist.Add(newNodes[newNodeIndex].Base.UniqueKey)
			newNodeIndex++
		} else {
			return
		}
	}
}
