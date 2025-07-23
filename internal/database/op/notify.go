package op

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/utils/cache"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var nr interfaces.NotifyRepository
var ntr interfaces.NotifyTemplateRepository
var notifyTemplateCache = cache.New[string, string](4)
var notifyCache = cache.New[uint16, notify.Data](4)

func notifyRepo() interfaces.NotifyRepository {
	if nr == nil {
		nr = repo.Notify()
	}
	return nr
}
func GetNotifyList() ([]notify.Data, error) {
	notifyList := notifyCache.GetAll()
	if len(notifyList) == 0 {
		err := refreshNotifyCache(context.Background())
		if err != nil {
			return nil, err
		}
		notifyList = notifyCache.GetAll()
	}
	var result = make([]notify.Data, 0, len(notifyList))
	for _, v := range notifyList {
		result = append(result, v)
	}
	return result, nil
}
func GetNotifyByID(id uint16) (notify.Data, error) {
	if value, ok := notifyCache.Get(id); ok {
		return value, nil
	}
	err := refreshNotifyCache(context.Background())
	if err != nil {
		return notify.Data{}, err
	}
	if value, ok := notifyCache.Get(id); ok {
		return value, nil
	}
	return notify.Data{}, fmt.Errorf("notify not found")
}
func UpdateNotify(ctx context.Context, n *notify.Data) error {
	if notifyCache.Len() == 0 {
		err := refreshNotifyCache(context.Background())
		if err != nil {
			return err
		}
	}
	notifyCache.Set(n.ID, *n)
	return notifyRepo().Update(ctx, n)
}
func CreateNotify(ctx context.Context, n *notify.Data) error {
	if notifyCache.Len() == 0 {
		err := refreshNotifyCache(context.Background())
		if err != nil {
			return err
		}
	}
	notifyCache.Set(n.ID, *n)
	return notifyRepo().Create(ctx, n)
}
func DeleteNotify(ctx context.Context, id uint16) error {
	if notifyCache.Len() == 0 {
		err := refreshNotifyCache(context.Background())
		if err != nil {
			return err
		}
	}
	notifyCache.Del(id)
	return notifyRepo().Delete(ctx, id)
}
func refreshNotifyCache(ctx context.Context) error {
	notifyList, err := notifyRepo().List(ctx)
	if err != nil {
		return err
	}
	for _, n := range *notifyList {
		notifyCache.Set(n.ID, n)
	}
	return nil
}

func NotifyTemplateRepo() interfaces.NotifyTemplateRepository {
	if ntr == nil {
		ntr = repo.NotifyTemplate()
	}
	return ntr
}
func GetNotifyTemplateList() ([]notify.Template, error) {
	notifyTemplateList := notifyTemplateCache.GetAll()
	if len(notifyTemplateList) == 0 {
		err := refreshNotifyTemplate(context.Background())
		if err != nil {
			return nil, err
		}
		notifyTemplateList = notifyTemplateCache.GetAll()
	}
	var result = make([]notify.Template, 0, len(notifyTemplateList))
	for k, v := range notifyTemplateList {
		result = append(result, notify.Template{Type: k, Template: v})
	}
	return result, nil
}

func GetNotifyTemplateByType(t string) (string, error) {
	if value, ok := notifyTemplateCache.Get(t); ok {
		return value, nil
	}
	err := refreshNotifyTemplate(context.Background())
	if err != nil {
		return "", err
	}
	if value, ok := notifyTemplateCache.Get(t); ok {
		return value, nil
	}
	return "", fmt.Errorf("notify template not found")
}
func UpdateNotifyTemplate(ctx context.Context, nt *notify.Template) error {
	if notifyTemplateCache.Len() == 0 {
		refreshNotifyTemplate(context.Background())
	}
	log.Debugf("Update Notify Template Len: %v", notifyTemplateCache.Len())
	notifyTemplateCache.Set(nt.Type, nt.Template)
	return NotifyTemplateRepo().Update(ctx, nt)
}
func refreshNotifyTemplate(ctx context.Context) error {
	notifyTemplateCache.Clear()
	notifyTemplates, err := NotifyTemplateRepo().List(ctx)
	if err != nil {
		return err
	}
	for _, t := range *notifyTemplates {
		notifyTemplateCache.Set(t.Type, t.Template)
	}
	return nil
}
