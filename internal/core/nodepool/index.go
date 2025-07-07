package nodepool

import (
	"sort"

	"github.com/bestruirui/bestsub/internal/models/node"
)

// updateIndexes 更新所有索引
func updateIndexes(uniqueKey uint64, info *node.Info) {
	// 更新数值索引（有序插入）
	insertSortedByDelay(uniqueKey, info)
	insertSortedBySpeedUp(uniqueKey, info)
	insertSortedBySpeedDown(uniqueKey, info)
	insertSortedByRisk(uniqueKey, info)

	// 更新映射索引
	updateCountryIndex(uniqueKey, info.Country)
	updateAliveStatusIndex(uniqueKey, info.AliveStatus)
}

// removeFromIndexes 从所有索引中删除
func removeFromIndexes(uniqueKey uint64, info *node.Info) {
	// 从数值索引中删除
	removeFromSortedSlice(&globalPool.index.delay, uniqueKey)
	removeFromSortedSlice(&globalPool.index.speedUp, uniqueKey)
	removeFromSortedSlice(&globalPool.index.speedDown, uniqueKey)
	removeFromSortedSlice(&globalPool.index.risk, uniqueKey)

	// 从映射索引中删除
	removeFromCountryIndex(uniqueKey, info.Country)
	removeFromAliveStatusIndex(uniqueKey, info.AliveStatus)
}

// insertSortedByDelay 按延迟有序插入
func insertSortedByDelay(uniqueKey uint64, info *node.Info) {
	avgDelay := calculateAverageDelay(info.Delay)
	insertSortedByValue(&globalPool.index.delay, uniqueKey, uint64(avgDelay), true) // 升序
}

// insertSortedBySpeedUp 按上传速度有序插入
func insertSortedBySpeedUp(uniqueKey uint64, info *node.Info) {
	avgSpeedUp := calculateAverageSpeedUp(info.SpeedUp)
	insertSortedByValue(&globalPool.index.speedUp, uniqueKey, uint64(avgSpeedUp), false) // 降序
}

// insertSortedBySpeedDown 按下载速度有序插入
func insertSortedBySpeedDown(uniqueKey uint64, info *node.Info) {
	avgSpeedDown := calculateAverageSpeedDown(info.SpeedDown)
	insertSortedByValue(&globalPool.index.speedDown, uniqueKey, uint64(avgSpeedDown), false) // 降序
}

// insertSortedByRisk 按风险等级有序插入
func insertSortedByRisk(uniqueKey uint64, info *node.Info) {
	insertSortedByValue(&globalPool.index.risk, uniqueKey, uint64(info.Risk), true) // 升序，低风险在前
}

// insertSortedByValue 通用的有序插入函数
func insertSortedByValue(slice *[]uint64, uniqueKey uint64, value uint64, ascending bool) {
	pos := sort.Search(len(*slice), func(i int) bool {
		existingKey := (*slice)[i]
		existingNode := globalPool.nodes[existingKey]

		var existingValue uint64
		switch slice {
		case &globalPool.index.delay:
			existingValue = uint64(calculateAverageDelay(existingNode.Info.Delay))
		case &globalPool.index.speedUp:
			existingValue = uint64(calculateAverageSpeedUp(existingNode.Info.SpeedUp))
		case &globalPool.index.speedDown:
			existingValue = uint64(calculateAverageSpeedDown(existingNode.Info.SpeedDown))
		case &globalPool.index.risk:
			existingValue = uint64(existingNode.Info.Risk)
		}

		if ascending {
			return existingValue >= value
		} else {
			return existingValue <= value
		}
	})

	*slice = append(*slice, 0)
	copy((*slice)[pos+1:], (*slice)[pos:])
	(*slice)[pos] = uniqueKey
}

