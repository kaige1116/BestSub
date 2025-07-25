package storage

import (
	storageModel "github.com/bestruirui/bestsub/internal/models/storage"
	"github.com/bestruirui/bestsub/internal/modules/register"
	_ "github.com/bestruirui/bestsub/internal/modules/storage/channel"
	"github.com/bestruirui/bestsub/internal/utils/desc"
)

type Desc = desc.Data

func Get(m string, c string) (storageModel.Instance, error) {
	return register.Get[storageModel.Instance]("storage", m, c)
}

func GetChannels() []string {
	return register.GetList("storage")
}

func GetInfoMap() map[string][]Desc {
	return register.GetInfoMap("storage")
}
