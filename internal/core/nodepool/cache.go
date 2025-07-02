package nodepool

import (
	"reflect"
)

// init 初始化反射缓存
func init() {
	reflectCacheInstance = &reflectCache{
		collectionFields: make(map[string]int),
		nodeConfigFields: make(map[string][]fieldInfo),
	}

	// 缓存Collection字段信息
	collectionType := reflect.TypeOf(Collection{})
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
	collectionType := reflect.TypeOf(Collection{})

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

		// 处理内嵌字段（如 BaseConfig）
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// 递归处理内嵌结构体的字段
			embeddedFields := cacheEmbeddedFields(field.Type, i)
			fields = append(fields, embeddedFields...)
		} else if field.Type.Kind() == reflect.String || field.Type.Kind() == reflect.Interface {
			fields = append(fields, fieldInfo{
				index:         i,
				name:          field.Name,
				isEmbedded:    false,
				embeddedIndex: -1,
			})
		}
	}
	reflectCacheInstance.nodeConfigFields[nodeType] = fields
}

// cacheEmbeddedFields 缓存内嵌结构体的字符串字段信息
func cacheEmbeddedFields(embeddedType reflect.Type, parentIndex int) []fieldInfo {
	var fields []fieldInfo

	for i := 0; i < embeddedType.NumField(); i++ {
		field := embeddedType.Field(i)
		if field.Type.Kind() == reflect.String || field.Type.Kind() == reflect.Interface {
			fields = append(fields, fieldInfo{
				index:         i,
				name:          field.Name,
				isEmbedded:    true,
				embeddedIndex: parentIndex,
			})
		}
	}

	return fields
}

// GetNodeTypeNames 获取所有节点类型名称（已缓存）
func GetNodeTypeNames() []string {
	return cachedTypeNames
}
