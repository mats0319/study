package parse

import (
	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/generate_ts/code_template"
)

func ParseUtils(apiIns *data.API) {
ALL:
	for filename, serviceItems := range apiIns.Service {
		for i := range serviceItems {
			messageName := serviceItems[i].Name + apiIns.Config.RequestMessageSuffix

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
				apiIns.Utils.ObjectToFormData = []byte(code_template.FuncCodeIndentation(apiIns.Config, code_template.FunctionCode_ObjectToFormData))
				break ALL
			}
		}
	}
}
