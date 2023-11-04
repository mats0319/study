package data

import "strings"

const functionName_ObjectToFormData = "objectToFormData"

const FunctionCode_ObjectToFormData = `
// objectToFormData 泛型用于解决'obj[key]'报错问题
export function objectToFormData<T extends object>(obj: T): FormData {
{{ $indentation }}let data: FormData = new FormData()
{{ $indentation }}for (let key in obj) {
{{ $indentation }}{{ $indentation }}if (typeof obj[key] == "object") { // if field type is another object
{{ $indentation }}{{ $indentation }}{{ $indentation }}objectToFormData(obj[key] as object).forEach((value: FormDataEntryValue, key: string) => {
{{ $indentation }}{{ $indentation }}{{ $indentation }}{{ $indentation }}data.append(key, value)
{{ $indentation }}{{ $indentation }}{{ $indentation }}})
{{ $indentation }}{{ $indentation }}} else { // normal
{{ $indentation }}{{ $indentation }}{{ $indentation }}data.append(key, obj[key] as string)
{{ $indentation }}{{ $indentation }}}
{{ $indentation }}}

{{ $indentation }}return data
}
`

func FuncCodeIndentation(funcCode string) string {
	res := funcCode
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(GetIndentation(1)))

	return res
}
