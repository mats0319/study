package data

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mats9693/study/go/goc-ts/utils"
)

func (ins *Generator) Initialize(configFile string) {
	fileBytes, err := os.ReadFile(configFile)

	configIns := &Config{}
	err2 := json.Unmarshal(fileBytes, configIns)

	if err != nil || err2 != nil {
		log.Println("> Goc_ts: Use default config, error: ", err, err2)
	} else {
		log.Println("> Goc_ts: Use config file: ", configFile)
	}

	configIns.mustValid()
	ins.Config = configIns

	// set type maps
	for _, typIns := range configIns.BasicGoType {
		for _, goTyp := range typIns.GoType {
			ins.TsType[goTyp] = typIns.TsType
		}
		ins.TsZeroValue[typIns.TsType] = typIns.TsZeroValue
	}

	// set indentation
	indentationBytes := make([]byte, 0, ins.Config.Indentation)
	for i := 0; i < ins.Config.Indentation; i++ {
		indentationBytes = append(indentationBytes, ' ')
	}
	ins.IndentationStr = string(indentationBytes)
}

// make sure all configs are valid, use default config cover empty ones
func (c *Config) mustValid() {
	if len(c.GoDir) < 1 {
		c.GoDir = DefaultGeneratorConfig.GoDir
	} else {
		c.GoDir = utils.MustSuffix(c.GoDir, "/")
	}
	if len(c.TsDir) < 1 {
		c.TsDir = DefaultGeneratorConfig.TsDir
	} else {
		c.TsDir = utils.MustSuffix(c.TsDir, "/")
	}
	utils.MustExistDir(c.GoDir)
	utils.MustExistDir(c.GoDir + "backup/")
	utils.EmptyDir(c.TsDir)

	if len(c.BaseURL) < 1 {
		c.BaseURL = DefaultGeneratorConfig.BaseURL
	}
	if c.Timeout < 1 {
		c.Timeout = DefaultGeneratorConfig.Timeout
	}

	if len(c.RequestStructureSuffix) < 1 {
		c.RequestStructureSuffix = DefaultGeneratorConfig.RequestStructureSuffix
	}
	if len(c.ResponseStructureSuffix) < 1 {
		c.ResponseStructureSuffix = DefaultGeneratorConfig.ResponseStructureSuffix
	}
	if len(c.RequestFileSuffix) < 1 {
		c.RequestFileSuffix = DefaultGeneratorConfig.RequestFileSuffix
	} else {
		c.RequestFileSuffix = utils.MustSuffix(c.RequestFileSuffix, ".ts")
	}
	if len(c.StructureFileSuffix) < 1 {
		c.StructureFileSuffix = DefaultGeneratorConfig.StructureFileSuffix
	} else {
		c.StructureFileSuffix = utils.MustSuffix(c.StructureFileSuffix, ".ts")
	}

	if len(c.BasicGoType) < 1 {
		c.BasicGoType = DefaultGeneratorConfig.BasicGoType
	}
	if c.Indentation < 1 {
		c.Indentation = DefaultGeneratorConfig.Indentation
	}
}
