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
	EnumTypeDefineRE = regexp.MustCompile(`type\s+(\w+)\s*=?\s*(\w+)`)
	EnumRE           = regexp.MustCompile(`const\s*\(([^)]*)\)`)
	EnumUnitRE       = regexp.MustCompile(`\s*(\w+)\s+(\w+)\s*=\s*(\w+)\s*`)
	StructureRE      = regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]*)}`)
	StructureFieldRE = regexp.MustCompile(`\w+\s+([\[\]\w]+)\s+.*json:"(\w+)"`)
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
	matchEnums(filename, fileBytes)
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

func matchEnums(filename string, fileBytes []byte) {
	matchEnumTypeDefines(filename, fileBytes)

	matchEnum(fileBytes)
}

func matchEnumTypeDefines(filename string, fileBytes []byte) {
	EnumTypeDefineREMatched := EnumTypeDefineRE.FindAllSubmatch(fileBytes, -1)
	for i := range EnumTypeDefineREMatched {
		if len(EnumTypeDefineREMatched[i]) < 3 || string(EnumTypeDefineREMatched[i][2]) == "struct" {
			continue
		}

		enumName := string(EnumTypeDefineREMatched[i][1])
		data.GeneratorIns.EnumAffiliation[filename] = append(data.GeneratorIns.EnumAffiliation[filename], enumName)
		data.GeneratorIns.TypeFrom[enumName] = filename

		//data.GeneratorIns.TsType[enumName] = enumName
		data.GeneratorIns.TsZeroValue[enumName] = "0"
	}
}

func matchEnum(fileBytes []byte) {
	enumREMatched := EnumRE.FindAllSubmatch(fileBytes, -1)
	for i := range enumREMatched {
		if len(enumREMatched[i]) < 2 {
			continue
		}

		matchEnumUnits(enumREMatched[i][1])
	}
}

func matchEnumUnits(fileBytes []byte) {
	enumUnitREMatched := EnumUnitRE.FindAllSubmatch(fileBytes, -1)
	for i := range enumUnitREMatched {
		if len(enumUnitREMatched[i]) < 4 {
			continue
		}

		enumName := string(enumUnitREMatched[i][2])
		enumUnitName := string(enumUnitREMatched[i][1])
		if !strings.HasPrefix(enumUnitName, enumName+"_") {
			continue
		}

		enumItemIns, ok := data.GeneratorIns.Enums[enumName]
		if !ok {
			enumItemIns = &data.EnumItem{Units: make([]*data.EnumUnit, 0)}
		}

		enumItemIns.Units = append(enumItemIns.Units, &data.EnumUnit{
			Name:  enumUnitName,
			Value: string(enumUnitREMatched[i][3]),
		})

		data.GeneratorIns.Enums[enumName] = enumItemIns
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
			data.GeneratorIns.TypeFrom[structureName] = filename
		}

		data.GeneratorIns.StructureAffiliation[filename] = append(data.GeneratorIns.StructureAffiliation[filename], structureName)
		data.GeneratorIns.Structures[structureName] = &data.StructureItem{
			Fields: matchStructureFields(structureREMatched[i][2]),
		}
	}
}

func matchStructureFields(fields []byte) []*data.StructureField {
	res := make([]*data.StructureField, 0)

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
