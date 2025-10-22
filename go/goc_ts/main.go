package main

import (
	"log"

	"github.com/mats9693/study/go/goc-ts/generate_ts"
	"github.com/mats9693/study/go/goc-ts/parse"
)

func main() {
	log.Println("> Goc_ts: Start.")
	defer log.Println("> Goc_ts: Finish.")

	parse.TraversalDir()
	parse.ParseUtils()

	generate_ts.GenerateConfig()
	generate_ts.GenerateUtils()
	generate_ts.GenerateRequestFiles()
	generate_ts.GenerateStructureFiles()
}
