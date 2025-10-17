package initialize

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// OnGenerate 实现-g参数：解析go接口文件的初始化文件(init.json)，并生成基础的go接口文件
func OnGenerate(initFileName string) {
	fileBytes, err := os.ReadFile(initFileName)
	if err != nil {
		log.Fatalln(fmt.Sprintf("read init file(%s) failed, error: %v", initFileName, err))
	}

	err = json.Unmarshal(fileBytes, &InitializerIns)
	if err != nil {
		log.Fatalln("json unmarshal failed, error: ", err)
	}

	for i := range InitializerIns.Files {
		InitializerIns.Files[i].writeFile(InitializerIns.PackageName)
	}
}
