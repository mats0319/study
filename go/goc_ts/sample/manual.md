# sample

```cmd
> go install github.com/mats9693/study/go/goc_ts
> goc_ts -h

2025/10/17 11:44:51 Options: 
(In this help, './go/' means our go files dir)
  -c string
        config file (default "./go/config.json")
  -g    generate go files from './go/init.json'
        
  -genFrom string
        generate go files from given file
  -h    this help
  -i    initializeFlag basic files
        overwrite './go/config_default.json' and './go/init_default.json'
  -v    show version
```

因为本工具有一些默认的约束条件（例如命名），所以我们为此提供了相应的支持（启动参数）：  
```cmd
goc_ts -i
```
命令可以生成一个**配置文件`(config_default.json)`**和一个**go接口文件的初始化文件`(init_default.json)`**

根据你的需求自行编辑配置文件和初始化文件，编辑完成后去掉文件名中的`_default`（也可以在后续步骤中输入文件名）
> 一个比较重要的配置项是`basic_go_type`，你应该把你用到的所有go语言基础变量类型都添加进去，还有它们对应的ts类型与0值

初始化文件则是用于生成符合约定的go接口文件，默认的初始化文件生成的go接口文件是这样的：  
```cmd
goc_ts -g （也可以通过`-genFrom`参数指定初始化文件）
```

```go
package api

const URI_ListUser = "/user/list"

type ListUserReq struct{}

type ListUserRes struct{}

const URI_CreateUser = "/user/create"

type CreateUserReq struct{}

type CreateUserRes struct{}
```

根据你的需求，编辑go接口文件，例如我们对`demo`文件作出如下补充：

```go
type ResBase struct {
IsSuccess bool   `json:"is_success"`
Err       string `json:"err"`
}

type Pagination struct {
PageNum  int `json:"page_num"`
PageSize int `json:"page_size"`
}

type ListUserReq struct {
Operator string     `json:"operator"`
Page     Pagination `json:"page"`
}

type ListUserRes struct {
Res     ResBase  `json:"res"`
Summary int64    `json:"summary"`
Users   []string `json:"users"`
}
```

最后生成ts代码
```cmd
goc_ts
```
