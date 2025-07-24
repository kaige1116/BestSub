package exec

import (
	"github.com/bestruirui/bestsub/internal/models/exec"
	_ "github.com/bestruirui/bestsub/internal/modules/exec/executor"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/desc"
)

type Desc = desc.Data

func Get(m string, c string) (exec.Instance, error) {
	return register.Get[exec.Instance]("exec", m, c)
}

func GetTypes() []string {
	return register.GetList("exec")
}

func GetInfoMap() map[string][]Desc {
	return register.GetInfoMap("exec")
}
