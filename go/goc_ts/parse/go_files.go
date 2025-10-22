package parse

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mats9693/study/go/goc-ts/data"
)

var (
	RequestRE        = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*"([\w/-]+)"`)
	StructureRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructureFieldRE = regexp.MustCompile(`\w+\s+([\[\]\w]+)\s+.*json:"(\w+)"`)
	EnumRE           = regexp.MustCompile(`type\s+(\w+)\s*=?\s*\w+\s+const\s*\(([^)]*)\)`)
	EnumUnitRE       = regexp.MustCompile(`(\w+)\s+(\w+)\s*=\s*(\w+)`)
)

func TraversalDir() {
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

		parseGoFile(fileInfo.Name())
	}

	setTsTypeAndZeroValue()
}

func parseGoFile(filename string) {
	absolutePath := data.GeneratorIns.Config.GoDir + filename
	fileBytes, err := os.ReadFile(absolutePath)
	if err != nil {
		log.Fatalln(fmt.Sprintf("read go file(%s) failed, error: %v", absolutePath, err))
	}

	filename = strings.TrimSuffix(filename, ".go")

	matchRequests(filename, fileBytes)
	matchStructures(filename, fileBytes)
	matchEnums(filename, fileBytes)
}

func matchRequests(filename string, fileBytes []byte) {
	requestREMatched := RequestRE.FindAllSubmatch(fileBytes, -1)
	for i := range requestREMatched {
		// 没什么意义，因为匹配函数如果找到了匹配项，这里的长度就必然不小于3,但是因为后续要直接使用下标访问，所以还是判断一下
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

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], structureName)
		data.GeneratorIns.StructureFrom[structureName] = filename
		data.GeneratorIns.Structures[structureName] = &data.StructureItem{
			Typ:    &data.StructureType{IsStruct: true},
			Fields: matchStructureFields(structureREMatched[i][2]),
		}
	}
}

func matchStructureFields(fields []byte) []*data.StructureField {
	fieldSlice := make([]*data.StructureField, 0)

	structureFieldREMatched := StructureFieldRE.FindAllSubmatch(fields, -1)
	for i := range structureFieldREMatched {
		if len(structureFieldREMatched[i]) < 3 {
			continue
		}

		fieldIns := &data.StructureField{
			Name:    string(structureFieldREMatched[i][2]),
			GoType:  string(bytes.TrimPrefix(structureFieldREMatched[i][1], []byte("[]"))),
			IsArray: bytes.HasPrefix(structureFieldREMatched[i][1], []byte("[]")),
		}

		fieldSlice = append(fieldSlice, fieldIns)
	}

	return fieldSlice
}

func matchEnums(filename string, fileBytes []byte) {
	enumREMatched := EnumRE.FindAllSubmatch(fileBytes, -1)
	for i := range enumREMatched {
		if len(enumREMatched[i]) < 3 {
			continue
		}

		enumName := string(enumREMatched[i][1])

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], enumName)
		data.GeneratorIns.StructureFrom[enumName] = filename
		data.GeneratorIns.Structures[enumName] = &data.StructureItem{
			Typ:    &data.StructureType{IsEnum: true},
			Fields: matchEnumUnits(enumREMatched[i][2]),
		}
	}
}

func matchEnumUnits(fileBytes []byte) []*data.StructureField {
	enumUnitSlice := make([]*data.StructureField, 0)

	enumUnitREMatched := EnumUnitRE.FindAllSubmatch(fileBytes, -1)
	for i := range enumUnitREMatched {
		if len(enumUnitREMatched[i]) < 4 {
			continue
		}

		enumName := string(enumUnitREMatched[i][2])
		enumUnitIns := &data.StructureField{
			Name:        strings.TrimPrefix(string(enumUnitREMatched[i][1]), enumName+"_"),
			GoType:      enumName,
			TSType:      enumName,
			TSZeroValue: string(enumUnitREMatched[i][3]),
		}

		enumUnitSlice = append(enumUnitSlice, enumUnitIns)
	}

	return enumUnitSlice
}

// according to 'field', record 'ts filed type' and 'ts field zero value'
func setTsTypeAndZeroValue() {
	for i, structureItemIns := range data.GeneratorIns.Structures {
		if structureItemIns.Typ.IsEnum {
			continue
		}

		for j, fieldIns := range structureItemIns.Fields {
			v, ok := data.GeneratorIns.Structures[fieldIns.GoType]
			if !ok { // basic type, in type map
				tsType, _ := data.GeneratorIns.TsType[fieldIns.GoType]
				tsZeroValue, _ := data.GeneratorIns.TsZeroValue[tsType]
				data.GeneratorIns.Structures[i].Fields[j].TSType = tsType
				data.GeneratorIns.Structures[i].Fields[j].TSZeroValue = tsZeroValue
			} else { // not in map, consider as self-define type
				data.GeneratorIns.Structures[i].Fields[j].TSType = fieldIns.GoType
				if v.Typ.IsEnum { // 这个go结构体字段的类型是枚举类型
					data.GeneratorIns.Structures[i].Fields[j].TSZeroValue = "0"
				} else {
					data.GeneratorIns.Structures[i].Fields[j].TSZeroValue = fmt.Sprintf("new %s()", fieldIns.GoType)
				}
			}

			if fieldIns.IsArray {
				data.GeneratorIns.Structures[i].Fields[j].TSType = fmt.Sprintf("Array<%s>", data.GeneratorIns.Structures[i].Fields[j].TSType)
				data.GeneratorIns.Structures[i].Fields[j].TSZeroValue = fmt.Sprintf("new %s()", data.GeneratorIns.Structures[i].Fields[j].TSType)
			}
		}
	}
}
