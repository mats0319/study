package data

type Generator struct {
	Config *Config
	Utils  *Utils

	RequestAffiliation   map[string][]string       // filename - request name(s)
	Requests             map[string]string         // request name - request uri
	StructureAffiliation map[string][]string       // filename - structure name(s)
	Structures           map[string]*StructureItem // structure name - structure item
	StructureFrom        map[string]string         // (self-define) structure name - from filename

	TsType      map[string]string // go type - ts type
	TsZeroValue map[string]string // ts type - ts zero value
}

var GeneratorIns = &Generator{
	Config:               &Config{},
	Utils:                &Utils{},
	RequestAffiliation:   make(map[string][]string),
	Requests:             make(map[string]string),
	StructureAffiliation: make(map[string][]string),
	Structures:           make(map[string]*StructureItem),
	StructureFrom:        make(map[string]string),
	TsType:               make(map[string]string),
	TsZeroValue:          make(map[string]string),
}

type Config struct {
	// work path
	GoDir string `json:"go_dir"`
	TsDir string `json:"ts_dir"`
	// axios config
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"` // unit: micro-second
	// naming conventions 命名规范/约定
	RequestStructureSuffix  string `json:"request_structure_suffix"`
	ResponseStructureSuffix string `json:"response_structure_suffix"`
	RequestFileSuffix       string `json:"request_file_suffix"`
	StructureFileSuffix     string `json:"structure_file_suffix"`
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

type StructureItem struct {
	Fields []*StructureField
}

type StructureField struct {
	Name        string // field name, from json tag of go struct field
	GoType      string
	IsArray     bool
	TSType      string
	TSZeroValue string
}

var DefaultGeneratorConfig = &Config{
	GoDir:                   "./go/",
	TsDir:                   "./ts/",
	BaseURL:                 "http://127.0.0.1:9693",
	Timeout:                 3_000,
	RequestStructureSuffix:  "Req",
	ResponseStructureSuffix: "Res",
	RequestFileSuffix:       ".http.ts",
	StructureFileSuffix:     ".go.ts",
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
			GoType: []string{"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64"},
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
