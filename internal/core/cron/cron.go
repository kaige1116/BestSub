package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

type cronFunc struct {
	fn       func()
	cronExpr string
}

var scheduler = cron.New(cron.WithLocation(time.Local))

func Start() {
	scheduler.Start()
}

func Stop() {
	scheduler.Stop()
}
