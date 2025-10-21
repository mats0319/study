package generate_ts

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func GenerateStructureFiles() {
	for filename := range data.GeneratorIns.StructureAffiliation {
		content := utils.Copyright

		externalStructures := make(map[string][]string)
		structuresStr := ""
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			structuresStr += serializeStructure(structureName, externalStructures)
		}
		delete(externalStructures, filename) // not import from current file

		importStructuresStr := serializeStructuresImport(externalStructures)
		if len(importStructuresStr) > 0 {
			content = append(content, '\n')
		}

		enumsStr := ""
		for _, enumName := range data.GeneratorIns.EnumAffiliation[filename] {
			enumsStr += serializeEnum(enumName)
		}

		content = append(content, importStructuresStr...)
		content = append(content, enumsStr...)
		content = append(content, structuresStr...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.StructureFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeStructure(structureName string, externalStructures map[string][]string) string {
	structureItemIns, _ := data.GeneratorIns.Structures[structureName]

	const template = "{{ $indentation }}{{ $fieldName }}: {{ $fieldType_Ts }} = {{ $fieldZeroValue_Ts }};\n"

	fieldsStr := ""
	for _, structureFieldIns := range structureItemIns.Fields {
		field := template
		field = strings.ReplaceAll(field, "{{ $fieldName }}", structureFieldIns.Name)
		field = strings.ReplaceAll(field, "{{ $fieldType_Ts }}", structureFieldIns.TSType)
		field = strings.ReplaceAll(field, "{{ $fieldZeroValue_Ts }}", structureFieldIns.TSZeroValue)

		fieldsStr += field

		// find self-define type, record it
		if _, ok := data.GeneratorIns.TsType[structureFieldIns.GoType]; !ok {
			fromFile, _ := data.GeneratorIns.TypeFrom[structureFieldIns.GoType]
			externalStructures[fromFile] = append(externalStructures[fromFile], structureFieldIns.TSType)
		}
	}
	fieldsStr = strings.ReplaceAll(fieldsStr, "{{ $indentation }}", string(data.GetIndentation()))

	classStr := "\n" +
		"export class {{ $structureName }} {\n" +
		"{{ $structureFields }}" +
		"}\n"
	classStr = strings.ReplaceAll(classStr, "{{ $structureName }}", structureName)
	classStr = strings.ReplaceAll(classStr, "{{ $structureFields }}", fieldsStr)

	return classStr
}

func serializeEnum(enumName string) string {
	enumItemIns, _ := data.GeneratorIns.Enums[enumName]

	const template = "{{ $indentation }}{{ $enumName }} = {{ $enumZeroValue_Ts }},\n"

	enumUnitsStr := ""
	for _, enumUnitIns := range enumItemIns.Units {
		unit := template
		unit = strings.ReplaceAll(unit, "{{ $enumName }}", enumUnitIns.Name)
		unit = strings.ReplaceAll(unit, "{{ $enumZeroValue_Ts }}", enumUnitIns.Value)

		enumUnitsStr += unit
	}
	enumUnitsStr = strings.ReplaceAll(enumUnitsStr, "{{ $indentation }}", string(data.GetIndentation()))

	enumStr := "\n" +
		"export enum {{ $enumName }} {\n" +
		"{{ $enumUnits }}" +
		"}\n"
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumName }}", enumName)
	enumStr = strings.ReplaceAll(enumStr, "{{ $enumUnits }}", enumUnitsStr)

	return enumStr
}

// structures: from filename - structures' name
func serializeStructuresImport(externalStructures map[string][]string) string {
	if len(externalStructures) < 1 {
		return ""
	}

	const template = `import { {{ $structures }} } from "./{{ $filename }}{{ $structureFileSuffix }}"` + "\n"

	importStructuresStr := ""
	for key, value := range externalStructures {
		str := template
		str = strings.ReplaceAll(str, "{{ $structures }}", utils.FormatStrSliceInLine(value))
		str = strings.ReplaceAll(str, "{{ $filename }}", key)

		importStructuresStr += str
	}
	importStructuresStr = strings.ReplaceAll(importStructuresStr, "{{ $structureFileSuffix }}", strings.TrimSuffix(data.GeneratorIns.Config.StructureFileSuffix, ".ts"))

	return importStructuresStr
}
