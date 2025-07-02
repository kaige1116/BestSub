package mihomo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/bestruirui/bestsub/internal/core/nodepool"
	"github.com/bestruirui/bestsub/internal/models/node"
	"gopkg.in/yaml.v2"
)

// 反射缓存结构
type parseCache struct {
	typeBytes map[string][]byte       // 节点类型名到字节切片的映射
	nodeTypes map[string]nodeTypeInfo // 节点类型信息缓存
}

type nodeTypeInfo struct {
	fieldIndex       int          // 在Collection中的字段索引
	nodeType         reflect.Type // 节点结构体类型
	configType       reflect.Type // Config结构体类型
	configArrayType  reflect.Type // Config数组类型 [1]ConfigType
	infoFieldIndex   int          // Info字段在节点结构体中的索引
	configFieldIndex int          // Config字段在节点结构体中的索引
}

var (
	parseCacheInstance *parseCache
	// 预分配常用的Info实例，避免重复创建
	emptyInfo = reflect.ValueOf(node.Info{})
)

// init 初始化反射缓存
func init() {
	parseCacheInstance = &parseCache{
		typeBytes: make(map[string][]byte),
		nodeTypes: make(map[string]nodeTypeInfo),
	}

	// 初始化Collection反射信息
	collectionType := reflect.TypeOf(nodepool.Collection{})

	// 缓存所有支持的节点类型
	for _, nodeTypeName := range node.SupportedNodes {
		// 创建类型检查用的字节切片
		parseCacheInstance.typeBytes[nodeTypeName] = []byte("type: " + nodeTypeName)

		// 获取Collection中对应字段的信息
		fieldName := capitalizeFirst(nodeTypeName)
		if field, found := collectionType.FieldByName(fieldName); found {
			fieldIndex := getFieldIndex(collectionType, fieldName)

			// 获取节点类型和Config类型
			sliceType := field.Type
			if sliceType.Kind() == reflect.Slice {
				nodeType := sliceType.Elem()
				if configField, found := nodeType.FieldByName("Config"); found {
					// 预缓存字段索引
					infoFieldIndex := getFieldIndex(nodeType, "Info")
					configFieldIndex := getFieldIndex(nodeType, "Config")

					// 预创建Config数组类型
					configArrayType := reflect.ArrayOf(1, configField.Type)
					parseCacheInstance.nodeTypes[nodeTypeName] = nodeTypeInfo{
						fieldIndex:       fieldIndex,
						nodeType:         nodeType,
						configType:       configField.Type,
						configArrayType:  configArrayType,
						infoFieldIndex:   infoFieldIndex,
						configFieldIndex: configFieldIndex,
					}
				}
			}
		}
	}
}

// capitalizeFirst 将字符串首字母大写
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}

// getFieldIndex 获取字段在结构体中的索引
func getFieldIndex(structType reflect.Type, fieldName string) int {
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Name == fieldName {
			return i
		}
	}
	return -1
}

// fastNodeTypeCheck 节点类型检查
func fastNodeTypeCheck(nodeYAML *[]byte) string {
	for nodeType, typeBytes := range parseCacheInstance.typeBytes {
		if bytes.Contains(*nodeYAML, typeBytes) {
			return nodeType
		}
	}
	return ""
}

// parseProxyNode 解析单个代理节点
func parseProxyNode(nodeYAML *[]byte, collection *nodepool.Collection) error {
	nodeType := fastNodeTypeCheck(nodeYAML)
	if len(nodeType) == 0 {
		return fmt.Errorf("无法确定节点类型")
	}

	typeInfo, exists := parseCacheInstance.nodeTypes[nodeType]
	if !exists {
		return fmt.Errorf("unsupported proxy type: %s", nodeType)
	}

	return parseNodeWithReflection(nodeYAML, collection, nodeType, typeInfo)
}

func parseNodeWithReflection(nodeYAML *[]byte, collection *nodepool.Collection, nodeType string, typeInfo nodeTypeInfo) error {
	configArrayValue := reflect.New(typeInfo.configArrayType).Elem()
	configArrayPtr := configArrayValue.Addr().Interface()
	if err := yaml.Unmarshal(*nodeYAML, configArrayPtr); err != nil {
		fmt.Printf("nodeYAML: %s\n", string(*nodeYAML))
		return fmt.Errorf("failed to parse %s config: %v", nodeType, err)
	}

	configValue := configArrayValue.Index(0)

	nodeValue := reflect.New(typeInfo.nodeType).Elem()

	if typeInfo.infoFieldIndex >= 0 {
		infoField := nodeValue.Field(typeInfo.infoFieldIndex)
		if infoField.CanSet() {
			infoField.Set(emptyInfo)
		}
	}

	if typeInfo.configFieldIndex >= 0 {
		configField := nodeValue.Field(typeInfo.configFieldIndex)
		if configField.CanSet() {
			configField.Set(configValue)
		}
	}

	collectionValue := reflect.ValueOf(collection).Elem()
	fieldValue := collectionValue.Field(typeInfo.fieldIndex)

	newSlice := reflect.Append(fieldValue, nodeValue)
	fieldValue.Set(newSlice)

	return nil
}
