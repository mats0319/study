package generate_ts

import (
	"strconv"
	"strings"

	"github.com/mats9693/study/go/goc-ts/data"
	"github.com/mats9693/study/go/goc-ts/utils"
)

func GenerateConfig() {
	content := utils.Copyright
	content = append(content, serializeConfigFile()...)

	utils.WriteFile(data.GeneratorIns.Config.TsDir+"config.ts", content)
}

func serializeConfigFile() string {
	res := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
{{ $indentation }}baseURL: "{{ $baseURL }}",
{{ $indentation }}timeout: {{ $timeout }},
});
`
	res = strings.ReplaceAll(res, "{{ $indentation }}", data.GeneratorIns.IndentationStr)
	res = strings.ReplaceAll(res, "{{ $baseURL }}", data.GeneratorIns.Config.BaseURL)
	res = strings.ReplaceAll(res, "{{ $timeout }}", strconv.Itoa(data.GeneratorIns.Config.Timeout))

	return res
}
