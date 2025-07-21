package op

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/notify"
)

var nr interfaces.NotifyRepository
var ntr interfaces.NotifyTemplateRepository
var notifyList *[]notify.Data
var notifyTemplateList map[string]string

func notifyRepo() interfaces.NotifyRepository {
	if nr == nil {
		nr = repo.Notify()
	}
	return nr
}
func GetNotifyByID(id uint16) (*notify.Data, error) {
	if notifyList == nil {
		err := refreshNotify(context.Background())
		if err != nil {
			return nil, err
		}
	}
	for _, n := range *notifyList {
		if n.ID == id {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("通知渠道不存在")
}

func GetNotifyList(ctx context.Context) ([]notify.Data, error) {
	if notifyList == nil {
		err := refreshNotify(ctx)
		if err != nil {
			return nil, err
		}
	}
	return *notifyList, nil
}
func CreateNotify(ctx context.Context, req *notify.Data) error {
	if notifyList == nil {
		err := refreshNotify(ctx)
		if err != nil {
			return err
		}
	}
	err := notifyRepo().Create(ctx, req)
	if err != nil {
		return err
	}
	*notifyList = append(*notifyList, *req)
	return nil
}
func UpdateNotify(ctx context.Context, notify *notify.Data) error {
	if notifyList == nil {
		err := refreshNotify(ctx)
		if err != nil {
			return err
		}
	}
	for i, n := range *notifyList {
		if n.ID == notify.ID {
			(*notifyList)[i] = *notify
			return notifyRepo().Update(ctx, notify)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func DeleteNotify(ctx context.Context, id uint16) error {
	if notifyList == nil {
		err := refreshNotify(ctx)
		if err != nil {
			return err
		}
	}
	for i, n := range *notifyList {
		if n.ID == id {
			*notifyList = append((*notifyList)[:i], (*notifyList)[i+1:]...)
			return notifyRepo().Delete(ctx, id)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func refreshNotify(ctx context.Context) error {
	var err error
	notifyList, err = notifyRepo().List(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NotifyTemplateRepo() interfaces.NotifyTemplateRepository {
	if ntr == nil {
		ntr = repo.NotifyTemplate()
	}
	return ntr
}
func GetNotifyTemplate(ctx context.Context, t uint16) (string, error) {
	if notifyTemplateList == nil {
		err := refreshNotifyTemplate(ctx)
		if err != nil {
			return "", err
		}
	}
	return notifyTemplateList[notify.TypeMap[t]], nil
}
func GetNotifyTemplateList(ctx context.Context) ([]notify.Template, error) {
	if notifyTemplateList == nil {
		err := refreshNotifyTemplate(ctx)
		if err != nil {
			return nil, err
		}
	}

	var list []notify.Template
	for k, v := range notifyTemplateList {
		list = append(list, notify.Template{
			Type:     k,
			Template: v,
		})
	}
	return list, nil
}
func UpdateNotifyTemplate(ctx context.Context, notify *notify.Template) error {
	if notifyTemplateList == nil {
		err := refreshNotifyTemplate(ctx)
		if err != nil {
			return err
		}
	}
	for k := range notifyTemplateList {
		if k == notify.Type {
			notifyTemplateList[notify.Type] = notify.Template
			return NotifyTemplateRepo().Update(ctx, notify)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func refreshNotifyTemplate(ctx context.Context) error {
	notifyTemplateList = make(map[string]string)
	list, err := NotifyTemplateRepo().List(ctx)
	if err != nil {
		return err
	}
	for _, n := range *list {
		notifyTemplateList[n.Type] = n.Template
	}
	return nil
}
