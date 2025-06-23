package main

import (
	"context"
	"flag"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils"
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

	utils.PrintBanner(version, commit, date, builtBy)

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

	shutdown.Register(func(ctx context.Context) error {
		log.Debug("正在关闭数据库连接...")
		database.Close()
		log.Debug("数据库连接已关闭")
		return nil
	})

	shutdown.Register(func(ctx context.Context) error {
		log.Debug("正在关闭日志系统...")
		log.Close()
		log.Debug("日志系统已关闭")
		return nil
	})

	log.Info("应用程序启动完成，按 Ctrl+C 退出")
	shutdown.Listen()
}
