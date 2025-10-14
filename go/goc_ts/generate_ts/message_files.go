package generate_ts

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func GenerateMessageFiles() {
	for filename := range data.GeneratorIns.Messages {
		content := utils.Copyright
		for i := range data.GeneratorIns.Messages[filename] {
			content = append(content, serializeMessage(data.GeneratorIns.Messages[filename][i])...)
		}

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.MessageFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeMessage(message *data.MessageItem) string {
	fieldsStr := ""
	for i := range message.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldTSType }} = {{ $fieldTSZeroValue }};\n"
		field = strings.ReplaceAll(field, "{{ $fieldName }}", message.Fields[i].Name)
		field = strings.ReplaceAll(field, "{{ $fieldTSType }}", message.Fields[i].TSType)
		field = strings.ReplaceAll(field, "{{ $fieldTSZeroValue }}", message.Fields[i].TSZeroValue)

		fieldsStr += field
	}
	fieldsStr = strings.ReplaceAll(fieldsStr, "{{ $indentation }}", string(data.GetIndentation()))

	res := "\n" +
		"export class {{ $messageName }} {\n" +
		"{{ $messageCode_Fields }}" +
		"}\n"
	res = strings.ReplaceAll(res, "{{ $messageName }}", message.Name)
	res = strings.ReplaceAll(res, "{{ $messageCode_Fields }}", fieldsStr)

	return res
}
