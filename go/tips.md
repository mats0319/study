# go小知识

## Commands

- go tool compile –S [main.go] >> [main.S] 生成main.go的go汇编，并保存在main.S文件中
- go test -coverprofile=[coverage.txt] 执行当前目录下的测试文件，保存结果为coverage.txt
- go tool cover -html=[coverage.txt] 生成网页形式的代码覆盖率报告
- go build -gcflags -m [file name] 查看编译过程中，优化了哪些代码

## 时间类型(`time.Duration`)

```go 
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
```

go语言时间类型的单位是纳秒（1 sec = 1 * 10^9 nano sec），允许配置时间的地方，例如数据库超时时间，一般的库函数不会做转换，  
即，`var timeout Duration = 1`表示超时时间为1 纳秒，如需设置超时时间为1s，应写成：`var timeout Duration = 1 * time.Second`

## 内置函数append的一个约束

```go 
var nums1 []interface{}
nums2 := []int{1, 3, 4}
nums3 := append(nums1, nums2...)
fmt.Println(len(nums3))
```

以上代码会在编译时报错：不能把[]int类型的变量，当做append函数的参数，赋值给[]interface{}类型  
题目本身更像是一个脑筋急转弯，所以这里不谈题目，只谈知识点：

1. go 拥有完整定义的slice，用作可变参数传入函数时，函数内部看到的可变参数，不是将slice的一个一个元素取出来组成的新的slice，而是函数外定义的slice本身
    1. 完整定义的slice：参考上述代码第二行，将slice定义为一个变量；可变参数位置这样写：`[]int{1,2,3}...`也算**有完整定义**
2. append函数，如果两个输入参数类型不一致，会在编译期间报错，错误提示见前文

如果append的可变参数写成：`1,2,3`，函数内看到的可变参数，是这几个元素组成的新的slice、可以被推断为[]interface{}类型，不会报错  
但像题目中的写法，函数内看到的可变参数，与函数外一致，是有类型的([]int)，所以会报错
