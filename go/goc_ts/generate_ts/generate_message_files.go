package generate_ts

import (
	"fmt"
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
)

func GenerateMessageFiles(apiIns *data.API, outDir string) {
	for filename := range apiIns.Message {
		absolutePath := outDir + filename + data.MessageFileSuffix
		generateMessageFile(absolutePath, apiIns.Message[filename])
	}
}

// generateMessageFile write in formatData single function, mainly for 'defer file.close', avoid hold many file handles
func generateMessageFile(absolutePath string, messageItems []*data.MessageItem) {
	file, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open message file(%s) failed, error: %v\n", absolutePath, err))
	}
	defer func() {
		_ = file.Close()
	}()

	content := data.Copyright
	for i := range messageItems {
		content = append(content, formatMessage(messageItems[i])...)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln(fmt.Sprintf("write message file(%s) failed, error: %v\n", absolutePath, err))
	}
}

// formatMessage format 'message item' to ts code
func formatMessage(message *data.MessageItem) []byte {
	res := []byte("\nexport class " + message.Name + " {\n")

	for i := range message.Fields {
		res = append(res, fmt.Sprintf("%s%s?: %s = %s\n",
			data.GetIndentation(1), message.Fields[i].Name,
			message.Fields[i].TSType, message.Fields[i].TSZeroValue)...)
	}

	res = append(res, "}\n"...)

	return res
}
