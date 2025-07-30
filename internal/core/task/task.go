package task

import (
	"time"

	"github.com/robfig/cron/v3"
)

var scheduler = cron.New(cron.WithLocation(time.Local))

func Start() {
	scheduler.Start()
}

func Stop() {
	scheduler.Stop()
}
