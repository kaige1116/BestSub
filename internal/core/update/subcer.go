package update

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func InitSubconverter() error {
	filePath := config.Base().SubConverter.Path + "/subconverter"
	if runtime.GOOS == "windows" {
		filePath += ".exe"
	}

	if _, err := os.Stat(filePath); err != nil {
		log.Infof("subconverter not found, downloading...")
		err = UpdateSubconverter()
		if err != nil {
			log.Warnf("auto update subconverter failed, please download subconverter manually from %s and move to %s: %v", op.GetSettingStr(setting.SUBCONVERTER_URL), config.Base().SubConverter.Path, err)
			os.Exit(1)
			return err
		}
	}
	log.Infof("subconverter is already up to date")
	return nil
}

func UpdateSubconverter() error {

	arch := runtime.GOARCH
	goos := runtime.GOOS

	var downloadUrl string
	baseUrl := op.GetSettingStr(setting.SUBCONVERTER_URL)

	var filename string
	switch goos {
	case "windows":
		switch arch {
		case "386":
			filename = "subconverter_win32.exe"
		case "amd64":
			filename = "subconverter_win64.exe"
		default:
			log.Errorf("unsupported windows architecture: %s", arch)
			return fmt.Errorf("unsupported windows architecture: %s", arch)
		}
	case "darwin":
		switch arch {
		case "amd64":
			filename = "subconverter_darwin64"
		case "arm64":
			filename = "subconverter_darwinarm"
		default:
			log.Errorf("unsupported darwin architecture: %s", arch)
			return fmt.Errorf("unsupported darwin architecture: %s", arch)
		}
	case "linux":
		switch arch {
		case "386":
			filename = "subconverter_linux32"
		case "amd64":
			filename = "subconverter_linux64"
		case "arm":
			filename = "subconverter_armv7"
		case "arm64":
			filename = "subconverter_aarch64"
		default:
			log.Errorf("unsupported linux architecture: %s", arch)
			return fmt.Errorf("unsupported linux architecture: %s", arch)
		}
	default:
		log.Errorf("unsupported operating system: %s", goos)
		return fmt.Errorf("unsupported operating system: %s", goos)
	}

	downloadUrl = baseUrl + "/" + filename

	bytes, err := download(downloadUrl, op.GetSettingBool(setting.SUBCONVERTER_URL_PROXY))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.Base().SubConverter.Path, 0755); err != nil {
		log.Errorf("failed to create directory: %v", err)
		return err
	}

	targetPath := config.Base().SubConverter.Path + "/subconverter"
	if goos == "windows" {
		targetPath += ".exe"
	}

	err = os.WriteFile(targetPath, bytes, 0755)
	if err != nil {
		return err
	}

	log.Infof("update subconverter success: %s", filename)

	return nil
}
