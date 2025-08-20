package update

import (
	"os"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func InitUI() error {
	if _, err := os.Stat(config.Base().Server.UIPath + "/index.html"); err != nil {
		err = UpdateUI()
		if err != nil {
			log.Warnf("auto update ui failed, please download ui manually from %s and unzip to %s: %v", op.GetSettingStr(setting.FRONTEND_URL), config.Base().Server.UIPath, err)
			return err
		}
	}
	log.Infof("UI is already up to date")
	return nil
}

func UpdateUI() error {

	bytes, err := download(op.GetSettingStr(setting.FRONTEND_URL), op.GetSettingBool(setting.FRONTEND_URL_PROXY))
	if err != nil {
		return err
	}

	if err := unzip(bytes, config.Base().Server.UIPath); err != nil {
		return err
	}

	log.Infof("update ui success")

	return nil
}
