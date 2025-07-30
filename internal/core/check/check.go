package check

import (
	_ "github.com/bestruirui/bestsub/internal/core/check/checker"
	"github.com/bestruirui/bestsub/internal/models/check"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/desc"
)

type Desc = desc.Data

func Get(m string, c string) (check.Instance, error) {
	return register.Get[check.Instance]("check", m, c)
}

func GetTypes() []string {
	return register.GetList("check")
}

func GetInfoMap() map[string][]Desc {
	return register.GetInfoMap("check")
}
