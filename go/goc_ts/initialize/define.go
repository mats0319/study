package initialize

type Initializer struct {
	PackageName string       `json:"package_name"` // 因为同文件夹下.go文件的包名应统一，这里就提出来了
	Files       []*GoAPIFile `json:"files"`
}

var InitializerIns = &Initializer{}

type GoAPIFile struct {
	FileName string     `json:"file_name"`
	APIList  []*APIItem `json:"api_list"`
}

type APIItem struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

var DefaultInitializer = &Initializer{
	PackageName: "api",
	Files: []*GoAPIFile{
		{
			FileName: "demo",
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
