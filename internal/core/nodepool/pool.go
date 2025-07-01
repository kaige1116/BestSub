package nodepool

import (
	"fmt"
	"reflect"
	"time"

	"github.com/bestruirui/bestsub/internal/models/node"
)

// Add 添加节点集合到池中
func Add(collection *node.Collection, subLinkID int64) error {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	timestamp := time.Now().Unix()
	addedCount := globalPool.processCollection(collection, subLinkID, timestamp)

	// 更新统计信息
	globalPool.totalNodes += int64(addedCount)

	return nil
}

// RemoveByUniqueKey 根据唯一键从池中删除节点
func RemoveByUniqueKey(uniqueKey uint64) bool {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	// 从唯一键索引中查找节点信息
	nodeInfo, exists := globalPool.indexes.uniqueKey[uniqueKey]
	if !exists {
		return false
	}

	// 找到节点所属的类型和订阅链接ID
	nodeType, subLinkID := globalPool.findNodeTypeAndSubLinkID(nodeInfo)
	if nodeType == "" {
		return false
	}

	// 从集合中移除节点
	if !globalPool.removeFromCollection(nodeType, nodeInfo.ArrayIndex) {
		return false
	}

	// 从所有索引中移除节点
	globalPool.removeFromAllIndexes(uniqueKey, nodeInfo, nodeType, subLinkID)

	// 更新统计信息
	globalPool.totalNodes--

	// 将 NodeInfo 对象返回到对象池
	nodeInfoPool.Put(nodeInfo)

	return true
}

// GetNextNode 获取指定订阅链接的下一个节点，多线程安全，按顺序返回
//
// 参数：
//   - subLinkID: 订阅链接ID，用于过滤特定订阅的节点
//
// 返回值：
//   - string: 节点类型名称
//   - string: 节点Config字段转换为字符串
//   - *node.Info: 节点Info字段的指针，可直接修改
//   - bool: 是否还有更多节点可以遍历
//
// 注意：每次调用都会返回下一个节点的信息，当所有节点遍历完成后返回 ("", "", nil, false)
func GetNextNode(subLinkID int64) (string, string, *node.Info, bool) {
	globalPool.mu.RLock()
	defer globalPool.mu.RUnlock()

	iterator := getOrCreateIterator(subLinkID)
	iterator.mu.Lock()
	defer iterator.mu.Unlock()

	// 如果已经遍历完成，返回空值
	if iterator.finished {
		return "", "", nil, false
	}

	collectionValue := reflect.ValueOf(globalPool.collection).Elem()

	// 寻找下一个有效节点
	for iterator.currentType < len(iterator.typeNames) {
		typeName := iterator.typeNames[iterator.currentType]
		cache := getReflectCache()

		// 使用反射缓存获取字段索引
		fieldIndex, exists := cache.collectionFields[typeName]
		if !exists {
			iterator.currentType++
			iterator.currentIndex = 0
			continue
		}

		sliceField := collectionValue.Field(fieldIndex)
		if sliceField.Kind() != reflect.Slice {
			iterator.currentType++
			iterator.currentIndex = 0
			continue
		}

		// 在当前类型中寻找属于指定subLinkID的节点
		for iterator.currentIndex < sliceField.Len() {
			nodeValue := sliceField.Index(iterator.currentIndex)
			currentArrayIndex := iterator.currentIndex

			// 移动到下一个位置
			iterator.currentIndex++

			// 检查该节点是否属于指定的subLinkID
			if globalPool.isNodeBelongsToSubLink(currentArrayIndex, subLinkID) {
				// 获取节点的Config字段并转换为字符串
				configField := nodeValue.FieldByName("Config")
				var configStr string
				if configField.IsValid() {
					configStr = fmt.Sprintf("%+v", configField.Interface())
				}

				// 获取节点的Info字段指针
				infoField := nodeValue.FieldByName("Info")
				var infoPtr *node.Info
				if infoField.IsValid() && infoField.CanAddr() {
					infoPtr = infoField.Addr().Interface().(*node.Info)
				}

				return typeName, configStr, infoPtr, true
			}
		}

		// 当前类型遍历完成，移动到下一个类型
		iterator.currentType++
		iterator.currentIndex = 0
	}

	// 所有节点都遍历完成
	iterator.finished = true
	return "", "", nil, false
}

// ResetIterator 重置指定订阅链接的迭代器
func ResetIterator(subLinkID int64) {
	mgr := getOrCreateIteratorManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if iterator, exists := mgr.iterators[subLinkID]; exists {
		iterator.mu.Lock()
		iterator.currentType = 0
		iterator.currentIndex = 0
		iterator.finished = false
		iterator.mu.Unlock()
	}
}

// ResetAllIterators 重置所有迭代器
func ResetAllIterators() {
	mgr := getOrCreateIteratorManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	for _, iterator := range mgr.iterators {
		iterator.mu.Lock()
		iterator.currentType = 0
		iterator.currentIndex = 0
		iterator.finished = false
		iterator.mu.Unlock()
	}
}

// Reset 重置节点池（仅用于测试）
func Reset() {
	globalPool.mu.Lock()
	defer globalPool.mu.Unlock()

	// 重置集合
	globalPool.collection = &node.Collection{}

	// 重置索引
	globalPool.indexes = createNodeIndexes()

	// 重置统计信息
	globalPool.totalNodes = 0

	// 重置所有迭代器
	ResetAllIterators()
}
