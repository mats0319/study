package generate_ts

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
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

	_, err = file.Write([]byte(data.FormatConfigCode(config)))
	if err != nil {
		log.Fatalln("write config file failed, error: ", err)
	}
}
