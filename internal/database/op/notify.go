package op

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/modules/register"
)

var nr interfaces.NotifyRepository
var ntr interfaces.NotifyTemplateRepository
var notifyList *[]notify.Data
var notifyTemplateList *[]notify.Template

func notifyRepo() interfaces.NotifyRepository {
	if nr == nil {
		nr = repo.Notify()
	}
	return nr
}
func GetNotifyTypes() []string {
	return register.GetNotifyTypes()
}
func GetNotifyByType(t string) (*notify.Data, error) {
	if notifyList == nil {
		err := refreshNotify(context.Background())
		if err != nil {
			return nil, err
		}
	}
	for _, n := range *notifyList {
		if n.Type == t {
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
	refreshNotify(ctx)
	for _, n := range *notifyList {
		if n.Type == req.Type {
			return fmt.Errorf("通知渠道已存在")
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
	refreshNotify(ctx)
	for i, n := range *notifyList {
		if n.ID == notify.ID {
			(*notifyList)[i] = *notify
			return notifyRepo().Update(ctx, notify)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func DeleteNotify(ctx context.Context, id uint16) error {
	refreshNotify(ctx)
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

func notifyTemplateRepo() interfaces.NotifyTemplateRepository {
	if ntr == nil {
		ntr = repo.NotifyTemplate()
	}
	return ntr
}
func GetNotifyTemplateList(ctx context.Context) ([]notify.Template, error) {
	if notifyTemplateList == nil {
		err := refreshNotifyTemplate(ctx)
		if err != nil {
			return nil, err
		}
	}
	return *notifyTemplateList, nil
}
func CreateNotifyTemplate(ctx context.Context, req *notify.Template) error {
	refreshNotifyTemplate(ctx)
	for _, n := range *notifyTemplateList {
		if n.Name == req.Name {
			return fmt.Errorf("通知渠道已存在")
		}
	}
	err := notifyTemplateRepo().Create(ctx, req)
	if err != nil {
		return err
	}
	*notifyTemplateList = append(*notifyTemplateList, *req)
	return nil
}
func UpdateNotifyTemplate(ctx context.Context, notify *notify.Template) error {
	refreshNotifyTemplate(ctx)
	for i, n := range *notifyTemplateList {
		if n.ID == notify.ID {
			(*notifyTemplateList)[i] = *notify
			return notifyTemplateRepo().Update(ctx, notify)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func DeleteNotifyTemplate(ctx context.Context, id uint16) error {
	refreshNotifyTemplate(ctx)
	for i, n := range *notifyTemplateList {
		if n.ID == id {
			*notifyTemplateList = append((*notifyTemplateList)[:i], (*notifyTemplateList)[i+1:]...)
			return notifyTemplateRepo().Delete(ctx, id)
		}
	}
	return fmt.Errorf("通知渠道不存在")
}
func refreshNotifyTemplate(ctx context.Context) error {
	var err error
	notifyTemplateList, err = notifyTemplateRepo().List(ctx)
	if err != nil {
		return err
	}
	return nil
}
