package op

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	subModel "github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/cache"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var subRepo interfaces.SubRepository
var subCache = cache.New[uint16, subModel.Data](16)

func SubRepo() interfaces.SubRepository {
	if subRepo == nil {
		subRepo = repo.Sub()
	}
	return subRepo
}
func GetSubList(ctx context.Context) ([]subModel.Data, error) {
	subList := subCache.GetAll()
	if len(subList) == 0 {
		err := refreshSubCache(context.Background())
		if err != nil {
			return nil, err
		}
		subList = subCache.GetAll()
	}
	var result = make([]subModel.Data, 0, len(subList))
	for _, v := range subList {
		result = append(result, v)
	}
	return result, nil
}

func GetSubByID(ctx context.Context, id uint16) (*subModel.Data, error) {
	if subCache.Len() == 0 {
		if err := refreshSubCache(ctx); err != nil {
			return nil, err
		}
	}
	if s, ok := subCache.Get(id); ok {
		return &s, nil
	}
	return nil, fmt.Errorf("sub not found")
}
func CreateSub(ctx context.Context, sub *subModel.Data) error {
	if subCache.Len() == 0 {
		if err := refreshSubCache(ctx); err != nil {
			return err
		}
	}
	if err := SubRepo().Create(ctx, sub); err != nil {
		return err
	}
	subCache.Set(sub.ID, *sub)
	return nil
}
func UpdateSub(ctx context.Context, sub *subModel.Data) error {
	if subCache.Len() == 0 {
		if err := refreshSubCache(ctx); err != nil {
			return err
		}
	}
	oldSub, ok := subCache.Get(sub.ID)
	if !ok {
		return fmt.Errorf("sub not found")
	}
	sub.Result = oldSub.Result
	sub.CreatedAt = oldSub.CreatedAt
	if err := SubRepo().Update(ctx, sub); err != nil {
		return err
	}
	subCache.Set(sub.ID, *sub)
	return nil
}
func UpdateSubResult(ctx context.Context, id uint16, result subModel.Result) error {
	if subCache.Len() == 0 {
		if err := refreshSubCache(ctx); err != nil {
			return err
		}
	}
	sub, ok := subCache.Get(id)
	if !ok {
		return fmt.Errorf("sub not found")
	}
	var oldStatus subModel.Result
	json.Unmarshal([]byte(sub.Result), &oldStatus)

	result.Success += oldStatus.Success
	result.Fail += oldStatus.Fail

	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	sub.Result = string(bytes)
	if err := SubRepo().Update(ctx, &sub); err != nil {
		return err
	}
	subCache.Set(id, sub)
	log.Infof("sub %d result updated", id)
	return nil
}
func DeleteSub(ctx context.Context, id uint16) error {
	if subCache.Len() == 0 {
		if err := refreshSubCache(ctx); err != nil {
			return err
		}
	}
	if err := SubRepo().Delete(ctx, id); err != nil {
		return err
	}
	subCache.Del(id)
	return nil
}
func refreshSubCache(ctx context.Context) error {
	subList, err := SubRepo().List(ctx)
	if err != nil {
		return err
	}
	for _, s := range *subList {
		subCache.Set(s.ID, s)
	}
	return nil
}
