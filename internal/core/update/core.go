package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/shutdown"
)

func UpdateCore() error {
	log.Infof("start update core")
	err := updateCore()
	if err != nil {
		log.Warnf("update core failed, please update manually", err)
		return err
	}
	log.Infof("update core success")
	return nil
}

func updateCore() error {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	var downloadUrl string

	var filename string
	switch goos {
	case "windows":
		switch arch {
		case "386":
			filename = "bestsub-windows-x86.zip"
		case "amd64":
			filename = "bestsub-windows-x86_64.zip"
		default:
			log.Errorf("unsupported windows architecture: %s", arch)
			return fmt.Errorf("unsupported windows architecture: %s", arch)
		}
	case "darwin":
		switch arch {
		case "amd64":
			filename = "bestsub-darwin-x86_64.zip"
		case "arm64":
			filename = "bestsub-darwin-arm64.zip"
		default:
			log.Errorf("unsupported darwin architecture: %s", arch)
			return fmt.Errorf("unsupported darwin architecture: %s", arch)
		}
	case "linux":
		switch arch {
		case "386":
			filename = "bestsub-linux-x86.zip"
		case "amd64":
			filename = "bestsub-linux-x86_64.zip"
		case "arm":
			filename = "bestsub-linux-armv7.zip"
		case "arm64":
			filename = "bestsub-linux-arm64.zip"
		default:
			log.Errorf("unsupported linux architecture: %s", arch)
			return fmt.Errorf("unsupported linux architecture: %s", arch)
		}
	default:
		log.Errorf("unsupported operating system: %s", goos)
		return fmt.Errorf("unsupported operating system: %s", goos)
	}

	downloadUrl = bestsubUpdateUrl + "/" + filename

	bytes, err := download(downloadUrl, op.GetSettingBool(setting.PROXY_URL))
	if err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	execDir := filepath.Dir(execPath)
	if err := unzip(bytes, execDir); err != nil {
		return err
	}
	go restartExecutable(execPath)
	return nil
}

func restartExecutable(execPath string) {
	var err error
	shutdown.All()
	if runtime.GOOS == "windows" {
		cmd := exec.Command(execPath, os.Args[1:]...)
		log.Infof("restarting: %q %q", execPath, os.Args[1:])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			log.Errorf("restarting: %s", err)
		}

		os.Exit(0)
	}

	log.Infof("restarting: %q %q", execPath, os.Args[1:])
	err = syscall.Exec(execPath, os.Args, os.Environ())
	if err != nil {
		log.Errorf("restarting: %s", err)
	}
	log.Infof("restarting success")
}
