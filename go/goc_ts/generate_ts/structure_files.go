package generate_ts

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func GenerateStructureFiles() {
	for filename := range data.GeneratorIns.StructureAffiliation {
		content := utils.Copyright

		structures := make(map[string][]string)
		messages := ""
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			messages += serializeStructure(structureName, structures)
		}
		delete(structures, filename) // not import from current file

		structuresStr := serializeStructuresImport(structures)
		if len(structuresStr) > 0 {
			content = append(content, '\n')
		}

		content = append(content, structuresStr...)
		content = append(content, messages...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.StructureFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeStructure(structureName string, structures map[string][]string) string {
	structureItemIns, _ := data.GeneratorIns.Structures[structureName]

	fieldsStr := ""
	for _, structureFieldIns := range structureItemIns.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldType_Ts }} = {{ $fieldZeroValue_Ts }};\n"
		field = strings.ReplaceAll(field, "{{ $fieldName }}", structureFieldIns.Name)
		field = strings.ReplaceAll(field, "{{ $fieldType_Ts }}", structureFieldIns.TSType)
		field = strings.ReplaceAll(field, "{{ $fieldZeroValue_Ts }}", structureFieldIns.TSZeroValue)

		fieldsStr += field

		// find self-define type, record it
		if _, ok := data.GeneratorIns.TsType[structureFieldIns.GoType]; !ok {
			fromFile, _ := data.GeneratorIns.StructureFrom[structureFieldIns.GoType]
			structures[fromFile] = append(structures[fromFile], structureFieldIns.GoType)
		}
	}
	fieldsStr = strings.ReplaceAll(fieldsStr, "{{ $indentation }}", string(data.GetIndentation()))

	res := "\n" +
		"export class {{ $structureName }} {\n" +
		"{{ $structureFields }}" +
		"}\n"
	res = strings.ReplaceAll(res, "{{ $structureName }}", structureName)
	res = strings.ReplaceAll(res, "{{ $structureFields }}", fieldsStr)

	return res
}

// structures: from filename - structures' name
func serializeStructuresImport(structures map[string][]string) string {
	if len(structures) < 1 {
		return ""
	}

	const template = `import { {{ $structures }} } from "./{{ $filename }}{{ $structureFileSuffix }}"` + "\n"
	res := ""
	for key, value := range structures {
		res += template
		res = strings.ReplaceAll(res, "{{ $structures }}", utils.FormatStrSliceInLine(value))
		res = strings.ReplaceAll(res, "{{ $filename }}", key)
	}
	res = strings.ReplaceAll(res, "{{ $structureFileSuffix }}", strings.TrimSuffix(data.GeneratorIns.Config.StructureFileSuffix, ".ts"))

	return res
}
