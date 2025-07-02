package nodepool

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/bestruirui/bestsub/internal/utils/log"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/cespare/xxhash/v2"
)

var (
	globalPool      *NodePool
	nodeInfoPool    sync.Pool // NodeInfo 对象池
	stringSlicePool sync.Pool // 字符串切片对象池
)

// init 初始化全局节点池
func init() {
	globalPool = &NodePool{
		collection: &Collection{},
		indexes:    createNodeIndexes(),
	}

	// 初始化对象池
	nodeInfoPool.New = func() any {
		return &NodeInfo{}
	}

	// 初始化字符串切片对象池
	stringSlicePool.New = func() any {
		slice := make([]string, 0, 8)
		return &slice
	}
}

// findNodeTypeAndSubLinkID 根据节点信息查找节点类型和订阅链接ID
func (p *NodePool) findNodeTypeAndSubLinkID(targetNodeInfo *NodeInfo) (string, int64) {
	// 遍历所有节点类型索引
	for nodeType, nodeInfos := range p.indexes.nodeType {
		for _, nodeInfo := range nodeInfos {
			if nodeInfo == targetNodeInfo {
				// 找到节点类型，现在查找订阅链接ID
				for subLinkID, subLinkNodeInfos := range p.indexes.subLinkID {
					for _, subLinkNodeInfo := range subLinkNodeInfos {
						if subLinkNodeInfo == targetNodeInfo {
							return nodeType, subLinkID
						}
					}
				}
			}
		}
	}
	return "", 0
}

// removeFromCollection 从集合中移除指定位置的节点
func (p *NodePool) removeFromCollection(nodeType string, arrayIndex int) bool {
	collectionValue := reflect.ValueOf(p.collection).Elem()
	cache := getReflectCache()

	// 使用反射缓存获取字段索引
	fieldIndex, exists := cache.collectionFields[nodeType]
	if !exists {
		return false
	}

	sliceField := collectionValue.Field(fieldIndex)
	if sliceField.Kind() != reflect.Slice {
		return false
	}

	// 检查索引是否有效
	if arrayIndex < 0 || arrayIndex >= sliceField.Len() {
		return false
	}

	// 移除元素：将最后一个元素移动到要删除的位置，然后缩短切片
	lastIndex := sliceField.Len() - 1
	if arrayIndex != lastIndex {
		// 将最后一个元素复制到要删除的位置
		lastElement := sliceField.Index(lastIndex)
		sliceField.Index(arrayIndex).Set(lastElement)

		// 更新被移动节点的索引信息
		p.updateMovedNodeIndex(nodeType, lastIndex, arrayIndex)
	}

	// 缩短切片
	newSlice := sliceField.Slice(0, lastIndex)
	sliceField.Set(newSlice)

	return true
}

// updateMovedNodeIndex 更新被移动节点的索引信息
func (p *NodePool) updateMovedNodeIndex(nodeType string, oldIndex, newIndex int) {
	// 遍历所有索引，找到原来在 oldIndex 位置的节点，更新其 ArrayIndex
	for _, nodeInfos := range p.indexes.nodeType[nodeType] {
		if nodeInfos.ArrayIndex == oldIndex {
			nodeInfos.ArrayIndex = newIndex
			break
		}
	}
}

// removeFromAllIndexes 从所有索引中移除节点
func (p *NodePool) removeFromAllIndexes(uniqueKey uint64, nodeInfo *NodeInfo, nodeType string, subLinkID int64) {
	// 从唯一键索引中移除
	delete(p.indexes.uniqueKey, uniqueKey)

	// 从节点类型索引中移除
	p.removeFromSliceIndex(p.indexes.nodeType[nodeType], nodeInfo, func(slice []*NodeInfo) {
		p.indexes.nodeType[nodeType] = slice
	})

	// 从订阅链接ID索引中移除
	p.removeFromSliceIndex(p.indexes.subLinkID[subLinkID], nodeInfo, func(slice []*NodeInfo) {
		p.indexes.subLinkID[subLinkID] = slice
	})
}

