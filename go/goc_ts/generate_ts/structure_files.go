package generate_ts

import (
	"strings"

	"github.com/mats9693/study/go/goc-ts/data"
	"github.com/mats9693/study/go/goc-ts/utils"
)

func GenerateStructureFiles() {
	for filename := range data.GeneratorIns.StructureAffiliation {
		content := utils.Copyright

		externalStructures := make(map[string][]string)
		structuresStr := ""
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			structuresStr += serializeStructure(structureName, externalStructures)
		}
		structuresStr = strings.ReplaceAll(structuresStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)
		delete(externalStructures, filename) // not import from current file

		importStructuresStr := serializeStructuresImport(externalStructures)
		if len(importStructuresStr) > 0 {
			content = append(content, '\n')
		}

		content = append(content, importStructuresStr...)
		content = append(content, structuresStr...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.StructureFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeStructure(structureName string, externalStructures map[string][]string) string {
	structureItemIns, _ := data.GeneratorIns.Structures[structureName]

	structureStr := ""
	if structureItemIns.Typ.IsStruct {
		structureStr = serializeStruct(structureName, structureItemIns, externalStructures)
	} else {
		structureStr = serializeEnum(structureName, structureItemIns)
	}

	return structureStr
}

func serializeStruct(structureName string, structureItemIns *data.StructureItem, externalStructures map[string][]string) string {
	fieldsStr := ""
	for _, fieldIns := range structureItemIns.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldType_Ts }} = {{ $fieldZeroValue_Ts }};\n"
		field = strings.ReplaceAll(field, "{{ $fieldName }}", fieldIns.Name)
		field = strings.ReplaceAll(field, "{{ $fieldType_Ts }}", fieldIns.TSType)
		field = strings.ReplaceAll(field, "{{ $fieldZeroValue_Ts }}", fieldIns.TSZeroValue)

		fieldsStr += field

		if _, ok := data.GeneratorIns.Structures[fieldIns.GoType]; ok {
			fromFile, _ := data.GeneratorIns.StructureFrom[fieldIns.GoType]
			externalStructures[fromFile] = append(externalStructures[fromFile], fieldIns.TSType)
		}
	}

	structStr := "\nexport class {{ $structureName }} {\n{{ $structureFields }}}\n"
	structStr = strings.ReplaceAll(structStr, "{{ $structureName }}", structureName)
	structStr = strings.ReplaceAll(structStr, "{{ $structureFields }}", fieldsStr)

	return structStr
}

func serializeEnum(enumName string, enumItemIns *data.StructureItem) string {
	enumUnitsStr := ""
	for _, enumUnitIns := range enumItemIns.Fields {
		unit := "{{ $indentation }}{{ $enumName }} = {{ $enumZeroValue_Ts }},\n"
		unit = strings.ReplaceAll(unit, "{{ $enumName }}", enumUnitIns.Name)
		unit = strings.ReplaceAll(unit, "{{ $enumZeroValue_Ts }}", enumUnitIns.TSZeroValue)

		enumUnitsStr += unit
	}

	enumStr := "\nexport enum {{ $enumName }} {\n{{ $enumUnits }}}\n"
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumName }}", enumName)
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumUnits }}", enumUnitsStr)

	return enumStr
}

// structures: from filename - structures' name
func serializeStructuresImport(externalStructures map[string][]string) string {
	if len(externalStructures) < 1 {
		return ""
	}

	importStructuresStr := ""
	for fromFile, structureNames := range externalStructures {
		str := "import { {{ $structures }} } from \"./{{ $filename }}{{ $structureFileSuffix }}\"\n"
		str = strings.ReplaceAll(str, "{{ $structures }}", utils.FormatStrSliceInLine(structureNames))
		str = strings.ReplaceAll(str, "{{ $filename }}", fromFile)

		importStructuresStr += str
	}
	importStructuresStr = strings.ReplaceAll(importStructuresStr, "{{ $structureFileSuffix }}", strings.TrimSuffix(data.GeneratorIns.Config.StructureFileSuffix, ".ts"))

	return importStructuresStr
}
