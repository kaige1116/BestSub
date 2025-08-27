package subcer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var (
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
)

func init() {
	prefPath := filepath.Join(config.Base().SubConverter.Path, "pref.yml")
	if _, err := os.Stat(prefPath); err != nil {
		os.MkdirAll(config.Base().SubConverter.Path, 0755)
		cfg := fmt.Sprintf(pref, config.Base().SubConverter.Host, config.Base().SubConverter.Port)
		if err := os.WriteFile(prefPath, []byte(cfg), 0644); err != nil {
			log.Errorf("failed to write subconverter config: %v", err)
			os.Exit(1)
		}
	}
}

func Start() error {
	cfg := config.Base()

	ctx, cancel = context.WithCancel(context.Background())
	binPath := filepath.Join(cfg.SubConverter.Path, "subconverter")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	cmd = exec.CommandContext(ctx, binPath)
	cmd.Dir = filepath.Dir(binPath)
	if utils.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
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

	log.Debug("subconverter service stopped")
	return nil
}

func Lock() {
	mu.Lock()
}
func RLock() {
	mu.RLock()
}

func RUnlock() {
	mu.RUnlock()
}

func Unlock() {
	mu.Unlock()
}

func GetBaseUrl() string {
	return fmt.Sprintf("http://127.0.0.1:%d", config.Base().SubConverter.Port)
}
func GetVersion() string {
	resp, err := http.Get(fmt.Sprintf("%s/version", GetBaseUrl()))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	parts := strings.Split(string(body), " ")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
