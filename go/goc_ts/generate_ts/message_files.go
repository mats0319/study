package generate_ts

import (
	"fmt"

	"log"
	"os"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts/code_template"
)

func GenerateMessageFiles(apiIns *data.API, outDir string) {
	for filename := range apiIns.Message {
		absolutePath := outDir + filename + apiIns.Config.MessageFileSuffix
		generateMessageFile(apiIns.Config, absolutePath, apiIns.Message[filename])
	}
}

// generateMessageFile write in formatData single function, mainly for 'defer file.close', avoid hold many file handles
func generateMessageFile(config *data.APIConfig, absolutePath string, messageItems []*data.MessageItem) {
	file, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open message file(%s) failed, error: %v\n", absolutePath, err))
	}
	defer func() {
		_ = file.Close()
	}()

	content := data.Copyright
	for i := range messageItems {
		content = append(content, code_template.FormatMessage(config, messageItems[i])...)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln(fmt.Sprintf("write message file(%s) failed, error: %v\n", absolutePath, err))
	}
}
