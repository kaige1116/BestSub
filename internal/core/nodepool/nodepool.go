package nodepool

import (
	"github.com/bestruirui/bestsub/internal/models/node"
)

const NodePoolSize = 100

// 初始化节点池
func init() {
	globalPool = &pool{
		nodes:     make(map[uint64]node.Data, NodePoolSize),
		subs:      make(map[int64][]uint64, NodePoolSize),
		iterators: make(map[int64]int, NodePoolSize),
		index: &index{
			delay:       make([]uint64, NodePoolSize),
			speedUp:     make([]uint64, NodePoolSize),
			speedDown:   make([]uint64, NodePoolSize),
			risk:        make([]uint64, NodePoolSize),
			country:     make(map[uint16][]uint64, NodePoolSize),
			aliveStatus: make(map[uint16][]uint64, NodePoolSize),
		},
	}
}

// Add 添加节点到节点池
func Add(nodes *[]node.Data, subLinkID int64) int {
	if len(*nodes) == 0 {
		return 0
	}

	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	addedCount := 0

	for _, nodeData := range *nodes {

		uniqueKey := nodeData.Info.UniqueKey

		if _, exists := globalPool.nodes[uniqueKey]; exists {
			continue
		}

		globalPool.nodes[uniqueKey] = nodeData

		globalPool.subs[subLinkID] = append(globalPool.subs[subLinkID], uniqueKey)

		updateIndexes(uniqueKey, &nodeData.Info)

		addedCount++
	}

	(*nodes) = nil

	return addedCount
}

// GetNextNode 获取下一个节点
func GetNextNode(subLinkID int64) node.Data {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	nodeKeys, exists := globalPool.subs[subLinkID]
	if !exists || len(nodeKeys) == 0 {
		return node.Data{}
	}

	currentPos := globalPool.iterators[subLinkID]
	if currentPos >= len(nodeKeys) {
		currentPos = 0
	}

	uniqueKey := nodeKeys[currentPos]
	nodeData, exists := globalPool.nodes[uniqueKey]
	if !exists {
		globalPool.iterators[subLinkID] = 0
		return node.Data{}
	}

	globalPool.mu.RUnlock()
	globalPool.mu.Lock()
	globalPool.iterators[subLinkID] = (currentPos + 1) % len(nodeKeys)
	globalPool.mu.Unlock()
	globalPool.mu.RLock()

	result := node.Data{
		Config: make([]byte, len(nodeData.Config)),
		Info:   nodeData.Info,
	}
	copy(result.Config, nodeData.Config)

	return result
}

// FilterNodes 根据条件筛选节点
func FilterNodes(condition FilterCondition) []node.Data {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	candidateKeys := filterByIndexes(condition)

	var results []node.Data
	for i, uniqueKey := range candidateKeys {
		nodeData, exists := globalPool.nodes[uniqueKey]
		if !exists {
			continue
		}

		result := node.Data{
			Config: make([]byte, len(nodeData.Config)),
			Info:   nodeData.Info,
		}
		copy(result.Config, nodeData.Config)
		results = append(results, result)

		candidateKeys[i] = 0
	}

	candidateKeys = nil

	return results
}

// RemoveNodes 批量删除节点
func RemoveNodes(uniqueKeys []uint64) int {
	if len(uniqueKeys) == 0 {
		return 0
	}

	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	removedCount := 0
	affectedSubLinks := make(map[int64]bool)

	for i, uniqueKey := range uniqueKeys {
		nodeData, exists := globalPool.nodes[uniqueKey]
		if !exists {
			continue
		}

		delete(globalPool.nodes, uniqueKey)

		for subLinkID, keys := range globalPool.subs {
			for j, key := range keys {
				if key == uniqueKey {
					globalPool.subs[subLinkID] = append(keys[:j], keys[j+1:]...)
					affectedSubLinks[subLinkID] = true
					break
				}
			}
		}

		removeFromIndexes(uniqueKey, &nodeData.Info)

		removedCount++

		uniqueKeys[i] = 0
	}

	for subLinkID := range affectedSubLinks {
		globalPool.iterators[subLinkID] = 0
	}

	affectedSubLinks = nil

	return removedCount
}

// RemoveBySubLink 删除订阅链接的所有节点
func RemoveBySubLink(subLinkID int64) int {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	nodeKeys, exists := globalPool.subs[subLinkID]
	if !exists || len(nodeKeys) == 0 {
		return 0
	}

	removedCount := 0

	for i, uniqueKey := range nodeKeys {
		if nodeData, exists := globalPool.nodes[uniqueKey]; exists {
			delete(globalPool.nodes, uniqueKey)
			removeFromIndexes(uniqueKey, &nodeData.Info)
			removedCount++
		}
		nodeKeys[i] = 0
	}

	delete(globalPool.subs, subLinkID)

	delete(globalPool.iterators, subLinkID)

	nodeKeys = nil

	return removedCount
}

// CleanupHighDelayNodes 清理高延迟节点
func CleanupHighDelayNodes(maxAvgDelay int) int {
	globalPool.mu.RLock()

	var toRemove []uint64

	for uniqueKey, nodeData := range globalPool.nodes {
		avgDelay := calculateAverageDelay(nodeData.Info.Delay)
		if avgDelay > uint16(maxAvgDelay) {
			toRemove = append(toRemove, uniqueKey)
		}
	}

	globalPool.mu.RUnlock()

	result := RemoveNodes(toRemove)

	toRemove = nil

	return result
}

// UpdateNodeInfo 更新节点信息
func UpdateNodeInfo(info node.Info) bool {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	nodeData, exists := globalPool.nodes[info.UniqueKey]
	if !exists {
		return false
	}

	removeFromIndexes(info.UniqueKey, &nodeData.Info)

	nodeData.Info = info

	updateIndexes(info.UniqueKey, &info)

	return true
}

// ResetIterator 重置订阅链接的迭代器位置
func ResetIterator(subLinkID int64) {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	globalPool.iterators[subLinkID] = 0
}

// Reset 重置节点池（清空所有数据）
func Reset() {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	globalPool.nodes = make(map[uint64]node.Data)
	globalPool.subs = make(map[int64][]uint64)
	globalPool.iterators = make(map[int64]int)

	globalPool.index.delay = make([]uint64, 0)
	globalPool.index.speedUp = make([]uint64, 0)
	globalPool.index.speedDown = make([]uint64, 0)
	globalPool.index.risk = make([]uint64, 0)
	globalPool.index.country = make(map[uint16][]uint64)
	globalPool.index.aliveStatus = make(map[uint16][]uint64)
}

// GetStats 获取节点统计信息
func GetStats() PoolStats {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	stats := PoolStats{
		TotalNodes:     len(globalPool.nodes),
		NodesByCountry: make(map[uint16]int),
		NodesByStatus:  make(map[uint16]int),
		SubLinkCount:   len(globalPool.subs),
	}

	for _, nodeData := range globalPool.nodes {
		stats.NodesByCountry[nodeData.Info.Country]++
		stats.NodesByStatus[nodeData.Info.AliveStatus]++
	}

	return stats
}
