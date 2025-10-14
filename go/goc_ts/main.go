package main

import (
	"log"

	"github.com/mats9693/study/go/goc_ts/generate_ts"
	"github.com/mats9693/study/go/goc_ts/parse"
)

func main() {
	log.Println("> Goc_ts: Start.")
	defer log.Println("> Goc_ts: Finish.")

	parse.ParseGoFiles()
	parse.ParseUtils()

	generate_ts.GenerateConfig()
	generate_ts.GenerateUtils()
	generate_ts.GenerateServiceFiles()
	generate_ts.GenerateMessageFiles()
}
