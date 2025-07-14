package exec

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

// 全局处理器注册表
var (
	execs     = make(map[string]RegisterInfo)
	execInfos = make(map[string]Info)
	mu        sync.RWMutex
)

// 前端数据类型常量
const (
	TypeBoolean = "boolean"
	TypeNumber  = "number"
	TypeString  = "string"
)

// RegisterInfo 处理器注册信息
type RegisterInfo struct {
	Type    string
	Handler Execer
	Config  any
}

// Info 前端响应结构
type Info struct {
	Type   string   `json:"type"`
	Config []Config `json:"config"`
}

// Config 配置字段信息
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     any    `json:"default"`
}

// TaskInfo 任务执行上下文信息
type TaskInfo struct {
	Type   string
	ID     int64
	Level  string
	Config []byte
}

type Execer interface {
	Do(ctx context.Context, logger *log.Logger, task *TaskInfo) error
}

// Register 注册处理器
func Register(info *RegisterInfo) {
	if info == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := execs[info.Type]; exists {
		log.Fatalf("任务类型处理器已存在: %s", info.Type)
		return
	}
	execs[info.Type] = *info

	execInfos[info.Type] = Info{
		Type:   info.Type,
		Config: parseFields(info.Config),
	}
	log.Debugf("注册任务类型处理器: %s", info.Type)
}

// Get 获取指定类型的处理器
func Get(taskType string) Execer {
	return withReadLock(func() Execer {
		info, exists := execs[taskType]
		if !exists {
			return nil
		}
		return info.Handler
	})
}

// GetAll 获取所有注册的处理器信息
func GetAll() []Info {
	return withReadLock(func() []Info {
		infos := make([]Info, 0, len(execInfos))
		for _, info := range execInfos {
			infos = append(infos, info)
		}
		return infos
	})
}

// GetTypes 获取所有已注册的任务类型
func GetTypes() []string {
	return withReadLock(func() []string {
		types := make([]string, 0, len(execs))
		for taskType := range execs {
			types = append(types, taskType)
		}
		return types
	})
}

// GetExecInfo 获取执行器信息
func GetExecInfo(execType string) Info {
	return withReadLock(func() Info {
		info, exists := execInfos[execType]
		if !exists {
			return Info{}
		}
		return info
	})
}

// Run 执行任务
func Run(ctx context.Context, task *TaskInfo) error {
	execer := Get(task.Type)
	if execer == nil {
		return fmt.Errorf("未找到任务类型处理器: %s", task.Type)
	}
	logger, err := log.NewTaskLogger(task.ID, task.Level)
	if err != nil {
		return err
	}
	return execer.Do(ctx, logger, task)
}

// withReadLock 执行需要读锁保护的操作
func withReadLock[T any](fn func() T) T {
	mu.RLock()
	defer mu.RUnlock()
	return fn()
}

// parseFields 解析配置结构体字段信息
func parseFields(config any) []Config {
	var fields []Config

	configType, configValue := getStructTypeAndValue(config)
	if configType == nil {
		return fields
	}

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldInfo := Config{
			Name: getFieldName(field),
			Type: getFieldType(field.Type),
		}

		if desc := field.Tag.Get("description"); desc != "" {
			fieldInfo.Description = desc
		}

		if required := field.Tag.Get("required"); required == "true" {
			fieldInfo.Required = true
		}

		if defaultValue := field.Tag.Get("default"); defaultValue != "" {
			fieldInfo.Default = parseDefaultValue(defaultValue, field.Type)
		} else if fieldValue.IsValid() && !fieldValue.IsZero() {
			fieldInfo.Default = fieldValue.Interface()
		}

		fields = append(fields, fieldInfo)
	}

	return fields
}

// getStructTypeAndValue 获取结构体的类型和值，处理指针解引用
func getStructTypeAndValue(config any) (reflect.Type, reflect.Value) {
	configType := reflect.TypeOf(config)
	configValue := reflect.ValueOf(config)

	if configType.Kind() == reflect.Ptr {
		configType = configType.Elem()
		configValue = configValue.Elem()
	}

	if configType.Kind() != reflect.Struct {
		return nil, reflect.Value{}
	}

	return configType, configValue
}
func getFieldName(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}
	return strings.ToLower(field.Name)
}

func getFieldType(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.Bool:
		return TypeBoolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return TypeNumber
	case reflect.String:
		return TypeString
	default:
		return TypeString
	}
}

func parseDefaultValue(value string, fieldType reflect.Type) any {
	switch fieldType.Kind() {
	case reflect.Bool:
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(value, 10, 64); err == nil {
			return u
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	case reflect.String:
		return value
	}
	return value
}
