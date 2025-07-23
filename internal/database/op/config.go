package op

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/config"
	"github.com/bestruirui/bestsub/internal/utils/cache"
)

var configRepo interfaces.ConfigRepository
var configCache = cache.New[string, string](4)

func ConfigRepo() interfaces.ConfigRepository {
	if configRepo == nil {
		configRepo = repo.Config()
	}
	return configRepo
}
func GetAllConfig(ctx context.Context) ([]config.GroupAdvance, error) {
	sysConf := config.DefaultAdvance()
	sysConfCache := configCache.GetAll()
	if len(sysConfCache) == 0 {
		err := refreshConfigCache(context.Background())
		if err != nil {
			return nil, err
		}
		sysConfCache = configCache.GetAll()
	}
	for i := range sysConf {
		for j := range sysConf[i].Data {
			value, ok := sysConfCache[sysConf[i].Data[j].Key]
			if ok {
				sysConf[i].Data[j].Value = value
			}
		}
	}

	return sysConf, nil
}
func GetConfigByKey(key string) (string, error) {
	if value, ok := configCache.Get(key); ok {
		return value, nil
	}
	err := refreshConfigCache(context.Background())
	if err != nil {
		return "", err
	}
	if value, ok := configCache.Get(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("config not found")
}
func GetConfigByGroup(groupName string) ([]config.GroupAdvance, error) {
	sysConf := config.DefaultAdvance()
	var result []config.Advance
	var description string
	sysConfCache := configCache.GetAll()
	if len(sysConfCache) == 0 {
		err := refreshConfigCache(context.Background())
		if err != nil {
			return nil, err
		}
		sysConfCache = configCache.GetAll()
	}
	for _, group := range sysConf {
		if group.GroupName == groupName {
			for _, data := range group.Data {
				value, ok := sysConfCache[data.Key]
				if ok {
					result = append(result, config.Advance{
						Key:   data.Key,
						Value: value,
					})
				}
			}
			description = group.Description
		}
	}
	return []config.GroupAdvance{
		{
			GroupName:   groupName,
			Description: description,
			Data:        result,
		},
	}, nil
}
func UpdateConfig(ctx context.Context, config *[]config.UpdateAdvance) error {
	if configCache.Len() == 0 {
		err := refreshConfigCache(context.Background())
		if err != nil {
			return err
		}
	}
	for _, item := range *config {
		configCache.Set(item.Key, item.Value)
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

func refreshConfigCache(ctx context.Context) error {
	configs, err := ConfigRepo().GetAll(ctx)
	if err != nil {
		return err
	}
	for _, config := range *configs {
		configCache.Set(config.Key, config.Value)
	}
	return nil
}
