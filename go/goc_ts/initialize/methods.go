package initialize

import (
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func (ins *GoAPIFile) writeFile(packageName string) {
	content := "package {{ $packageName }}\n"
	content = strings.ReplaceAll(content, "{{ $packageName }}", packageName)

	for i := range ins.APIList {
		content += ins.APIList[i].toGo()
	}

	ins.FileName = utils.MustSuffix(ins.FileName, ".go")
	utils.WriteFile(data.GeneratorIns.Config.GoDir+ins.FileName, []byte(content))
}

func (ins *APIItem) toGo() string {
	res := `
const URI_{{ $apiName }} = "{{ $apiURI }}"

type {{ $apiName }}{{ $reqSuffix }} struct {}

type {{ $apiName }}{{ $resSuffix }} struct {}
`

	res = strings.ReplaceAll(res, "{{ $apiName }}", ins.Name)
	res = strings.ReplaceAll(res, "{{ $apiURI }}", ins.URI)
	res = strings.ReplaceAll(res, "{{ $reqSuffix }}", data.GeneratorIns.Config.RequestStructureSuffix)
	res = strings.ReplaceAll(res, "{{ $resSuffix }}", data.GeneratorIns.Config.ResponseStructureSuffix)

	return res
}
