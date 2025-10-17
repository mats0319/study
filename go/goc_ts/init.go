package main

import (
	"flag"
	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/initialize"
	"github.com/mats9693/study/go/goc_ts/utils"
)

var (
	help       bool
	version    bool
	configFile string

	initializeFlag bool

	generateGoFlag bool
	generateGoFrom string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.StringVar(&configFile, "c", "./go/config.json", "config file")
	flag.BoolVar(&initializeFlag, "i", false, "initializeFlag basic files\n"+
		"overwrite './go/config_default.json' and './go/init_default.json'")
	flag.BoolVar(&generateGoFlag, "g", false, "generate go files from './go/init.json'\n")
	flag.StringVar(&generateGoFrom, "genFrom", "", "generate go files from given file")

	flag.Parse()

	if help {
		log.Println("Options: \n(In this help, './go/' means our go files dir)")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if version {
		log.Println(utils.Version)
		os.Exit(0)
	}

	data.Initialize(configFile)

	if initializeFlag {
		initialize.OnInitialize()
		os.Exit(0)
	}

	if generateGoFlag || len(generateGoFrom) > 0 {
		initFileName := ""
		if generateGoFlag {
			initFileName = data.GeneratorIns.Config.GoDir + "init.json"
		} else {
			initFileName = generateGoFrom
		}

		initialize.OnGenerate(initFileName)
		os.Exit(0)
	}
}
