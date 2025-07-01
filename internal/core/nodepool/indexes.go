package nodepool

// GetNodesBySubLinkID 根据订阅链接ID获取节点列表
func GetNodesBySubLinkID(subLinkID int64) []*NodeInfo {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	if nodes, exists := globalPool.indexes.subLinkID[subLinkID]; exists {
		return copyNodeInfoSlice(nodes)
	}
	return nil
}

// GetNodesByType 根据节点类型获取节点列表
func GetNodesByType(nodeType string) []*NodeInfo {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	if nodes, exists := globalPool.indexes.nodeType[nodeType]; exists {
		return copyNodeInfoSlice(nodes)
	}
	return nil
}

// GetNodesByMultipleFilters 根据多个条件组合查询节点
func GetNodesByMultipleFilters(subLinkID *int64, nodeType *string) []*NodeInfo {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	// 单条件查询
	if subLinkID != nil && nodeType == nil {
		return getNodesBySubLinkIDUnsafe(*subLinkID)
	}
	if nodeType != nil && subLinkID == nil {
		return getNodesByTypeUnsafe(*nodeType)
	}

	// 双条件查询：找到两个索引的交集
	if subLinkID != nil && nodeType != nil {
		return getNodesIntersection(*subLinkID, *nodeType)
	}

	// 无条件查询：返回所有节点
	return getAllNodesUnsafe()
}

// getNodesBySubLinkIDUnsafe 内部使用，不加锁
func getNodesBySubLinkIDUnsafe(subLinkID int64) []*NodeInfo {
	if nodes, exists := globalPool.indexes.subLinkID[subLinkID]; exists {
		return copyNodeInfoSlice(nodes)
	}
	return nil
}

// getNodesByTypeUnsafe 内部使用，不加锁
func getNodesByTypeUnsafe(nodeType string) []*NodeInfo {
	if nodes, exists := globalPool.indexes.nodeType[nodeType]; exists {
		return copyNodeInfoSlice(nodes)
	}
	return nil
}

// getNodesIntersection 获取两个条件的交集
func getNodesIntersection(subLinkID int64, nodeType string) []*NodeInfo {
	subLinkNodes, subLinkExists := globalPool.indexes.subLinkID[subLinkID]
	typeNodes, typeExists := globalPool.indexes.nodeType[nodeType]

	if !subLinkExists || !typeExists {
		return nil
	}

	// 选择较小的集合进行遍历
	var smaller, larger []*NodeInfo
	if len(subLinkNodes) <= len(typeNodes) {
		smaller, larger = subLinkNodes, typeNodes
	} else {
		smaller, larger = typeNodes, subLinkNodes
	}

	// 对于小集合使用线性搜索，大集合使用map查找
	if len(larger) <= 50 {
		return findIntersectionLinear(smaller, larger)
	}
	return findIntersectionWithMap(smaller, larger)
}

// findIntersectionLinear 线性搜索交集
func findIntersectionLinear(smaller, larger []*NodeInfo) []*NodeInfo {
	result := make([]*NodeInfo, 0, len(smaller))
	for _, smallNode := range smaller {
		for _, largeNode := range larger {
			if smallNode == largeNode {
				result = append(result, smallNode)
				break
			}
		}
	}
	return result
}

// findIntersectionWithMap 使用map查找交集
func findIntersectionWithMap(smaller, larger []*NodeInfo) []*NodeInfo {
	largerMap := make(map[*NodeInfo]bool, len(larger))
	for _, node := range larger {
		largerMap[node] = true
	}

	result := make([]*NodeInfo, 0, len(smaller))
	for _, node := range smaller {
		if largerMap[node] {
			result = append(result, node)
		}
	}
	return result
}

// getAllNodesUnsafe 获取所有节点，内部使用，不加锁
func getAllNodesUnsafe() []*NodeInfo {
	result := make([]*NodeInfo, 0, len(globalPool.indexes.uniqueKey))
	for _, nodeInfo := range globalPool.indexes.uniqueKey {
		result = append(result, nodeInfo)
	}
	return result
}

// copyNodeInfoSlice 复制节点信息切片，避免外部修改
func copyNodeInfoSlice(nodes []*NodeInfo) []*NodeInfo {
	result := make([]*NodeInfo, len(nodes))
	copy(result, nodes)
	return result
}

// addToBasicIndexes 将节点信息添加到索引
func (p *NodePool) addToBasicIndexes(uniqueKey uint64, nodeInfo *NodeInfo, nodeType string, subLinkID int64) {
	p.indexes.uniqueKey[uniqueKey] = nodeInfo
	p.indexes.subLinkID[subLinkID] = append(p.indexes.subLinkID[subLinkID], nodeInfo)
	p.indexes.nodeType[nodeType] = append(p.indexes.nodeType[nodeType], nodeInfo)
}
