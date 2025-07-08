package main

import (
	"flag"

	"github.com/bestruirui/bestsub/internal/api/server"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils/banner"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/shutdown"
)

func main() {

	configPath := flag.String("c", config.DefaultConfigPath, "config file path")
	flag.Parse()

	banner.Print()

	if err := config.Initialize(*configPath); err != nil {
		panic(err)
	}

	cfg := config.Get()

	if err := log.Initialize(cfg.Log.Level, cfg.Log.Output, cfg.Log.Dir); err != nil {
		panic(err)
	}
	if err := database.Initialize(cfg.Database.Type, cfg.Database.Path); err != nil {
		log.Fatal(err)
	}

	if err := server.Initialize(); err != nil {
		panic(err)
	}
	if err := task.Initialize(); err != nil {
		panic(err)
	}

	task.Start()
	server.Start()

	shutdown.Register("Log", log.Close)
	shutdown.Register("Database", database.Close)
	shutdown.Register("HTTP Server", server.Close)
	shutdown.Register("Task", task.Shutdown)

	shutdown.Listen()
}
