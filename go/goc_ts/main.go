package main

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts"
	"github.com/mats9693/study/go/goc_ts/parse"
	"log"
	"os"
)

var apiIns = &data.API{
	Config: &data.APIConfig{
		BaseURL: "http://127.0.0.1:9693",
		Timeout: 3_000,
	},
	Utils:   &data.APIUtils{},
	Service: make(map[string][]*data.ServiceItem),
	Message: make(map[string][]*data.MessageItem),
}

func main() {
	log.Println("> Goc_ts: Start.")
	defer log.Println("> Goc_ts: Finish.")

	emptyOutDir(outDir)

	parse.ParseConfig(apiIns, dir)
	parse.ParseGoFiles(apiIns, dir)
	parse.ParseUtils(apiIns)

	generate_ts.GenerateConfig(apiIns.Config, outDir)
	generate_ts.GenerateUtils(apiIns.Utils, outDir)
	generate_ts.GenerateServiceFiles(apiIns, outDir)
	generate_ts.GenerateMessageFiles(apiIns, outDir)
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
