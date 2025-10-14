package initialize

type Initializer struct {
	Files []*GoAPIFile `json:"files"`
}

type GoAPIFile struct {
	FileName    string     `json:"file_name"`
	PackageName string     `json:"package_name"`
	APIList     []*APIItem `json:"api_list"`
}

type APIItem struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

var DefaultInitializer = &Initializer{
	Files: []*GoAPIFile{
		{
			FileName:    "demo",
			PackageName: "api",
			APIList: []*APIItem{
				{
					Name: "ListUser",
					URI:  "/user/list",
				},
				{
					Name: "CreateUser",
					URI:  "/user/create",
				},
			},
		},
	},
}
