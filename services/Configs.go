package services

import (
	"go-bootstrap-demo/models"
	"os"

	"github.com/hicode0101/go-logger"
	"github.com/hicode0101/go-utils"
)

// resolveDataPath 定位数据文件：优先当前工作目录，找不到时尝试上一级目录
// （兼容 go test 以包目录为工作目录、而 config.json 在项目根目录的场景）
func resolveDataPath(name string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}
	parent := "../" + name
	if _, err := os.Stat(parent); err == nil {
		return parent
	}
	return name
}

func LoadConfig() {
	//加载配置文件
	configPath := resolveDataPath("config.json")
	logger.Debug("Config Path：", configPath)
	configData, configErr := os.ReadFile(configPath)
	if configErr != nil {
		logger.Error("Read config File err:", configErr)
		return
	}

	ConfigS = new(models.Config)
	err := utils.Json.FromJson(configData, ConfigS)
	if err != nil {
		panic(err)
		return
	}
}

func StoresConfig() {
	//保存配置文件
	jsonStrs, err := utils.Json.ToPrettyJson(ConfigS)
	if err != nil {
		logger.Error("ToPrettyJson config File err:", err)
		return
	}

	err = os.WriteFile(resolveDataPath("config.json"), jsonStrs, os.ModePerm)
	if err != nil {
		logger.Error("Write config File err:", err)
		return
	}

}

func ReLoadConfig() {
	LoadConfig()
	logger.Debug(utils.Json.ToPrettyJsonString(ConfigS))
}
