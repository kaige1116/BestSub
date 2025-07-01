package nodepool

import (
	"reflect"

	"github.com/bestruirui/bestsub/internal/models/node"
)

// init 初始化反射缓存
func init() {
	reflectCacheInstance = &reflectCache{
		collectionFields: make(map[string]int),
		nodeConfigFields: make(map[string][]fieldInfo),
	}

	// 缓存Collection字段信息
	collectionType := reflect.TypeOf(node.Collection{})
	for i := 0; i < collectionType.NumField(); i++ {
		fieldName := collectionType.Field(i).Name
		reflectCacheInstance.collectionFields[fieldName] = i
	}

	// 缓存所有节点类型的Config字段信息
	cacheAllNodeConfigFields()

	// 缓存节点类型名称
	cachedTypeNames = make([]string, 0, len(reflectCacheInstance.collectionFields))
	for typeName := range reflectCacheInstance.collectionFields {
		cachedTypeNames = append(cachedTypeNames, typeName)
	}
}

// getReflectCache 获取反射缓存
func getReflectCache() *reflectCache {
	return reflectCacheInstance
}

// cacheAllNodeConfigFields 缓存所有节点类型的Config字段信息
func cacheAllNodeConfigFields() {
	collectionType := reflect.TypeOf(node.Collection{})

	// 遍历Collection中的所有字段（每个字段代表一种节点类型）
	for i := 0; i < collectionType.NumField(); i++ {
		field := collectionType.Field(i)
		nodeType := field.Name

		// 获取节点切片的元素类型
		sliceType := field.Type
		if sliceType.Kind() != reflect.Slice {
			continue
		}

		nodeStructType := sliceType.Elem()
		if nodeStructType.Kind() != reflect.Struct {
			continue
		}

		// 查找Config字段并直接缓存其字段信息
		if configField, found := nodeStructType.FieldByName("Config"); found {
			cacheConfigFields(nodeType, configField.Type)
		}
	}
}

// cacheConfigFields 缓存Config字段的字符串类型字段信息
func cacheConfigFields(nodeType string, configType reflect.Type) {
	fields := make([]fieldInfo, 0, configType.NumField())

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type.Kind() == reflect.String {
			fields = append(fields, fieldInfo{
				index: i,
				name:  field.Name,
			})
		}
	}

	reflectCacheInstance.nodeConfigFields[nodeType] = fields
}

// GetNodeTypeNames 获取所有节点类型名称（已缓存）
func GetNodeTypeNames() []string {
	return cachedTypeNames
}
