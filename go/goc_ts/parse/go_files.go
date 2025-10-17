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
	RequestRE        = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*"([/\w-]+)"`)
	StructureRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructureFieldRE = regexp.MustCompile(`\w+\s+([\[\]\w]+)\s+.*json:"(\w+)".*`)
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
		log.Fatalln(fmt.Sprintf("read go file(%s) failed, error: %v", absolutePath, err))
	}

	filename = strings.TrimSuffix(filename, ".go")

	matchRequests(filename, fileBytes)
	matchStructures(filename, fileBytes)
}

func matchRequests(filename string, fileBytes []byte) {
	requestREMatched := RequestRE.FindAllSubmatch(fileBytes, -1)
	for i := range requestREMatched {
		if len(requestREMatched[i]) < 3 {
			continue
		}

		requestName := string(requestREMatched[i][1])
		requestURI := string(requestREMatched[i][2])

		data.GeneratorIns.RequestAffiliation[filename] = append(data.GeneratorIns.RequestAffiliation[filename], requestName)
		data.GeneratorIns.Requests[requestName] = requestURI
	}
}

func matchStructures(filename string, fileBytes []byte) {
	structureREMatched := StructureRE.FindAllSubmatch(fileBytes, -1)
	for i := range structureREMatched {
		if len(structureREMatched[i]) < 3 {
			continue
		}

		structureName := string(structureREMatched[i][1])
		if !strings.HasSuffix(structureName, data.GeneratorIns.Config.RequestStructureSuffix) &&
			!strings.HasSuffix(structureName, data.GeneratorIns.Config.ResponseStructureSuffix) {
			// self-define struct, record it's from
			data.GeneratorIns.StructureFrom[structureName] = filename
		}

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], structureName)
		data.GeneratorIns.Structures[structureName] = &data.StructureItem{
			Fields: matchStructureFields(structureREMatched[i][2]),
		}
	}
}

func matchStructureFields(field []byte) []*data.StructureField {
	res := make([]*data.StructureField, 0)

	structureFieldREMatched := StructureFieldRE.FindAllSubmatch(field, -1)
	for i := range structureFieldREMatched {
		if len(structureFieldREMatched[i]) < 3 {
			continue
		}

		fieldIns := &data.StructureField{
			Name:    string(structureFieldREMatched[i][2]),
			GoType:  string(bytes.TrimPrefix(structureFieldREMatched[i][1], []byte("[]"))),
			IsArray: bytes.HasPrefix(structureFieldREMatched[i][1], []byte("[]")),
		}

		getTsTypeAndZeroValue(fieldIns)

		res = append(res, fieldIns)
	}

	return res
}

// getTsTypeAndZeroValue according to 'field', record 'ts filed type' and 'ts field zero value'
func getTsTypeAndZeroValue(field *data.StructureField) {
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
