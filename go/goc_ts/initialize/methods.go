package initialize

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func (ins *GoAPIFile) toGo() {
	content := "package {{ $PackageName }}\n"
	content = strings.ReplaceAll(content, "{{ $PackageName }}", ins.PackageName)

	for i := range ins.APIList {
		content += ins.APIList[i].toGo()
	}

	ins.FileName = utils.MustGoFileName(ins.FileName)
	utils.WriteFile(data.GeneratorIns.Config.GoDir+ins.FileName, []byte(content))

	return
}

func (ins *APIItem) toGo() string {
	res := `
const URI_{{ $APIName }} = "{{ $APIURI }}"

type {{ $APIName }}{{ $ReqSuffix }} struct {}

type {{ $APIName }}{{ $ResSuffix }} struct {}
`

	res = strings.ReplaceAll(res, "{{ $APIName }}", ins.Name)
	res = strings.ReplaceAll(res, "{{ $APIURI }}", ins.URI)
	res = strings.ReplaceAll(res, "{{ $ReqSuffix }}", data.GeneratorIns.Config.RequestMessageSuffix)
	res = strings.ReplaceAll(res, "{{ $ResSuffix }}", data.GeneratorIns.Config.ResponseMessageSuffix)

	return res
}
