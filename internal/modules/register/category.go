package register

import (
	"github.com/bestruirui/bestsub/internal/models/check"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/models/storage"
)

func Notify(i notify.Instance) {
	register("notify", i)
}
func Check(i check.Instance) {
	register("check", i)
}
func Storage(i storage.Instance) {
	register("storage", i)
}
