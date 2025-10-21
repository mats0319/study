package generate_ts

import (
	"fmt"
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
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
	res := `
import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
{{ $structures }}{{ $externalFuncs }}

class {{ $filenameBig }}Axios {{{ $requests }}}

export const {{ $filename }}Axios: {{ $filenameBig }}Axios = new {{ $filenameBig }}Axios()
`

	structuresStr, externalFuncs, requestStr := prepareData(filename)

	res = strings.ReplaceAll(res, "{{ $structures }}", structuresStr)
	res = strings.ReplaceAll(res, "{{ $filename }}", filename)
	res = strings.ReplaceAll(res, "{{ $externalFuncs }}", externalFuncs)
	res = strings.ReplaceAll(res, "{{ $filenameBig }}", utils.MustBig(filename))
	res = strings.ReplaceAll(res, "{{ $requests }}", requestStr)
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation()))

	return res
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
					fromFile, _ := data.GeneratorIns.TypeFrom[structureField.GoType]
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
	res := "\n{{ $indentation }}public {{ $requestNameSmall }}({{ $functionInputs }}): " +
		"Promise<AxiosResponse<{{ $requestName }}{{ $responseStructureSuffix }}>> {{{ $requestInputs }}\n" +
		"{{ $indentation }}{{ $indentation }}return axiosWrapper.post(\"{{ $requestURI }}\"{{ $invokeInputs }})\n" +
		"{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $requestNameSmall }}", utils.MustSmall(requestName))
	res = strings.ReplaceAll(res, "{{ $requestName }}", requestName)
	res = strings.ReplaceAll(res, "{{ $responseStructureSuffix }}", data.GeneratorIns.Config.ResponseStructureSuffix)
	res = strings.ReplaceAll(res, "{{ $requestURI }}", data.GeneratorIns.Requests[requestName])

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

	res = strings.ReplaceAll(res, "{{ $functionInputs }}", functionInputs)
	res = strings.ReplaceAll(res, "{{ $requestInputs }}", requestInputs)
	res = strings.ReplaceAll(res, "{{ $invokeInputs }}", invokeInputs)

	return res
}

func serializeOneHttpRequestInput(requestName string, requestInputs []string) string {
	res := "\n{{ $indentation }}{{ $indentation }}let req: {{ $requestName }}{{ $requestMessageSuffix }} = {\n"

	for i := range requestInputs {
		res = res + "{{ $indentation }}{{ $indentation }}{{ $indentation }}" + requestInputs[i]
	}

	res += "{{ $indentation }}{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $requestName }}", requestName)
	res = strings.ReplaceAll(res, "{{ $requestMessageSuffix }}", data.GeneratorIns.Config.RequestStructureSuffix)

	return res
}
