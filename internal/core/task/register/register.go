package register

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// 前端数据类型常量
const (
	typeBoolean = "boolean"
	typeNumber  = "number"
	typeString  = "string"
)

// Handler 任务处理器接口
type Handler interface {
	Execute(ctx context.Context, config string) error
	Validate(config string) error
}

// Info 处理器注册信息
type Info struct {
	Type    string  // 任务类型
	Handler Handler // 处理器实例
	Config  any     // 配置项结构示例
}

// Field 配置字段信息
type Field struct {
	Name        string `json:"name"`        // 字段名
	Type        string `json:"type"`        // 字段类型
	Description string `json:"description"` // 字段说明
	Required    bool   `json:"required"`    // 是否必填
	Default     any    `json:"default"`     // 默认值
}

// Response 前端响应结构
type Response struct {
	Type   string  `json:"type"`   // 任务类型
	Config []Field `json:"config"` // 配置项字段列表
}

// 全局处理器注册表
var (
	handlers = make(map[string]Info)
	mu       sync.RWMutex
)

// Add 注册处理器
func Add(info *Info) {
	if info == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	handlers[info.Type] = *info
}

// Get 获取指定类型的处理器
func Get(taskType string) (Handler, bool) {
	mu.RLock()
	defer mu.RUnlock()

	info, exists := handlers[taskType]
	if !exists {
		return nil, false
	}
	return info.Handler, true
}

// GetAll 获取所有注册的处理器信息，返回可直接发给前端的JSON结构
func GetAll() []Response {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Response, 0, len(handlers))
	for _, info := range handlers {
		var configFields []Field

		if info.Config != nil {
			configFields = parseFields(info.Config)
		}

		result = append(result, Response{
			Type:   info.Type,
			Config: configFields,
		})
	}
	return result
}

// GetTypes 获取所有已注册的任务类型
func GetTypes() []string {
	mu.RLock()
	defer mu.RUnlock()

	types := make([]string, 0, len(handlers))
	for taskType := range handlers {
		types = append(types, taskType)
	}
	return types
}

// GetConfig 根据任务类型获取配置字段信息
func GetConfig(taskType string) []Field {
	mu.RLock()
	defer mu.RUnlock()

	info, exists := handlers[taskType]
	if !exists {
		return nil
	}

	if info.Config == nil {
		return nil
	}

	return parseFields(info.Config)
}

// parseFields 解析配置结构体字段信息
func parseFields(config any) []Field {
	var fields []Field

	configType := reflect.TypeOf(config)
	configValue := reflect.ValueOf(config)

	// 如果是指针，获取其指向的类型和值
	if configType.Kind() == reflect.Ptr {
		configType = configType.Elem()
		configValue = configValue.Elem()
	}

	// 只处理结构体类型
	if configType.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 解析json标签，获取字段名（去掉选项部分）
		jsonName, _, _ := strings.Cut(jsonTag, ",")

		configField := Field{
			Name:        jsonName,
			Type:        getFieldType(field.Type),
			Description: field.Tag.Get("description"),
			Required:    field.Tag.Get("required") == "true",
			Default:     getDefaultValueFromTag(field, fieldValue),
		}

		fields = append(fields, configField)
	}

	return fields
}

// getFieldType 获取前端数据类型字符串
func getFieldType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return typeBoolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return typeNumber
	case reflect.String:
		return typeString
	default:
		return typeString
	}
}

// getDefaultValueFromTag 从标签或字段值获取默认值
func getDefaultValueFromTag(field reflect.StructField, fieldValue reflect.Value) any {
	// 首先尝试从 default 标签获取
	defaultTag := field.Tag.Get("default")
	if defaultTag != "" {
		return parseDefaultValue(defaultTag, field.Type)
	}

	// 如果没有 default 标签，返回字段的零值
	return getDefaultValue(fieldValue)
}

// parseDefaultValue 解析默认值字符串
func parseDefaultValue(defaultStr string, fieldType reflect.Type) any {
	switch fieldType.Kind() {
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultStr); err == nil {
			return val
		}
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(defaultStr, 10, 64); err == nil {
			return val
		}
		return int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(defaultStr, 10, 64); err == nil {
			return val
		}
		return uint64(0)
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(defaultStr, 64); err == nil {
			return val
		}
		return float64(0)
	case reflect.String:
		return defaultStr
	default:
		return defaultStr
	}
}

// getDefaultValue 获取字段的默认值
func getDefaultValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		return v.String()
	default:
		return v.Interface()
	}
}