// removeFromSliceIndex 从切片索引中移除指定的节点信息
func (p *NodePool) removeFromSliceIndex(slice []*NodeInfo, targetNodeInfo *NodeInfo, updateFunc func([]*NodeInfo)) {
	for i, nodeInfo := range slice {
		if nodeInfo == targetNodeInfo {
			// 移除元素：将最后一个元素移动到当前位置，然后缩短切片
			lastIndex := len(slice) - 1
			if i != lastIndex {
				slice[i] = slice[lastIndex]
			}
			slice = slice[:lastIndex]
			updateFunc(slice)
			break
		}
	}
}

// processCollection 处理整个集合
func (p *NodePool) processCollection(collection *Collection, subLinkID int64) int {
	addedCount := 0
	collectionValue := reflect.ValueOf(collection).Elem()
	for i := 0; i < collectionValue.NumField(); i++ {
		field := collectionValue.Field(i)
		if field.Kind() != reflect.Slice || field.Len() == 0 {
			continue
		}
		fieldName := collectionValue.Type().Field(i).Name
		p.ensureCapacityDirect(field, field.Len())
		addedCount += p.processSlice(field, fieldName, subLinkID)
	}

	return addedCount
}

// processSlice 处理单个切片
func (p *NodePool) processSlice(sliceValue reflect.Value, nodeType string, subLinkID int64) int {
	addedCount := 0
	for i := 0; i < sliceValue.Len(); i++ {
		nodeValue := sliceValue.Index(i)
		if p.addNode(nodeValue, nodeType, subLinkID) {
			addedCount++
		}
	}

	return addedCount
}

// addNode 添加单个节点
func (p *NodePool) addNode(nodeValue reflect.Value, nodeType string, subLinkID int64) bool {
	// 生成唯一键
	uniqueKey := p.generateUniqueKey(nodeValue, nodeType)
	if uniqueKey == 0 {
		return false
	}

	// 检查是否已存在
	if _, exists := p.indexes.uniqueKey[uniqueKey]; exists {
		log.Debugf("node already exists: %v", nodeValue.FieldByName("Config").FieldByName("Name").String())
		return false
	}

	infoField := nodeValue.FieldByName("Info")
	if infoField.IsValid() && infoField.CanSet() {
		fieldName := infoField.FieldByName("UniqueKey")
		if fieldName.IsValid() && fieldName.CanSet() && fieldName.Kind() == reflect.Uint64 {
			fieldName.SetUint(uniqueKey)
		}
		fieldName = infoField.FieldByName("Id")
		if fieldName.IsValid() && fieldName.CanSet() && fieldName.Kind() == reflect.Int64 {
			fieldName.SetInt(timeutils.Now().Unix())
		}
	}

	// 添加到集合并获取索引
	arrayIndex := p.appendToCollection(nodeValue, nodeType)
	if arrayIndex == -1 {
		return false
	}

	// 创建节点信息
	nodeInfo := p.createNodeInfo(arrayIndex)

	// 添加到索引
	p.addToBasicIndexes(uniqueKey, nodeInfo, nodeType, subLinkID)

	return true
}

// createNodeInfo 创建节点信息（使用对象池）
func (p *NodePool) createNodeInfo(arrayIndex int) *NodeInfo {
	nodeInfo := nodeInfoPool.Get().(*NodeInfo)
	nodeInfo.ArrayIndex = arrayIndex
	return nodeInfo
}

