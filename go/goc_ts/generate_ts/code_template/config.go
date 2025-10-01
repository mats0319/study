package code_template

import (
	"strconv"
	"strings"

	"github.com/mats9693/study/go/goc_ts/data"
)

func FormatConfigCode(config *data.APIConfig) string {
	res := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
{{ $indentation }}baseURL: "{{ $baseURL }}",
{{ $indentation }}timeout: {{ $timeout }},
});
`
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(config.GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $baseURL }}", config.BaseURL)
	res = strings.ReplaceAll(res, "{{ $timeout }}", strconv.Itoa(int(config.Timeout)))

	return res
}
