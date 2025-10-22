package generate_ts

import (
	"github.com/mats9693/study/go/goc-ts/data"
	"github.com/mats9693/study/go/goc-ts/utils"
)

func GenerateUtils() {
	if !data.GeneratorIns.Utils.NeedObjectToFormData { // if all bool variables are false, not create 'utils' file
		return
	}

	content := utils.Copyright

	// functions
	if data.GeneratorIns.Utils.NeedObjectToFormData {
		content = append(content, data.GeneratorIns.Utils.ObjectToFormData...)
	}

	utils.WriteFile(data.GeneratorIns.Config.TsDir+"utils.ts", content)
}