// generateUniqueKey 生成节点唯一键
func (p *NodePool) generateUniqueKey(nodeValue reflect.Value, nodeType string) uint64 {
	configField := nodeValue.FieldByName("Config")
	if !configField.IsValid() {
		return 0
	}

	keyPartsPtr := stringSlicePool.Get().(*[]string)
	keyParts := *keyPartsPtr
	keyParts = keyParts[:0]
	defer func() {
		if cap(keyParts) <= maxStringSliceCapacity {
			*keyPartsPtr = keyParts[:0]
			stringSlicePool.Put(keyPartsPtr)
		}
	}()

	cache := getReflectCache()

	for _, fieldInfo := range cache.nodeConfigFields[nodeType] {

		var field reflect.Value
		if fieldInfo.isEmbedded {
			if fieldInfo.embeddedIndex < configField.NumField() {
				embeddedField := configField.Field(fieldInfo.embeddedIndex)
				if embeddedField.IsValid() && fieldInfo.index < embeddedField.NumField() {
					field = embeddedField.Field(fieldInfo.index)
				}
			}
		} else {
			if fieldInfo.index < configField.NumField() {
				field = configField.Field(fieldInfo.index)
			}
		}
		if field.IsValid() && uniqueKeyFields[fieldInfo.name] {
			keyParts = append(keyParts, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return generateUniqueKey(keyParts...)
}

// ensureCapacityDirect 确保切片容量足够
func (p *NodePool) ensureCapacityDirect(sliceField reflect.Value, additionalCount int) {
	if additionalCount <= 0 || sliceField.Kind() != reflect.Slice {
		return
	}

	currentLen := sliceField.Len()
	currentCap := sliceField.Cap()

	if currentCap < currentLen+additionalCount {
		newCap := max(currentLen+additionalCount, currentCap*3/2)
		newSlice := reflect.MakeSlice(sliceField.Type(), currentLen, newCap)
		reflect.Copy(newSlice, sliceField)
		sliceField.Set(newSlice)
	}
}

// appendToCollection 添加节点到集合
func (p *NodePool) appendToCollection(nodeValue reflect.Value, nodeType string) int {
	collectionValue := reflect.ValueOf(p.collection).Elem()
	cache := getReflectCache()

	// 使用反射缓存获取字段索引
	fieldIndex, exists := cache.collectionFields[nodeType]
	if !exists {
		return -1
	}

	sliceField := collectionValue.Field(fieldIndex)
	if sliceField.Kind() != reflect.Slice {
		return -1
	}

	newSlice := reflect.Append(sliceField, nodeValue)
	sliceField.Set(newSlice)
	return newSlice.Len() - 1
}

// getOrCreateIteratorManager 获取或创建迭代器管理器
func getOrCreateIteratorManager() *IteratorManager {
	iteratorMgrOnce.Do(func() {
		globalIteratorMgr = &IteratorManager{
			iterators: make(map[int64]*NodeIterator),
		}
	})
	return globalIteratorMgr
}

// getOrCreateIterator 获取或创建指定订阅链接的迭代器
func getOrCreateIterator(subLinkID int64) *NodeIterator {
	mgr := getOrCreateIteratorManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if iterator, exists := mgr.iterators[subLinkID]; exists {
		return iterator
	}

	// 创建新的迭代器，使用反射获取节点类型名称
	iterator := &NodeIterator{
		typeNames: GetNodeTypeNames(),
		subLinkID: subLinkID,
	}
	mgr.iterators[subLinkID] = iterator
	return iterator
}

// isNodeBelongsToSubLink 检查节点是否属于指定的订阅链接
func (p *NodePool) isNodeBelongsToSubLink(arrayIndex int, subLinkID int64) bool {
	// 从subLinkID索引中查找，直接比较ArrayIndex即可
	// 因为同一个subLinkID下的节点ArrayIndex是唯一的
	if nodes, exists := p.indexes.subLinkID[subLinkID]; exists {
		for _, nodeInfo := range nodes {
			if nodeInfo.ArrayIndex == arrayIndex {
				return true
			}
		}
	}
	return false
}

// generateUniqueKey 生成唯一键
func generateUniqueKey(args ...string) uint64 {
	if len(args) == 0 {
		return 0
	}
	d := xxhash.New()
	for _, arg := range args {
		d.WriteString(arg)
	}
	return d.Sum64()
}
