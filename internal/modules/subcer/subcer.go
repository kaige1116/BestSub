package subcer

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var (
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
)

func Start() error {
	cfg := config.Base()

	ctx, cancel = context.WithCancel(context.Background())
	binPath := filepath.Join(cfg.SubConverter.Path, "subconverter")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	cmd = exec.CommandContext(ctx, binPath)
	cmd.Dir = filepath.Dir(binPath)

	if err := cmd.Start(); err != nil {
		cancel()
		log.Warnf("failed to start subconverter process: %v", err)
		return err
	}

	log.Info("subconverter service started")
	return nil
}

func Stop() error {
	if cancel != nil {
		cancel()
	}

	if cmd != nil && cmd.Process != nil {
		cmd.Wait()
	}

	log.Info("subconverter service stopped")
	return nil
}

func GetBaseUrl() string {
	return fmt.Sprintf("http://127.0.0.1:%d", config.Base().SubConverter.Port)
}
