package data

type API struct {
	Config  *APIConfig
	Utils   *APIUtils
	Service map[string][]*ServiceItem // filename - service items
	Message map[string][]*MessageItem // filename - message items
}

type APIConfig struct {
	BaseURL string `json:"base_url"`
	Timeout int64  `json:"timeout"` // unit: micro-second
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
	TSType      string
	TSZeroValue string
	IsArray     bool
}
