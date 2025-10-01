package generate_ts

import (
	"fmt"

	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts/code_template"
)

func GenerateServiceFiles(apiIns *data.API, outDir string) {
	for filename := range apiIns.Service {
		absolutePath := outDir + filename + apiIns.Config.ServiceFileSuffix
		generateServiceFile(apiIns.Config, absolutePath, apiIns.Service[filename], apiIns.Message[filename], filename)
	}
}

func generateServiceFile(
	config *data.APIConfig,
	absolutePath string,
	serviceItems []*data.ServiceItem,
	messageItems []*data.MessageItem,
	filename string,
) {
	file, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open service file(%s) failed, error: %v\n", absolutePath, err))
	}
	defer func() {
		_ = file.Close()
	}()

	content := data.Copyright
	content = append(content, code_template.FormatServiceCode(config, serviceItems, messageItems, filename)...)

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln(fmt.Sprintf("write service file(%s) failed, error: %v\n", absolutePath, err))
	}
}
