package utils

var Copyright = []byte("// Generate File, Should not Edit.\n" +
	"// Author : mario. github.com/mats0319\n" +
	"// Code   : github.com/mats0319/study/go/goc_ts\n" +
	"// Version: " + Version + "\n")
var Version = "goc_ts v0.3.1"

const FunctionName_ObjectToFormData = "objectToFormData"

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
