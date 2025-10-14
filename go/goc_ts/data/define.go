package data

type Generator struct {
	Config   *Config
	Utils    *Utils
	Services map[string][]*ServiceItem // filename - service items
	Messages map[string][]*MessageItem // filename - message items

	TsType      map[string]string // go type - ts type
	TsZeroValue map[string]string // ts type - ts zero value
}

var GeneratorIns = &Generator{
	Config:      &Config{},
	Utils:       &Utils{},
	Services:    make(map[string][]*ServiceItem),
	Messages:    make(map[string][]*MessageItem),
	TsType:      make(map[string]string),
	TsZeroValue: make(map[string]string),
}

type Config struct {
	// work path
	GoDir string `json:"go_dir"`
	TsDir string `json:"ts_dir"`
	// axios config
	BaseURL string `json:"base_url"`
	Timeout int64  `json:"timeout"` // unit: micro-second
	// naming conventions 命名规范/约定
	RequestMessageSuffix  string `json:"request_message_suffix"`
	ResponseMessageSuffix string `json:"response_message_suffix"`
	ServiceFileSuffix     string `json:"service_file_suffix"`
	MessageFileSuffix     string `json:"message_file_suffix"`
	// params
	BasicGoType []struct {
		GoType      []string `json:"go_type"`
		TsType      string   `json:"ts_type"`
		TsZeroValue string   `json:"ts_zero_value"`
	} `json:"basic_go_type"`
	Indentation int `json:"indentation"`
}

type Utils struct {
	NeedObjectToFormData bool
	ObjectToFormData     []byte
}

type ServiceItem struct {
	Name string
	URI  string
}

type MessageItem struct {
	Name   string
	Fields []*MessageField
}

type MessageField struct {
	Name        string // field name, from json tag of go struct field
	GoType      string
	IsArray     bool
	TSType      string
	TSZeroValue string
}

var DefaultGeneratorConfig = &Config{
	GoDir:                 "./go/",
	TsDir:                 "./ts/",
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
