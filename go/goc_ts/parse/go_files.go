package parse

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
)

var (
	ServiceRE      = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*"([/\w-]+)"`)
	MessageRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	MessageFieldRE = regexp.MustCompile(`\w+\s+([\[\]\w]+)\s+.*json:"(\w+)".*`)
)

func ParseGoFiles() {
	entry, err := os.ReadDir(data.GeneratorIns.Config.GoDir)
	if err != nil {
		log.Fatalln("read dir failed, error: ", err)
	}

	for i := range entry {
		if entry[i].IsDir() { // ignore folder
			continue
		}

		var fileInfo fs.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			log.Println("get go file info failed, error: ", err)
			continue
		}

		if !strings.HasSuffix(fileInfo.Name(), ".go") { // ignore not go files
			continue
		}

		parseFile(fileInfo.Name())
	}
}

func parseFile(filename string) {
	absolutePath := data.GeneratorIns.Config.GoDir + filename
	fileBytes, err := os.ReadFile(absolutePath)
	if err != nil {
		log.Println(fmt.Sprintf("read go file(%s) failed, error: %v\n", absolutePath, err))
		return
	}

	filename = strings.TrimSuffix(filename, ".go")

	matchService(filename, fileBytes)
	matchMessage(filename, fileBytes)
}

func matchService(filename string, fileBytes []byte) {
	serviceREMatched := ServiceRE.FindAllSubmatch(fileBytes, -1)
	for i := range serviceREMatched {
		if len(serviceREMatched[i]) < 3 {
			continue
		}

		data.GeneratorIns.Services[filename] = append(data.GeneratorIns.Services[filename], &data.ServiceItem{
			Name: string(serviceREMatched[i][1]),
			URI:  string(serviceREMatched[i][2]),
		})
	}
}

func matchMessage(filename string, fileBytes []byte) {
	messageREMatched := MessageRE.FindAllSubmatch(fileBytes, -1)
	for i := range messageREMatched {
		if len(messageREMatched[i]) < 3 {
			continue
		}

		data.GeneratorIns.Messages[filename] = append(data.GeneratorIns.Messages[filename], &data.MessageItem{
			Name:   string(messageREMatched[i][1]),
			Fields: matchMessageFields(messageREMatched[i][2]),
		})
	}
}

func matchMessageFields(field []byte) []*data.MessageField {
	res := make([]*data.MessageField, 0)

	messageFieldREMatched := MessageFieldRE.FindAllSubmatch(field, -1)
	for i := range messageFieldREMatched {
		if len(messageFieldREMatched[i]) < 3 {
			continue
		}

		fieldIns := &data.MessageField{
			Name:    string(messageFieldREMatched[i][2]),
			GoType:  string(bytes.TrimPrefix(messageFieldREMatched[i][1], []byte("[]"))),
			IsArray: bytes.HasPrefix(messageFieldREMatched[i][1], []byte("[]")),
		}

		messageFieldToTs(fieldIns)

		res = append(res, fieldIns)
	}

	return res
}

// messageFieldToTs according to 'field', generate 'ts filed type' and 'ts field zero value'
func messageFieldToTs(field *data.MessageField) {
	v, ok := data.GeneratorIns.TsType[field.GoType]
	if ok { // basic type, in type map
		field.TSType = v
		field.TSZeroValue = data.GeneratorIns.TsZeroValue[v]
	} else { // not in map, consider as self-define type
		field.TSType = field.GoType
		field.TSZeroValue = fmt.Sprintf("new %s()", field.GoType)
	}

	if field.IsArray {
		field.TSType = fmt.Sprintf("Array<%s>", field.TSType)
		field.TSZeroValue = fmt.Sprintf("new %s()", field.TSType)
	}
}
