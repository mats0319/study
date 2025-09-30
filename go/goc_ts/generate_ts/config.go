package generate_ts

import (
	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts/code_template"
)

func GenerateConfig(config *data.APIConfig, outDir string) {
	file, err := os.OpenFile(outDir+"config.ts", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("open config file failed, error: ", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(data.Copyright)
	if err != nil {
		log.Fatalln("write config file failed, error: ", err)
	}

	_, err = file.Write([]byte(code_template.FormatConfigCode(config)))
	if err != nil {
		log.Fatalln("write config file failed, error: ", err)
	}
}