// removeFromSortedSlice 从有序切片中删除元素
func removeFromSortedSlice(slice *[]uint64, uniqueKey uint64) {
	for i, key := range *slice {
		if key == uniqueKey {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}

// updateCountryIndex 更新国家索引
func updateCountryIndex(uniqueKey uint64, country uint16) {
	globalPool.index.country[country] = append(globalPool.index.country[country], uniqueKey)
}

// removeFromCountryIndex 从国家索引中删除
func removeFromCountryIndex(uniqueKey uint64, country uint16) {
	keys := globalPool.index.country[country]
	for i, key := range keys {
		if key == uniqueKey {
			globalPool.index.country[country] = append(keys[:i], keys[i+1:]...)
			if len(globalPool.index.country[country]) == 0 {
				delete(globalPool.index.country, country)
			}
			keys = nil
			break
		}
	}
}

// updateAliveStatusIndex 更新存活状态索引
func updateAliveStatusIndex(uniqueKey uint64, aliveStatus uint16) {
	globalPool.index.aliveStatus[aliveStatus] = append(globalPool.index.aliveStatus[aliveStatus], uniqueKey)
}

// removeFromAliveStatusIndex 从存活状态索引中删除
func removeFromAliveStatusIndex(uniqueKey uint64, aliveStatus uint16) {
	keys := globalPool.index.aliveStatus[aliveStatus]
	for i, key := range keys {
		if key == uniqueKey {
			globalPool.index.aliveStatus[aliveStatus] = append(keys[:i], keys[i+1:]...)
			if len(globalPool.index.aliveStatus[aliveStatus]) == 0 {
				delete(globalPool.index.aliveStatus, aliveStatus)
			}
			keys = nil
			break
		}
	}
}

// filterByIndexes 使用索引进行筛选
func filterByIndexes(condition FilterCondition) []uint64 {
	var candidateKeys []uint64

	if condition.Country > 0 {
		if keys, exists := globalPool.index.country[condition.Country]; exists {
			candidateKeys = append(candidateKeys, keys...)
		} else {
			return []uint64{}
		}
	} else {
		for uniqueKey := range globalPool.nodes {
			candidateKeys = append(candidateKeys, uniqueKey)
		}
	}

	var filteredKeys []uint64
	for i, uniqueKey := range candidateKeys {
		nodeData, exists := globalPool.nodes[uniqueKey]
		if !exists {
			continue
		}

		info := &nodeData.Info

		if condition.AliveStatus > 0 && (info.AliveStatus&condition.AliveStatus) == 0 {
			continue
		}

		if condition.DelayLessThan > 0 {
			avgDelay := calculateAverageDelay(info.Delay)
			if avgDelay >= condition.DelayLessThan {
				continue
			}
		}

		if condition.SpeedUpMore > 0 {
			avgSpeedUp := calculateAverageSpeedUp(info.SpeedUp)
			if avgSpeedUp <= condition.SpeedUpMore {
				continue
			}
		}

		if condition.SpeedDownMore > 0 {
			avgSpeedDown := calculateAverageSpeedDown(info.SpeedDown)
			if avgSpeedDown <= condition.SpeedDownMore {
				continue
			}
		}

		if condition.RiskLessThan > 0 && info.Risk >= condition.RiskLessThan {
			continue
		}

		filteredKeys = append(filteredKeys, uniqueKey)

		candidateKeys[i] = 0
	}

	candidateKeys = nil

	return filteredKeys
}

// calculateAverageDelay 计算平均延迟
func calculateAverageDelay(delays [5]uint16) uint16 {
	var sum uint32
	count := 0
	for _, delay := range delays {
		if delay > 0 {
			sum += uint32(delay)
			count++
		}
	}
	if count == 0 {
		return 65535
	}
	return uint16(sum / uint32(count))
}

// calculateAverageSpeedUp 计算平均上传速度
func calculateAverageSpeedUp(speeds [5]uint32) uint32 {
	var sum uint64
	count := 0
	for _, speed := range speeds {
		if speed > 0 {
			sum += uint64(speed)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return uint32(sum / uint64(count))
}

// calculateAverageSpeedDown 计算平均下载速度
func calculateAverageSpeedDown(speeds [5]uint32) uint32 {
	var sum uint64
	count := 0
	for _, speed := range speeds {
		if speed > 0 {
			sum += uint64(speed)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return uint32(sum / uint64(count))
}
