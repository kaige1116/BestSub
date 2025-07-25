package main

import (
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/server/auth"
	"github.com/bestruirui/bestsub/internal/server/server"
	"github.com/bestruirui/bestsub/internal/utils/info"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/shutdown"
)

func main() {

	info.Banner()

	cfg := config.Base()

	if err := log.Initialize(cfg.Log.Level, cfg.Log.Path, cfg.Log.Output); err != nil {
		panic(err)
	}
	if err := database.Initialize(cfg.Database.Type, cfg.Database.Path); err != nil {
		panic(err)
	}

	if err := server.Initialize(); err != nil {
		panic(err)
	}

	task.Check.Start()
	task.Fetch.Start()
	task.Load()

	server.Start()
	log.CleanupOldLogs(5)

	shutdown.Register(server.Close)      // 关闭顺序
	shutdown.Register(task.Check.Stop)   //   ↓↓
	shutdown.Register(task.Fetch.Stop)   //   ↓↓
	shutdown.Register(database.Close)    //   ↓↓
	shutdown.Register(auth.CloseSession) //   ↓↓
	shutdown.Register(log.Close)         //   ↓↓

	shutdown.Listen()
}
