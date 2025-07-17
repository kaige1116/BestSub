package database

import (
	"context"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/models/task"
	"golang.org/x/crypto/bcrypt"
)

type memberCache struct {
	auth              *auth.Data
	SystemConfig      *system.Data
	Notify            *notify.Data
	NotifyTemplate    *notify.Template
	Task              *task.Data
	SubStorageConfig  *sub.StorageConfig
	SubOutputTemplate *sub.OutputTemplate
	SubNodeFilterRule *sub.NodeFilterRule
	SubLink           *sub.Data
}

type repositoryCache struct {
	auth              interfaces.AuthRepository
	systemConfig      interfaces.ConfigRepository
	notify            interfaces.NotifyRepository
	notifyTemplate    interfaces.NotifyTemplateRepository
	task              interfaces.TaskRepository
	sub               interfaces.SubRepository
	subStorageConfig  interfaces.SubStorageConfigRepository
	subSaveConfig     interfaces.SubSaveRepository
	subShare          interfaces.SubShareRepository
	subOutputTemplate interfaces.SubOutputTemplateRepository
	subNodeFilterRule interfaces.SubNodeFilterRuleRepository
}

var member memberCache
var repository repositoryCache

func AuthRepo() interfaces.AuthRepository {
	if repository.auth == nil {
		repository.auth = repo.Auth()
	}
	return repository.auth
}
func ConfigRepo() interfaces.ConfigRepository {
	if repository.systemConfig == nil {
		repository.systemConfig = repo.Config()
	}
	return repository.systemConfig
}
func NotifyRepo() interfaces.NotifyRepository {
	if repository.notify == nil {
		repository.notify = repo.Notify()
	}
	return repository.notify
}
func NotifyTemplateRepo() interfaces.NotifyTemplateRepository {
	if repository.notifyTemplate == nil {
		repository.notifyTemplate = repo.NotifyTemplate()
	}
	return repository.notifyTemplate
}
func TaskRepo() interfaces.TaskRepository {
	if repository.task == nil {
		repository.task = repo.Task()
	}
	return repository.task
}
func SubStorageConfigRepo() interfaces.SubStorageConfigRepository {
	if repository.subStorageConfig == nil {
		repository.subStorageConfig = repo.SubStorageConfig()
	}
	return repository.subStorageConfig
}
func SubOutputTemplateRepo() interfaces.SubOutputTemplateRepository {
	if repository.subOutputTemplate == nil {
		repository.subOutputTemplate = repo.SubOutputTemplate()
	}
	return repository.subOutputTemplate
}
func SubNodeFilterRuleRepo() interfaces.SubNodeFilterRuleRepository {
	if repository.subNodeFilterRule == nil {
		repository.subNodeFilterRule = repo.SubNodeFilterRule()
	}
	return repository.subNodeFilterRule
}
func SubRepo() interfaces.SubRepository {
	if repository.sub == nil {
		repository.sub = repo.Sub()
	}
	return repository.sub
}
func SubSaveConfigRepo() interfaces.SubSaveRepository {
	if repository.subSaveConfig == nil {
		repository.subSaveConfig = repo.SubSaveConfig()
	}
	return repository.subSaveConfig
}
func SubShareRepo() interfaces.SubShareRepository {
	if repository.subShare == nil {
		repository.subShare = repo.SubShareLink()
	}
	return repository.subShare
}

func AuthGet() (*auth.Data, error) {
	var err error
	if member.auth == nil {
		member.auth, err = AuthRepo().Get(context.Background())
	}
	return member.auth, err
}
func AuthUpdateName(name string) error {
	auth, err := AuthGet()
	if err != nil {
		return err
	}
	auth.UserName = name
	updatedAt, err := AuthRepo().UpdateName(context.Background(), name)
	if err != nil {
		return err
	}
	auth.UpdatedAt = updatedAt
	return nil
}
func AuthUpdatePassWord(password string) error {
	auth, err := AuthGet()
	if err != nil {
		return err
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	auth.Password = string(hashedBytes)
	updatedAt, err := AuthRepo().UpdatePassword(context.Background(), auth.Password)
	if err != nil {
		return err
	}
	auth.UpdatedAt = updatedAt
	return nil
}
func AuthVerify(username, password string) error {
	auth, err := AuthGet()
	if err != nil {
		return err
	}
	if auth.UserName != username {
		return fmt.Errorf("用户名不匹配")
	}
	return bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password))
}
