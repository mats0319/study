package data

import "strings"

// FormatMessage format 'message item' to ts code
func FormatMessage(message *MessageItem) string {
	fieldsStr := ""
	for i := range message.Fields {
		field := "{{ $indentation }}{{ $fieldName }}?: {{ $fieldTSType }} = {{ $fieldTSZeroValue }}\n"
		field = strings.ReplaceAll(field, "{{ $indentation }}", string(GetIndentation(1)))
		field = strings.ReplaceAll(field, "{{ $fieldName }}", message.Fields[i].Name)
		field = strings.ReplaceAll(field, "{{ $fieldTSType }}", message.Fields[i].TSType)
		field = strings.ReplaceAll(field, "{{ $fieldTSZeroValue }}", message.Fields[i].TSZeroValue)

		fieldsStr += field
	}

	res := "\n"+
		"export class {{ $messageName }} {\n"+
		"{{ $messageCode_Fields }}" +
		"}\n"
	res = strings.ReplaceAll(res, "{{ $messageName }}", message.Name)
	res = strings.ReplaceAll(res, "{{ $messageCode_Fields }}", fieldsStr)

	return res
}
