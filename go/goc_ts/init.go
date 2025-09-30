package main

import (
	"flag"
	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
)

var (
	help       bool
	version    bool
	configFile string
	initialize bool
	generateGo string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.StringVar(&configFile, "c", "./go/config.json", "config file")
	flag.BoolVar(&initialize, "i", false, "initialize basic files\n"+
		"overwrite 'config_default.json' and 'init_sample.json'")
	flag.StringVar(&generateGo, "g", "./go/init.json", "generate go api definition files\n"+
		"make sure param names are equal")

	flag.Parse()

	if help {
		log.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if version {
		log.Println(data.Version)
		os.Exit(0)
	}
}
