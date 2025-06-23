package start

import (
	"flag"
	"fmt"
	"os"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func InitConfig() error {
	configPath := flag.String("f", "", "Configuration file path, default is ./bestsub/data/config.json")
	flag.Parse()
	config, err := config.LoadConfigFromPath(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}
	if err := log.Init(config); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}
	return nil
}
