package generate_ts

import (
	"strconv"
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
	"github.com/mats9693/study/go/goc_ts/utils"
)

func GenerateConfig() {
	content := utils.Copyright
	content = append(content, []byte(serializeConfigCode())...)

	utils.WriteFile(data.GeneratorIns.Config.TsDir+"config.ts", content)
}

func serializeConfigCode() string {
	res := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
{{ $indentation }}baseURL: "{{ $baseURL }}",
{{ $indentation }}timeout: {{ $timeout }},
});
`
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(data.GetIndentation()))
	res = strings.ReplaceAll(res, "{{ $baseURL }}", data.GeneratorIns.Config.BaseURL)
	res = strings.ReplaceAll(res, "{{ $timeout }}", strconv.Itoa(int(data.GeneratorIns.Config.Timeout)))

	return res
}
