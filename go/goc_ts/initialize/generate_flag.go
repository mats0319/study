package initialize

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var InitializerIns = &Initializer{}

// OnGenerate 实现-g参数：解析go接口文件的初始化文件(init.json)，并生成基础的go接口文件
func OnGenerate(initFileName string) {
	deserializeInitFile(initFileName)
	generateBasicGoAPIFiles()
}

func deserializeInitFile(initFileName string) {
	fileBytes, err := os.ReadFile(initFileName)
	if err != nil {
		log.Println(fmt.Sprintf("read init file(%s) failed, error: %v\n", initFileName, err))
		return
	}

	err = json.Unmarshal(fileBytes, &InitializerIns)
	if err != nil {
		log.Fatalln("json unmarshal failed, error: ", err)
	}
}

func generateBasicGoAPIFiles() {
	for i := range InitializerIns.Files {
		InitializerIns.Files[i].toGo()
	}
}
