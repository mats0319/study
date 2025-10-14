package generate_ts

import (
	"fmt"
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func GenerateServiceFiles() {
	for filename := range data.GeneratorIns.Services {
		content := utils.Copyright
		content = append(content, serializeServiceFile(filename)...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.ServiceFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeServiceFile(filename string) string {
	res := `
import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { {{ $structures }} } from "./{{ $filename }}.go"
{{ $importExternalFunctions }}

class {{ $filenameBig }}Axios {{{ $serviceCode_Requests }}}

export const {{ $filename }}Axios: {{ $filenameBig }}Axios = new {{ $filenameBig }}Axios()
`

	structuresStr, utilsStr, requestStr := prepareData(filename)

	res = strings.ReplaceAll(res, "{{ $structures }}", structuresStr)
	res = strings.ReplaceAll(res, "{{ $filename }}", filename)
	res = strings.ReplaceAll(res, "{{ $importExternalFunctions }}", utilsStr)
	res = strings.ReplaceAll(res, "{{ $filenameBig }}", utils.MustBig(filename))
	res = strings.ReplaceAll(res, "{{ $serviceCode_Requests }}", requestStr)
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation()))

	return res
}

// prepareData prepare 'structures' / 'import functions' / 'http request invokes' of service code
func prepareData(filename string) (string, string, string) {
	var (
		structures   []string
		functions    = make(map[string]struct{}) // key: func name
		httpRequests []byte
	)

	for i := range data.GeneratorIns.Services[filename] {
		serviceItemIns := data.GeneratorIns.Services[filename][i]

		messageFields := make([]*data.MessageField, 0)
		messageReqName := serviceItemIns.Name + data.GeneratorIns.Config.RequestMessageSuffix
		for j := range data.GeneratorIns.Messages[filename] {
			if data.GeneratorIns.Messages[filename][j].Name == messageReqName { // find 'xxxReq' message
				messageFields = data.GeneratorIns.Messages[filename][j].Fields
				break
			}
		}

		structures = append(structures, serviceItemIns.Name+data.GeneratorIns.Config.ResponseMessageSuffix) // import 'xxxRes' message

		if len(messageFields) > 0 { // a http request need input param(s)
			structures = append(structures, messageReqName)

			functions[utils.FunctionName_ObjectToFormData] = struct{}{} // if 'xxxReq' has field(s), need this func
		}

		httpRequests = append(httpRequests, serializeOneHttpRequestInvoke(serviceItemIns, messageFields)...)
	}

	return utils.FormatStrSliceInLine(structures), serializeExternalFunctionsImport(functions), string(httpRequests)
}

// serializeExternalFunctionsImport serialize import external functions statement
func serializeExternalFunctionsImport(functions map[string]struct{}) string {
	if len(functions) < 1 {
		return ""
	}

	funcSlice := make([]string, 0)
	for funcName := range functions {
		funcSlice = append(funcSlice, funcName)
	}

	// 因为可能有的http请求文件不需要引入外部函数，所以这个引用的整个格式都写在这里
	res := `import { {{ $functionNames }} } from "./utils"`
	res = strings.ReplaceAll(res, "{{ $functionNames }}", utils.FormatStrSliceInLine(funcSlice))

	return res
}

// serializeOneHttpRequestInvoke serialize a http request invoke
func serializeOneHttpRequestInvoke(serviceItem *data.ServiceItem, messageFields []*data.MessageField) string {
	res := "\n{{ $indentation }}public {{ $serviceNameSmall }}({{ $functionInputs }}): " +
		"Promise<AxiosResponse<{{ $serviceName }}Res>> {{{ $requestInputs }}\n" +
		"{{ $indentation }}{{ $indentation }}return axiosWrapper.post(\"{{ $serviceURI }}\"{{ $invokeInputs }})\n" +
		"{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $serviceNameSmall }}", utils.MustSmall(serviceItem.Name))
	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceItem.Name)
	res = strings.ReplaceAll(res, "{{ $serviceURI }}", serviceItem.URI)

	var (
		functionInputs string
		requestInputs  string
		invokeInputs   string
	)
	if len(messageFields) > 0 { // http request has input param(s)
		// format fields to 'ts input params' and 'ts object init'
		functionInputsSlice := make([]string, 0)
		requestInputsSlice := make([]string, 0)
		for i := range messageFields {
			functionInputsSlice = append(functionInputsSlice, fmt.Sprintf("%s: %s", messageFields[i].Name, messageFields[i].TSType))
			requestInputsSlice = append(requestInputsSlice, fmt.Sprintf("%s: %s,\n", messageFields[i].Name, messageFields[i].Name))
		}

		functionInputs = utils.FormatStrSliceInLine(functionInputsSlice)
		requestInputs = serializeOneHttpRequestInput(serviceItem.Name, requestInputsSlice)
		invokeInputs = ", " + utils.FunctionName_ObjectToFormData + "(req)"
	}

	res = strings.ReplaceAll(res, "{{ $functionInputs }}", functionInputs)
	res = strings.ReplaceAll(res, "{{ $requestInputs }}", requestInputs)
	res = strings.ReplaceAll(res, "{{ $invokeInputs }}", invokeInputs)

	return res
}

func serializeOneHttpRequestInput(serviceName string, requestInputsSlice []string) string {
	res := "\n{{ $indentation }}{{ $indentation }}let req: {{ $serviceName }}Req = {\n"

	for i := range requestInputsSlice {
		res = res + "{{ $indentation }}{{ $indentation }}{{ $indentation }}" + requestInputsSlice[i]
	}

	res += "{{ $indentation }}{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceName)

	return res
}
