# go小知识

## Commands

- go tool compile –S [main.go] >> [main.S] 生成main.go的go汇编，并保存在main.S文件中
- go test -coverprofile=coverage.txt 执行当前目录下的测试文件，保存结果为coverage.txt
- go tool cover -html=coverage.txt 生成网页形式的代码覆盖率报告
- go build -gcflags -m [file name] 查看编译过程中，优化了哪些代码

## GMP调度模型

todo

## 时间类型(```time.Duration```)

```go 
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
```

go语言时间类型的单位是纳秒（1 sec = 1 * 10^9 nano sec），  
允许配置时间的地方，例如数据库超时时间，一般的库函数不会做转换，  
即，```var timeout Duration = 1```表示超时时间为1 纳秒，  
应写成：```var timeout Duration = 1 * time.Second```

## go命令影响范围

[reference](https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns)

```text
Directory and file names that begin with "." or "_" are ignored by the go tool, as are directories named "testdata".
```

go工具通常会遍历当前目录及其全部子目录，包括```node_modules```。  
所以假设这样一种情况：```/go```路径下包含```ui```项目，  
我们下载了```ui```项目的依赖，即产生了一个```node_modules```文件夹。  
此时，在```/go```路径执行```go mod tidy```命令，是会扫描整个```node_modules```文件夹的，  
想要跳过这种扫描，在Goland编辑器里，可以选择对应目录，设置```make directory as - excluded```；  
或者将对应目录改名：以```.```和```_```开头的文件和文件夹，以及名为```testdata```的文件夹会被go工具忽略
