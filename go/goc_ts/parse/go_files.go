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

func ParseGoFiles(apiIns *data.API, dir string) {
	entry, err := os.ReadDir(dir)
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

		parseFile(apiIns, dir, fileInfo.Name())
	}
}

func parseFile(apiIns *data.API, dir string, filename string) {
	fileBytes, err := os.ReadFile(dir + filename)
	if err != nil {
		log.Println(fmt.Sprintf("read go file(%s) failed, error: %v\n", dir+filename, err))
		return
	}

	filename = strings.TrimSuffix(filename, ".go")

	matchService(apiIns, filename, fileBytes)
	matchMessage(apiIns, filename, fileBytes)
}

func matchService(apiIns *data.API, filename string, fileBytes []byte) {
	serviceREMatched := ServiceRE.FindAllSubmatch(fileBytes, -1)
	for i := range serviceREMatched {
		if len(serviceREMatched[i]) < 3 {
			continue
		}

		apiIns.Service[filename] = append(apiIns.Service[filename], &data.ServiceItem{
			Name: string(serviceREMatched[i][1]),
			URI:  string(serviceREMatched[i][2]),
		})
	}
}

func matchMessage(apiIns *data.API, filename string, fileBytes []byte) {
	messageREMatched := MessageRE.FindAllSubmatch(fileBytes, -1)
	for i := range messageREMatched {
		if len(messageREMatched[i]) < 3 {
			continue
		}

		apiIns.Message[filename] = append(apiIns.Message[filename], &data.MessageItem{
			Name:   string(messageREMatched[i][1]),
			Fields: matchMessageField(apiIns, messageREMatched[i][2]),
		})
	}
}

func matchMessageField(apiIns *data.API, field []byte) []*data.MessageField {
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

		messageFieldToTs(apiIns, fieldIns)

		res = append(res, fieldIns)
	}

	return res
}

// messageFieldToTs according to 'field', generate 'ts filed type' and 'ts field zero value'
func messageFieldToTs(apiIns *data.API, field *data.MessageField) {
	v, ok := apiIns.TsType[field.GoType]
	if ok { // base type, in type map
		field.TSType = v
		field.TSZeroValue = apiIns.TsZeroValue[v]
	} else { // not in map, consider as self-define type
		field.TSType = field.GoType
		field.TSZeroValue = fmt.Sprintf("new %s()", field.GoType)
	}

	if field.IsArray {
		field.TSType = fmt.Sprintf("Array<%s>", field.TSType)
		field.TSZeroValue = fmt.Sprintf("new %s()", field.TSType)
	}
}
