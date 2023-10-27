# Goc_ts

根据go语言定义的结构，生成http请求对应的ts代码

命名参考了`protoc_gen_go_grpc`

ts代码：vue3 axios

## 使用

> 参考`goc_ts/sample/manual.md`

注意事项：

1. 结构体字段要有`json`tag，例如：```ID string `json:"id"` ```
2. http请求定义以`const URI_`开头，值为请求的uri，例如：`const URI_GetList = "/list/get"`
3. 根据http请求定义中的名称，自动派生其请求参数、响应参数结构，接上例：`GetListReq`/`GetListRes`结构体将被视为http请求的请求参数和响应参数结构
4. 工具内还是有一些固定的字符串表述，例如配置文件的文件名固定为`config.json`等，这些没有做成启动项（flag），主要是我认为没有必要

## 实现

使用正则匹配go语言结构定义文件，按照一定格式记录下来。  
根据记录的内容，生成ts代码，包括：axios配置初始化(`config.ts`)、工具函数(`utils.ts`)、http请求的数据结构(`xxx.go.ts`)
以及http请求函数(`xxx.http.ts`)。

1. 若根据分析，http请求没有用到任何工具函数，则不会生成工具函数文件
2. 工具通过预置代码模板、根据模板中的锚点替换实际参数的方式，实现生成ts代码
    1. 举个例子，以下是一个http请求函数的ts代码模板，在模板中预留了很多形如`{{ $xxx }}`的锚点。
       工具通过将这些锚点替换成带有实际意义的字符串，实现生成代码、控制缩进等功能。
       值得一提的是，部分较为简单的模板分散在代码各处、没有集中定义，所以如果想要改ts代码总体的结构，可能比较麻烦
        ```go
        var ServiceCode_Request = "\n" +
        "{{ $indentation }}public {{ $serviceNameSmall }}({{ $paramsWithType }}): " +
        "Promise<{{ $serviceName }}Res> {{{ $serviceCode_ReqStruct }}\n" +
        "{{ $indentation }}{{ $indentation }}return axiosWrapper.post(\"{{ $serviceURI }}\"{{ $requestParams }})\n" +
        "{{ $indentation }}}\n"
        ```

## FAQ

1. 为什么想做这样一个工具？
   为了统一定义、管理http请求结构，避免出错。  
   曾考虑过使用protocol buffer来做这个事情，放弃的原因是不想只为了使用它的序列化和反序列化而引入一门新的语言（pb）。
2. 为什么选择根据go代码生成ts代码？
   接上一条，因为不想引入一门新的语言，所以工具要么根据go生成ts、要么根据ts生成go。  
   至于为什么选择了根据go生成ts，主要是个人喜好。
