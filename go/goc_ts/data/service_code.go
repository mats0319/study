package data

import (
	"fmt"
	"strings"
)

func FormatServiceCode(serviceItems []*ServiceItem, messageItems []*MessageItem, filename string) string {
	res := `
import { axiosWrapper } from "./config"
import { {{ $structures }} } from "./{{ $filename }}.go"
{{ $serviceCode_Utils }}

class {{ $filenameBig }}Axios {{{ $serviceCode_Requests }}}

export const {{ $filename }}Axios: {{ $filenameBig }}Axios = new {{ $filenameBig }}Axios()
`

	structuresStr, utilsStr, requestStr := prepareData(serviceItems, messageItems)

	res = strings.ReplaceAll(res, "{{ $structures }}", structuresStr)
	res = strings.ReplaceAll(res, "{{ $filename }}", filename)
	res = strings.ReplaceAll(res, "{{ $serviceCode_Utils }}", utilsStr)
	res = strings.ReplaceAll(res, "{{ $filenameBig }}", mustBig(filename))
	res = strings.ReplaceAll(res, "{{ $serviceCode_Requests }}", requestStr)

	return res
}

// prepareData prepare 'structures' / 'utils' / 'requests' of service code
func prepareData(serviceItems []*ServiceItem, messageItems []*MessageItem) (string, string, string) {
	var (
		structures []string
		functions  = make(map[string]struct{}) // key: func name
		requests   []byte
	)

	for i := range serviceItems {
		serviceItemIns := serviceItems[i]

		messageFields := make([]*MessageField, 0)
		messageReqName := serviceItemIns.Name + RequestMessageSuffix
		for j := range messageItems {
			if messageItems[j].Name == messageReqName { // find 'xxxReq' message
				messageFields = messageItems[j].Fields
				break
			}
		}

		structures = append(structures, serviceItemIns.Name+ResponseMessageSuffix) // import 'xxxRes' message

		if len(messageFields) > 0 { // 'xxxReq' message exist and has one or more field(s)
			structures = append(structures, messageReqName)

			functions[functionName_ObjectToFormData] = struct{}{} // if 'xxxReq' has field(s), need this func
		}

		requests = append(requests, fillServiceCodeRequest(
			serviceItemIns.Name,
			serviceItemIns.URI,
			messageFields,
		)...)
	}

	return formatStrSliceInLine(structures), formatServiceUtils(functions), string(requests)
}

func formatServiceUtils(functions map[string]struct{}) string {
	res := ""
	if len(functions) > 0 {
		funcSlice := make([]string, 0)
		for funcName := range functions {
			funcSlice = append(funcSlice, funcName)
		}

		res = `import { {{ $functionNames }} } from "./utils"`
		res = strings.ReplaceAll(res, "{{ $functionNames }}", formatStrSliceInLine(funcSlice))
	}

	return res
}

func fillServiceCodeRequest(serviceName string, serviceURI string, messageFields []*MessageField) string {
	res := "\n" +
		"{{ $indentation }}public {{ $serviceNameSmall }}({{ $paramsWithType }}): " +
		"Promise<{{ $serviceName }}Res> {{{ $serviceCode_ReqStruct }}\n" +
		"{{ $indentation }}{{ $indentation }}return axiosWrapper.post(\"{{ $serviceURI }}\"{{ $requestParams }})\n" +
		"{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $indentation }}", string(GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $serviceNameSmall }}", mustSmall(serviceName))
	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceName)
	res = strings.ReplaceAll(res, "{{ $serviceURI }}", serviceURI)

	var (
		reqParamsStr string
		reqStructStr string
		paramsWithTypeStr string
	)
	if len(messageFields) > 0 { // message has field(s)
		// format fields to 'ts input params' and 'ts object init'
		paramsWithType := make([]string, 0)
		paramsWithInit := make([]string, 0)
		for i := range messageFields {
			paramsWithType = append(paramsWithType, fmt.Sprintf("%s: %s", messageFields[i].Name, messageFields[i].TSType))
			paramsWithInit = append(paramsWithInit, fmt.Sprintf("%s: %s,\n", messageFields[i].Name, messageFields[i].Name))
		}

		reqParamsStr = ", " + functionName_ObjectToFormData + "(req)"
		reqStructStr = fillServiceCodeReqStruct(serviceName, formatStrSliceInLine(paramsWithInit))
		paramsWithTypeStr = formatStrSliceInLine(paramsWithType)
	}

	res = strings.ReplaceAll(res, "{{ $requestParams }}", reqParamsStr)
	res = strings.ReplaceAll(res, "{{ $paramsWithType }}", paramsWithTypeStr)
	res = strings.ReplaceAll(res, "{{ $serviceCode_ReqStruct }}", reqStructStr)

	return res
}

func fillServiceCodeReqStruct(serviceName string, paramsWithInit string) string {
	res := "\n" +
	"{{ $indentation }}{{ $indentation }}let req: {{ $serviceName }}Req = {\n" +
	"{{ $indentation }}{{ $indentation }}{{ $indentation }}{{ $paramsWithInit }}" +
	"{{ $indentation }}{{ $indentation }}}\n"

	res = strings.ReplaceAll(res, "{{ $indentation }}", string(GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceName)
	res = strings.ReplaceAll(res, "{{ $paramsWithInit }}", paramsWithInit)

	return res
}

// formatStrSliceInLine e.g. ["a","b","c"] => "a, b, c"
func formatStrSliceInLine(data []string) string {
	if len(data) < 1 {
		return ""
	}

	var resBytes []byte
	for i := range data {
		resBytes = append(resBytes, ", "+data[i]...)
	}

	return string(resBytes[2:])
}

// mustSmall make sure first char is small-case, e.g. "GetData" => "getData"
func mustSmall(str string) string {
	if len(str) > 0 && 'A' < str[0] && str[0] < 'Z' {
		str = string(str[0]-'A'+'a') + str[1:]
	}

	return str
}

// mustBig make sure first char is big-case, e.g. "getData" => "DetData"
func mustBig(str string) string {
	if len(str) > 0 && 'a' < str[0] && str[0] < 'z' {
		str = string(str[0]-'a'+'A') + str[1:]
	}

	return str
}
