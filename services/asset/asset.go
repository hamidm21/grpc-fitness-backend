package asset

import (
	"fmt"

	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
)

type AssetManager struct {
	relPath   string
	serverUrl string
}

var manager *AssetManager

func (instance AssetManager) Initialize() func() {
	manager = new(AssetManager)
	manager.serverUrl = config.Config().GetString(utils.KeyMediaServer)
	manager.relPath = config.Config().GetString(utils.KeyMediaRelPath)
	return nil
}

func (instance AssetManager) MediaUrlExt(fileName string, fileExt string) string {
	return fmt.Sprintf("%s%s%s.%s", instance.serverUrl, instance.relPath, fileName, fileExt)
}

func (instance AssetManager) MediaUrl(fileName string) string {
	return fmt.Sprintf("%s%s%s", instance.serverUrl, instance.relPath, fileName)
}

func GetAssetManager() *AssetManager {
	return manager
}

func init() {
	initializer.Register(AssetManager{}, initializer.LowPriority)
}
