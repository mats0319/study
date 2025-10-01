package main

import (
	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts"
	"github.com/mats9693/study/go/goc_ts/parse"
)

var apiIns = &data.API{
	Config:      &data.APIConfig{},
	Utils:       &data.APIUtils{},
	Service:     make(map[string][]*data.ServiceItem),
	Message:     make(map[string][]*data.MessageItem),
	TsType:      make(map[string]string),
	TsZeroValue: make(map[string]string),
}

func main() {
	log.Println("> Goc_ts: Start.")
	defer log.Println("> Goc_ts: Finish.")

	apiIns.SetConfigFromFile(configFile)

	emptyOutDir(apiIns.Config.OutDir)

	parse.ParseGoFiles(apiIns, apiIns.Config.Dir)
	parse.ParseUtils(apiIns)

	generate_ts.GenerateConfig(apiIns.Config, apiIns.Config.OutDir)
	generate_ts.GenerateUtils(apiIns.Utils, apiIns.Config.OutDir)
	generate_ts.GenerateServiceFiles(apiIns, apiIns.Config.OutDir)
	generate_ts.GenerateMessageFiles(apiIns, apiIns.Config.OutDir)
}

func emptyOutDir(outDir string) {
	err := os.RemoveAll(outDir)
	if err != nil {
		log.Fatalln("clean out dir failed, error: ", err)
	}

	err = os.MkdirAll(outDir, 0777)
	if err != nil {
		log.Fatalln("'mkdir' on out dir failed, error: ", err)
	}
}
