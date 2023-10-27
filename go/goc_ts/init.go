package main

import (
	"flag"
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
	"strings"
)

var (
	help        bool
	version     bool
	dir         string
	outDir      string
	indentation int
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.StringVar(&dir, "d", "./go/", "input dir\n"+
		"we will handle all '*.go' files in it")
	flag.StringVar(&outDir, "o", "./ts/", "output dir")
	flag.IntVar(&indentation, "i", 4, "indentation of generate files")

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

	dir = mustDir(dir)
	outDir = mustDir(outDir)

	data.SetIndentation(indentation)
}

// mustDir make sure 'path' is end with '/'
func mustDir(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}
