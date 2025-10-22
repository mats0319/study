package generate_ts

import (
	"fmt"
	"strings"

	"github.com/mats9693/study/go/goc-ts/data"
	"github.com/mats9693/study/go/goc-ts/utils"
)

func GenerateRequestFiles() {
	for filename := range data.GeneratorIns.RequestAffiliation {
		content := utils.Copyright
		content = append(content, serializeRequestFile(filename)...)

		absolutePath := data.GeneratorIns.Config.TsDir + filename + data.GeneratorIns.Config.RequestFileSuffix
		utils.WriteFile(absolutePath, content)
	}
}

func serializeRequestFile(filename string) string {
	reqFileStr := `
import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
{{ $structures }}{{ $externalFuncs }}

class {{ $filenameBig }}Axios {{{ $requests }}}

export const {{ $filename }}Axios: {{ $filenameBig }}Axios = new {{ $filenameBig }}Axios()
`

	structuresStr, externalFuncs, requestStr := prepareData(filename)

	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $structures }}", structuresStr)
	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $filename }}", filename)
	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $externalFuncs }}", externalFuncs)
	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $filenameBig }}", utils.MustBig(filename))
	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $requests }}", requestStr)
	reqFileStr = strings.ReplaceAll(reqFileStr, "{{ $indentation }}", data.GeneratorIns.IndentationStr)

	return reqFileStr
}

// prepareData prepare 'structures' / 'import functions' / 'http request invokes' of service code
func prepareData(filename string) (string, string, string) {
	var (
		externalStructures = make(map[string][]string) // from filename - structures' name
		functions          = make(map[string]struct{}) // key: func name
		requests           = ""
	)

	for _, requestName := range data.GeneratorIns.RequestAffiliation[filename] {
		structureItemIns := &data.StructureItem{Fields: make([]*data.StructureField, 0)}

		reqStructureName := requestName + data.GeneratorIns.Config.RequestStructureSuffix
		for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
			if structureName == reqStructureName { // find 'xxxReq' message
				structureItemIns, _ = data.GeneratorIns.Structures[structureName]
				break
			}
		}

		responseResName := requestName + data.GeneratorIns.Config.ResponseStructureSuffix
		externalStructures[filename] = append(externalStructures[filename], responseResName)

		if len(structureItemIns.Fields) > 0 { // a http request need input param(s)
			externalStructures[filename] = append(externalStructures[filename], reqStructureName)

			functions[utils.FunctionName_ObjectToFormData] = struct{}{} // if 'xxxReq' has field(s), need this func

			for _, structureField := range structureItemIns.Fields {
				if _, ok := data.GeneratorIns.TsType[structureField.GoType]; !ok {
					fromFile, _ := data.GeneratorIns.StructureFrom[structureField.GoType]
					externalStructures[fromFile] = append(externalStructures[fromFile], structureField.TSType)
				}
			}
		}

		requests += serializeOneHttpRequestInvoke(requestName, structureItemIns)
	}

	return serializeStructuresImport(externalStructures), serializeExternalFunctionsImport(functions), requests
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
func serializeOneHttpRequestInvoke(requestName string, structureItemIns *data.StructureItem) string {
	httpReqInvokeStr := `
{{ $indentation }}public {{ $requestNameSmall }}({{ $functionInputs }}): Promise<AxiosResponse<{{ $requestName }}{{ $responseStructureSuffix }}>> {{{ $requestInputs }}
{{ $indentation }}{{ $indentation }}return axiosWrapper.post("{{ $requestURI }}"{{ $invokeInputs }})
{{ $indentation }}}
`

	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestNameSmall }}", utils.MustSmall(requestName))
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestName }}", requestName)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $responseStructureSuffix }}", data.GeneratorIns.Config.ResponseStructureSuffix)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestURI }}", data.GeneratorIns.Requests[requestName])

	var (
		functionInputs string
		requestInputs  string
		invokeInputs   string
	)
	if len(structureItemIns.Fields) > 0 { // http request has input param(s)
		// format fields to 'ts input params' and 'ts object init'
		functionInputsSlice := make([]string, 0)
		requestInputsSlice := make([]string, 0)
		for _, v := range structureItemIns.Fields {
			functionInputsSlice = append(functionInputsSlice, fmt.Sprintf("%s: %s", v.Name, v.TSType))
			requestInputsSlice = append(requestInputsSlice, fmt.Sprintf("%s: %s,\n", v.Name, v.Name))
		}

		functionInputs = utils.FormatStrSliceInLine(functionInputsSlice)
		requestInputs = serializeOneHttpRequestInput(requestName, requestInputsSlice)
		invokeInputs = ", " + utils.FunctionName_ObjectToFormData + "(req)"
	}

	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $functionInputs }}", functionInputs)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $requestInputs }}", requestInputs)
	httpReqInvokeStr = strings.ReplaceAll(httpReqInvokeStr, "{{ $invokeInputs }}", invokeInputs)

	return httpReqInvokeStr
}

func serializeOneHttpRequestInput(requestName string, requestInputs []string) string {
	fieldStr := ""
	for i := range requestInputs {
		field := "{{ $indentation }}{{ $indentation }}{{ $indentation }}{{ $requestInput }}"
		field = strings.ReplaceAll(field, "{{ $requestInput }}", requestInputs[i])

		fieldStr += field
	}

	httpReqInputStr := `
{{ $indentation }}{{ $indentation }}let req: {{ $requestName }}{{ $requestMessageSuffix }} = {
{{ $requestFields }}{{ $indentation }}{{ $indentation }}}
`
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestFields }}", fieldStr)
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestName }}", requestName)
	httpReqInputStr = strings.ReplaceAll(httpReqInputStr, "{{ $requestMessageSuffix }}", data.GeneratorIns.Config.RequestStructureSuffix)

	return httpReqInputStr
}
