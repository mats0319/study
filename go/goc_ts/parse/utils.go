package parse

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func ParseUtils() {
ALL:
	for filename, serviceItems := range data.GeneratorIns.Services {
		for i := range serviceItems {
			messageName := serviceItems[i].Name + data.GeneratorIns.Config.RequestMessageSuffix

			hasValidMessage := false
			for j := range data.GeneratorIns.Messages[filename] {
				messageItemIns := data.GeneratorIns.Messages[filename][j]
				if messageItemIns.Name == messageName && len(messageItemIns.Fields) > 0 { // exist 'xxxReq' message, and not empty
					hasValidMessage = true
					break
				}
			}

			if hasValidMessage {
				data.GeneratorIns.Utils.NeedObjectToFormData = true
				data.GeneratorIns.Utils.ObjectToFormData = []byte(funcCodeIndentation(utils.FunctionCode_ObjectToFormData))
				break ALL
			}
		}
	}
}

func funcCodeIndentation(funcCode string) string {
	res := funcCode
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation()))

	return res
}
