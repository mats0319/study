package data

var Copyright = []byte("// Generate File, Should not Edit.\n" +
	"// Author: mario. https://github.com/mats9693\n" +
	"// Version: " + Version + "\n")
var Version = "goc_ts v0.3.0"

var DefaultAPIConfig = &APIConfig{
	Dir:                   "./go/",
	OutDir:                "./ts/",
	BaseURL:               "http://127.0.0.1:9693",
	Timeout:               3_000,
	RequestMessageSuffix:  "Req",
	ResponseMessageSuffix: "Res",
	ServiceFileSuffix:     ".http.ts",
	MessageFileSuffix:     ".go.ts",
	BasicGoType: []struct {
		GoType      []string `json:"go_type"`
		TsType      string   `json:"ts_type"`
		TsZeroValue string   `json:"ts_zero_value"`
	}{
		{
			GoType:      []string{"string"},
			TsType:      "string",
			TsZeroValue: `""`,
		},
		{
			GoType:      []string{"int", "int8", "int32", "int64"},
			TsType:      "number",
			TsZeroValue: `0`,
		},
		{
			GoType:      []string{"bool"},
			TsType:      "boolean",
			TsZeroValue: `false`,
		},
	},
	Indentation: 4,
}
