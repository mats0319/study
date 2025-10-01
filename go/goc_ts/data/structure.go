package data

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type API struct {
	Config  *APIConfig
	Utils   *APIUtils
	Service map[string][]*ServiceItem // filename - service items
	Message map[string][]*MessageItem // filename - message items

	TsType      map[string]string // go type - ts type
	TsZeroValue map[string]string // ts type - ts zero value
}

func (ins *API) SetConfigFromFile(configFile string) {
	file, err := os.Open(configFile)
	if err != nil { // no valid config file, just use default value
		log.Println("open config file failed, use default value, error: ", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln("get config file stat failed, error: ", err)
	}

	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Fatalln("read config file failed, error: ", err)
	}

	apiConfigIns := &APIConfig{}
	err = json.Unmarshal(fileBytes, apiConfigIns)
	if err != nil {
		log.Fatalln("deserialize api config failed, error: ", err)
	}

	// use default config cover empty config
	if len(apiConfigIns.Dir) < 1 {
		apiConfigIns.Dir = DefaultAPIConfig.Dir
	}
	if len(apiConfigIns.OutDir) < 1 {
		apiConfigIns.OutDir = DefaultAPIConfig.OutDir
	}
	if len(apiConfigIns.BaseURL) < 1 {
		apiConfigIns.BaseURL = DefaultAPIConfig.BaseURL
	}
	if apiConfigIns.Timeout < 1 {
		apiConfigIns.Timeout = DefaultAPIConfig.Timeout
	}
	if len(apiConfigIns.RequestMessageSuffix) < 1 {
		apiConfigIns.RequestMessageSuffix = DefaultAPIConfig.RequestMessageSuffix
	}
	if len(apiConfigIns.ResponseMessageSuffix) < 1 {
		apiConfigIns.ResponseMessageSuffix = DefaultAPIConfig.ResponseMessageSuffix
	}
	if len(apiConfigIns.ServiceFileSuffix) < 1 {
		apiConfigIns.ServiceFileSuffix = DefaultAPIConfig.ServiceFileSuffix
	}
	if len(apiConfigIns.MessageFileSuffix) < 1 {
		apiConfigIns.MessageFileSuffix = DefaultAPIConfig.MessageFileSuffix
	}
	if len(apiConfigIns.BasicGoType) < 1 {
		apiConfigIns.BasicGoType = DefaultAPIConfig.BasicGoType
	}
	if apiConfigIns.Indentation < 1 {
		apiConfigIns.Indentation = DefaultAPIConfig.Indentation
	}

	ins.Config = apiConfigIns

	mustDir(ins.Config.Dir)
	mustDir(ins.Config.OutDir)

	// set type maps
	for i, typ := range ins.Config.BasicGoType {
		for _, goTyp := range ins.Config.BasicGoType[i].GoType {
			ins.TsType[goTyp] = typ.TsType
		}
		ins.TsZeroValue[typ.TsType] = typ.TsZeroValue
	}
}

type APIConfig struct {
	// work path
	Dir    string `json:"dir"`
	OutDir string `json:"out_dir"`
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

func (c *APIConfig) GetIndentation(level int) []byte {
	res := make([]byte, 0, c.Indentation*level)
	for i := 0; i < c.Indentation*level; i++ {
		res = append(res, ' ')
	}
	return res
}

type APIUtils struct {
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

// mustDir make sure 'path' is end with '/'
func mustDir(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}
