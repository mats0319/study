package parse

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"strings"
)

func ParseUtils(apiIns *data.API) {
ALL:
	for filename, serviceItems := range apiIns.Service {
		for i := range serviceItems {
			messageName := serviceItems[i].Name + data.RequestMessageSuffix

			hasValidMessage := false
			for j := range apiIns.Message[filename] {
				messageItemIns := apiIns.Message[filename][j]
				if messageItemIns.Name == messageName && len(messageItemIns.Fields) > 0 { // exist 'xxxReq' message, and not empty
					hasValidMessage = true
					break
				}
			}

			if hasValidMessage {
				apiIns.Utils.NeedObjectToFormData = true
				apiIns.Utils.ObjectToFormData = []byte(funcCodeIndentation(data.FunctionCode_ObjectToFormData))
				break ALL
			}
		}
	}
}

func funcCodeIndentation(funcCode string) string {
	res := funcCode
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation(1)))

	return res
}
