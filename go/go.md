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
