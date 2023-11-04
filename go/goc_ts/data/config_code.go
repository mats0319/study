package data

import (
	"strconv"
	"strings"
)

func FormatConfigCode(config *APIConfig) string {
	res := `
import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
{{ $indentation }}baseURL: "{{ $baseURL }}",
{{ $indentation }}timeout: {{ $timeout }},
});
`
	res = strings.ReplaceAll(res, "{{ $indentation }}", string(GetIndentation(1)))
	res = strings.ReplaceAll(res, "{{ $baseURL }}", config.BaseURL)
	res = strings.ReplaceAll(res, "{{ $timeout }}", strconv.Itoa(int(config.Timeout)))

	return res
}
