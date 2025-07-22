package exec

import (
	"github.com/bestruirui/bestsub/internal/models/exec"
	_ "github.com/bestruirui/bestsub/internal/modules/exec/executor"
	"github.com/bestruirui/bestsub/internal/modules/register"
)

func Get(m string, c string) (exec.Instance, error) {
	return register.Get[exec.Instance]("exec", m, c)
}

func GetTypes() []string {
	return register.GetList("exec")
}

func GetInfoMap() map[string][]register.Desc {
	return register.GetInfoMap("exec")
}
