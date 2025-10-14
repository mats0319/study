package data

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/utils"
)

func SetConfig(configFile string) {
	fileBytes, err := os.ReadFile(configFile)

	configIns := &Config{}
	err2 := json.Unmarshal(fileBytes, configIns)

	if err != nil || err2 != nil {
		log.Println("> Goc_ts: Use default config, error: ", err, err2)
	} else {
		log.Println("> Goc_ts: Use config file: ", configFile)
	}

	configIns.useDefaultConfigForEmpty()

	utils.MustDir(configIns.GoDir)
	utils.MustDir(configIns.TsDir)
	utils.MustExistDir(configIns.GoDir)
	utils.EmptyDir(configIns.TsDir)

	// set type maps
	for _, typ := range configIns.BasicGoType {
		for _, goTyp := range typ.GoType {
			GeneratorIns.TsType[goTyp] = typ.TsType
		}
		GeneratorIns.TsZeroValue[typ.TsType] = typ.TsZeroValue
	}

	GeneratorIns.Config = configIns
}

func GetIndentation() []byte {
	res := make([]byte, 0, GeneratorIns.Config.Indentation)
	for i := 0; i < GeneratorIns.Config.Indentation; i++ {
		res = append(res, ' ')
	}
	return res
}

// use default config cover empty config
func (c *Config) useDefaultConfigForEmpty() {
	if len(c.GoDir) < 1 {
		c.GoDir = DefaultGeneratorConfig.GoDir
	}
	if len(c.TsDir) < 1 {
		c.TsDir = DefaultGeneratorConfig.TsDir
	}
	if len(c.BaseURL) < 1 {
		c.BaseURL = DefaultGeneratorConfig.BaseURL
	}
	if c.Timeout < 1 {
		c.Timeout = DefaultGeneratorConfig.Timeout
	}
	if len(c.RequestMessageSuffix) < 1 {
		c.RequestMessageSuffix = DefaultGeneratorConfig.RequestMessageSuffix
	}
	if len(c.ResponseMessageSuffix) < 1 {
		c.ResponseMessageSuffix = DefaultGeneratorConfig.ResponseMessageSuffix
	}
	if len(c.ServiceFileSuffix) < 1 {
		c.ServiceFileSuffix = DefaultGeneratorConfig.ServiceFileSuffix
	}
	if len(c.MessageFileSuffix) < 1 {
		c.MessageFileSuffix = DefaultGeneratorConfig.MessageFileSuffix
	}
	if len(c.BasicGoType) < 1 {
		c.BasicGoType = DefaultGeneratorConfig.BasicGoType
	}
	if c.Indentation < 1 {
		c.Indentation = DefaultGeneratorConfig.Indentation
	}
}
