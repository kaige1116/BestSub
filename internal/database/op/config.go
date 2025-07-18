package op

import (
	"context"
	"strconv"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/system"
)

var configRepo interfaces.ConfigRepository
var cache = fastcache.New(1024 * 1024)

func ConfigRepo() interfaces.ConfigRepository {
	if configRepo == nil {
		configRepo = repo.Config()
	}
	return configRepo
}
func GetAllConfig(ctx context.Context) ([]system.GroupData, error) {
	sysConf := system.DefaultDbConfig()
	notExistConfig := make([]string, 0)

	for i := range sysConf {
		for j := range sysConf[i].Data {
			if value, ok := cache.HasGet(nil, []byte(sysConf[i].Data[j].Key)); ok {
				sysConf[i].Data[j].Value = string(value)
			} else {
				notExistConfig = append(notExistConfig, sysConf[i].Data[j].Key)
			}
		}
	}

	if len(notExistConfig) > 0 {
		configs, err := ConfigRepo().GetByKey(ctx, notExistConfig)
		if err != nil {
			return nil, err
		}

		configMap := make(map[string]string)
		for _, config := range *configs {
			configMap[config.Key] = config.Value
			cache.Set([]byte(config.Key), []byte(config.Value))
		}

		for i := range sysConf {
			for j := range sysConf[i].Data {
				if value, exists := configMap[sysConf[i].Data[j].Key]; exists {
					sysConf[i].Data[j].Value = value
				}
			}
		}
	}

	return sysConf, nil
}
func GetConfigByKey(key string) (string, error) {
	if value, ok := cache.HasGet(nil, []byte(key)); ok {
		return string(value), nil
	}
	configs, err := ConfigRepo().GetByKey(context.Background(), []string{key})
	if err != nil {
		return "", err
	}
	cache.Set([]byte(key), []byte((*configs)[0].Value))
	return (*configs)[0].Value, nil
}
func GetConfigByGroup(groupName string) ([]system.Data, error) {
	sysConf := system.DefaultDbConfig()
	var result []system.Data
	notExistConfig := make([]string, 0)

	for _, group := range sysConf {
		if group.GroupName == groupName {
			for _, data := range group.Data {
				if value, ok := cache.HasGet(nil, []byte(data.Key)); ok {
					data.Value = string(value)
					result = append(result, data)
				} else {
					notExistConfig = append(notExistConfig, data.Key)
					result = append(result, data)
				}
			}
			break
		}
	}

	if len(notExistConfig) > 0 {
		configs, err := ConfigRepo().GetByKey(context.Background(), notExistConfig)
		if err != nil {
			return nil, err
		}

		configMap := make(map[string]string)
		for _, config := range *configs {
			configMap[config.Key] = config.Value
			cache.Set([]byte(config.Key), []byte(config.Value))
		}

		for i := range result {
			if value, exists := configMap[result[i].Key]; exists {
				result[i].Value = value
			}
		}
	}

	return result, nil
}
func SetConfig(key string, value string) error {
	cache.Set([]byte(key), []byte(value))
	err := ConfigRepo().Update(context.Background(), &[]system.UpdateData{
		{
			Key:   key,
			Value: value,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func UpdateConfig(ctx context.Context, config *[]system.UpdateData) error {
	for _, item := range *config {
		cache.Set([]byte(item.Key), []byte(item.Value))
	}
	err := ConfigRepo().Update(ctx, config)
	if err != nil {
		return err
	}
	return nil

}
func GetConfigStr(key string) string {
	value, err := GetConfigByKey(key)
	if err != nil {
		return ""
	}
	return value
}
func GetConfigInt(key string) int {
	value, err := GetConfigByKey(key)
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}
func GetConfigBool(key string) bool {
	value, err := GetConfigByKey(key)
	if err != nil {
		return false
	}
	return value == "true"
}
