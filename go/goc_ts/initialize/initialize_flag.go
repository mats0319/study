package initialize

import (
	"encoding/json"
	"log"

	"github.com/mats9693/study/go/goc-ts/data"
	"github.com/mats9693/study/go/goc-ts/utils"
)

// OnInitialize 实现-i参数：写默认配置文件和go接口文件示例文件(demo)
func OnInitialize() {
	writeDefaultConfigFile()
	writeDefaultInitializerFile()
}

func writeDefaultConfigFile() {
	content, err := json.Marshal(data.DefaultGeneratorConfig)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	utils.WriteFile(data.GeneratorIns.Config.GoDir+"config_default.json", content)
}

func writeDefaultInitializerFile() {
	content, err := json.Marshal(DefaultInitializer)
	if err != nil {
		log.Fatalln("json marshal failed, error: ", err)
	}

	utils.WriteFile(data.GeneratorIns.Config.GoDir+"init_default.json", content)
}
