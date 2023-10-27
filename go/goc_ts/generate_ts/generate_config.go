package generate_ts

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
	"strconv"
	"strings"
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

	_, err = file.Write([]byte(fillConfigCode(config)))
	if err != nil {
		log.Fatalln("write config file failed, error: ", err)
	}
}

func fillConfigCode(config *data.APIConfig) string {
	res := data.ConfigCode
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $baseURL }}", config.BaseURL)
	res = strings.ReplaceAll(res, "{{ $timeout }}", strconv.Itoa(int(config.Timeout)))

	return res
}
