package parse

import (
	"encoding/json"
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
)

func ParseConfig(apiIns *data.API, dir string) {
	file, err := os.Open(dir + "config.json")
	if err != nil { // no valid config file, just use default value
		log.Println("open config file failed, error: ", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln("get config file stat failed, error: ", err)
	}

	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Fatalln("read config file failed, error: ", err)
	}

	apiConfigIns := &data.APIConfig{}
	err = json.Unmarshal(fileBytes, apiConfigIns)
	if err != nil {
		log.Fatalln("deserialize api config failed, error: ", err)
	}

	apiIns.Config = apiConfigIns
}
