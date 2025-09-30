package code_template

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
)

// FormatMessage format 'message item' to ts code
func FormatMessage(config *data.APIConfig, message *data.MessageItem) string {
	fieldsStr := ""
	for i := range message.Fields {
		field := "{{ $indentation }}{{ $fieldName }}: {{ $fieldTSType }} = {{ $fieldTSZeroValue }};\n"
		field = strings.ReplaceAll(field, "{{ $indentation }}", string(config.GetIndentation(1)))
		field = strings.ReplaceAll(field, "{{ $fieldName }}", message.Fields[i].Name)
		field = strings.ReplaceAll(field, "{{ $fieldTSType }}", message.Fields[i].TSType)
		field = strings.ReplaceAll(field, "{{ $fieldTSZeroValue }}", message.Fields[i].TSZeroValue)

		fieldsStr += field
	}

	res := "\n" +
		"export class {{ $messageName }} {\n" +
		"{{ $messageCode_Fields }}" +
		"}\n"
	res = strings.ReplaceAll(res, "{{ $messageName }}", message.Name)
	res = strings.ReplaceAll(res, "{{ $messageCode_Fields }}", fieldsStr)

	return res
}
