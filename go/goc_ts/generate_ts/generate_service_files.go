package generate_ts

import (
	"fmt"
	"github.com/mats9693/study/go/goc_ts/data"
	"log"
	"os"
	"strings"
)

func GenerateServiceFiles(apiIns *data.API, outDir string) {
	for filename := range apiIns.Service {
		absolutePath := outDir + filename + data.ServiceFileSuffix
		generateServiceFile(absolutePath, apiIns.Service[filename], apiIns.Message[filename], filename)
	}
}

func generateServiceFile(
	absolutePath string,
	serviceItems []*data.ServiceItem,
	messageItems []*data.MessageItem,
	filename string,
) {
	file, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open service file(%s) failed, error: %v\n", absolutePath, err))
	}
	defer func() {
		_ = file.Close()
	}()

	structures, utils, requests := formatData(serviceItems, messageItems)

	content := data.Copyright
	content = append(content, fillServiceCode(structures, filename, utils, requests)...)

	_, err = file.Write(content)
	if err != nil {
		log.Fatalln(fmt.Sprintf("write service file(%s) failed, error: %v\n", absolutePath, err))
	}
}

// formatData format 'structures'/'utils'/'requests' data
func formatData(serviceItems []*data.ServiceItem, messageItems []*data.MessageItem) (string, string, string) {
	var (
		structures []string
		functions  = make(map[string]struct{})
		requests   []byte
	)

	for i := range serviceItems {
		serviceItemIns := serviceItems[i]

		messageFields := make([]*data.MessageField, 0)
		messageReqName := serviceItemIns.Name + data.RequestMessageSuffix
		for j := range messageItems {
			if messageItems[j].Name == messageReqName { // find 'xxxReq' message
				messageFields = messageItems[j].Fields
				break
			}
		}

		structures = append(structures, serviceItemIns.Name+data.ResponseMessageSuffix) // import 'xxxRes' message

		if len(messageFields) < 1 { // 'xxxReq' message not exist or empty
			requests = append(requests, fillServiceCodeRequest(
				serviceItemIns.Name,
				serviceItemIns.URI,
				"",
				"",
				"",
			)...)
		} else { // 'xxxReq' message exist and has one or more field(s)
			structures = append(structures, messageReqName)

			functions["objectToFormData"] = struct{}{} // function name, use in 'utils' part

			// format fields to 'ts input params' and 'ts object init'
			paramsWithType := make([]string, 0)
			paramsWithInit := make([]string, 0)
			for j := range messageFields {
				paramsWithType = append(paramsWithType, fmt.Sprintf("%s: %s", messageFields[j].Name, messageFields[j].TSType))
				paramsWithInit = append(paramsWithInit, fmt.Sprintf("%s: %s,\n", messageFields[j].Name, messageFields[j].Name))
			}

			requests = append(requests, fillServiceCodeRequest(
				serviceItemIns.Name,
				serviceItemIns.URI,
				", objectToFormData(req)",
				formatStrSliceInLine(paramsWithType),
				fillServiceCodeReqStruct(serviceItemIns.Name, formatStrSliceInLine(paramsWithInit)),
			)...)
		}
	}

	return formatStrSliceInLine(structures), formatServiceUtils(functions), string(requests)
}

func fillServiceCode(
	structures string,
	filename string,
	serviceCode_Utils string,
	serviceCode_Request string,
) string {
	res := data.ServiceCode_Template
	res = strings.ReplaceAll(res, "{{ $structures }}", structures)
	res = strings.ReplaceAll(res, "{{ $filename }}", filename)
	res = strings.ReplaceAll(res, "{{ $serviceCode_Utils }}", serviceCode_Utils)
	res = strings.ReplaceAll(res, "{{ $filenameBig }}", mustBig(filename))
	res = strings.ReplaceAll(res, "{{ $serviceCode_Requests }}", serviceCode_Request)

	return res
}

func fillServiceCodeRequest(
	serviceName string,
	serviceURI string,
	requestParams string,
	paramsWithType string,
	serviceCode_ReqStruct string,
) string {
	res := data.ServiceCode_Request
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $serviceNameSmall }}", mustSmall(serviceName))
	res = strings.ReplaceAll(res, "{{ $paramsWithType }}", paramsWithType)
	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceName)
	res = strings.ReplaceAll(res, "{{ $serviceCode_ReqStruct }}", serviceCode_ReqStruct)
	res = strings.ReplaceAll(res, "{{ $serviceURI }}", serviceURI)
	res = strings.ReplaceAll(res, "{{ $requestParams }}", requestParams)

	return res
}

func fillServiceCodeReqStruct(serviceName string, paramsWithInit string) string {
	res := data.ServiceCode_ReqStruct
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $serviceName }}", serviceName)
	res = strings.ReplaceAll(res, "{{ $paramsWithInit }}", paramsWithInit)

	return res
}

func formatServiceUtils(functions map[string]struct{}) string {
	res := ""
	if len(functions) > 0 {
		funcSlice := make([]string, 0)
		for funcName := range functions {
			funcSlice = append(funcSlice, funcName)
		}

		res = fmt.Sprintf(`import { %s } from "./utils"`, formatStrSliceInLine(funcSlice))
	}

	return res
}

// formatStrSliceInLine e.g. ["a","b","c"] => "a, b, c"
func formatStrSliceInLine(data []string) string {
	var resBytes []byte
	for i := range data {
		resBytes = append(resBytes, ", "+data[i]...)
	}

	if len(resBytes) > 2 {
		resBytes = resBytes[2:]
	}

	return string(resBytes)
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
