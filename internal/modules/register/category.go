package register

import (
	"github.com/bestruirui/bestsub/internal/models/exec"
	"github.com/bestruirui/bestsub/internal/models/notify"
)

func Notify(i notify.Instance) {
	register("notify", i)
}
func Exec(i exec.Instance) {
	register("exec", i)
}
