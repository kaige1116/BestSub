package register

import (
	"github.com/bestruirui/bestsub/internal/models/exec"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/models/storage"
)

func Notify(i notify.Instance) {
	register("notify", i)
}
func Exec(i exec.Instance) {
	register("exec", i)
}
func Storage(i storage.Instance) {
	register("storage", i)
}
