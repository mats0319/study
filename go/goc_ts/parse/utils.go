package parse

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func ParseUtils() {
ALL:
	for filename, requestNames := range data.GeneratorIns.RequestAffiliation {
		for _, requestName := range requestNames {
			reqStructureName := requestName + data.GeneratorIns.Config.RequestStructureSuffix

			for _, structureName := range data.GeneratorIns.StructureAffiliation[filename] {
				if structureName != reqStructureName {
					continue
				}

				structureItemIns, ok := data.GeneratorIns.Structures[structureName]
				if ok && len(structureItemIns.Fields) > 0 { // exist 'xxxReq' message, and not empty
					data.GeneratorIns.Utils.NeedObjectToFormData = true
					data.GeneratorIns.Utils.ObjectToFormData = []byte(funcCodeIndentation(utils.FunctionCode_ObjectToFormData))
					break ALL
				}
			}
		}
	}
}

func funcCodeIndentation(funcCode string) string {
	res := funcCode
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation()))

	return res
}
