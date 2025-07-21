package register

import (
	"github.com/bestruirui/bestsub/internal/models/exec"
	"github.com/bestruirui/bestsub/internal/models/notify"
)

func Notify(m string, i notify.Instance) {
	register("notify", m, i)
}
func Exec(m string, i exec.Instance) {
	register("exec", m, i)
}
