package main

import (
	"flag"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils/banner"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/shutdown"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {

	configPath := flag.String("c", config.DefaultConfigPath, "config file path")
	flag.Parse()

	banner.Print(version, commit, date, builtBy)

	if err := config.Initialize(*configPath); err != nil {
		panic(err)
	}

	cfg := config.GetConfig()

	if err := log.Initialize(cfg.Log); err != nil {
		panic(err)
	}
	if err := database.Initialize(cfg.Database); err != nil {
		log.Fatal(err)
	}
	shutdown.Register("Log", log.Close)
	shutdown.Register("Database", database.Close)

	shutdown.Listen()
}
