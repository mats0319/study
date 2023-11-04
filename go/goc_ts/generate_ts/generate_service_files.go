package generate_ts

import (
	"fmt"
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
)

func GenerateServiceFiles(apiIns *data.API, outDir string) {
	for filename := range apiIns.Service {
		absolutePath := outDir + filename + data.ServiceFileSuffix
		generateServiceFile(absolutePath, apiIns.Service[filename], apiIns.Message[filename], filename)
	}
}

func generateServiceFile(
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
	content = append(content, data.FormatServiceCode(serviceItems, messageItems, filename)...)

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln(fmt.Sprintf("write service file(%s) failed, error: %v\n", absolutePath, err))
	}
}
