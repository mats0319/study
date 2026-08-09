# go小知识

## 逃逸分析

虽然官方文档中表示我 **不需要知道**一个值是分配在栈上还是堆上，但知道它可以让你更容易理解逃逸分析检查结果，并且这不只在go语言中有用

如果go不能确定一个值在函数退出时已经使用完毕，该值就会逃逸到堆上（保守规则）。触发逃逸分析的常见原因：

1. 返回局部变量的指针
2. 闭包和goroutine
3. 接口和动态值
4. slice、map和struct：如果一个值保存在生命周期超过当前函数的结构中，则该值会逃逸到堆上

检查逃逸：`go build -gcflags="-m"`

逃逸分析不能作为性能好坏的证据，基准测试和性能测试才可以。  
逃逸分析只是性能优化的一个方向，它本身不是一个完整的策略，而且不是每个逃逸都需要花精力去处理

## go的标识符可见性：protected

go通常使用首字母大小写控制标识符可见性，但还有一条控制规则：internal文件夹

- 这条规则并非go语言规范的一部分，而是go工具链强制执行的
- 标识符：包括常量、变量、类型、函数、方法
- 规则：假设一个模块有一个internal文件夹，那么：
    - 该模块以外无法访问internal文件夹
    - 该模块内，internal文件夹的父文件夹以外，无法访问internal文件夹
  ```text
  module_1
  - go.mod
  - main.go // &#10006;
  
  module_2
  - go.mod
  - main.go // &#10006;
  - parent
      - main.go // &#10004;
      - external
          - external.go // &#10004;
      - internal
          - child
              - main.go // &#10004;
          - internal.go // define: PublicVar
  
  &#10004;：允许使用`internal.PublicVar`
  &#10006;：不允许使用`internal.PublicVar`
  ```
- 错误信息：在编译时，提示import对应行出错，错误描述为 **不允许使用internal包**
  （这里用词是`package`，但实际上与包名无关，而是指文件夹名称）

## 编译时常量与命令行启动参数

### 编译时常量

``` go
//go:build !release

package tag

const IsDebug = true
```

编译时提供，编译出的二进制文件享受编译器优化

- 编译优化：`if tag.IsDebug {}`，在前面这段代码中，如果bool值为false，则编译后二进制文件中不包含整个if代码块
- 条件编译按文件生效

### 命令行启动参数

``` go 
package flags

import "flag"

var IsTestMode bool

func init() {
	flag.BoolVar(&IsTestMode, "t", false, "test mode flag")
	flag.Parse()
}
```

程序启动时提供，使用相同的二进制文件，类似从配置文件读配置

## 时间类型 (`time.Duration`)

``` 
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
```

go语言时间类型的单位是纳秒（1 sec = 1 * 10^9 nano sec），允许配置时间的地方，例如数据库超时时间，一般的库函数不会做转换，  
即，`var timeout Duration = 1`表示超时时间为1 纳秒，如需设置超时时间为1s，应写成：`var timeout Duration = 1 * time.Second`

## 内置函数append的一个约束

```code 
var nums1 []interface{}
nums2 := []int{1, 3, 4}
nums3 := append(nums1, nums2...)
fmt.Println(len(nums3))
```

以上代码会在编译时报错：不能把[]int类型的变量，当做append函数的参数，赋值给[]interface{}类型  
题目本身更像是一个脑筋急转弯，所以这里不谈题目，只谈知识点：

1. go 拥有完整定义的slice，用作可变参数传入函数时，函数内部看到的可变参数，是函数外定义的slice本身，而不是将slice的一个一个元素取出来组成的新的slice
    1. 完整定义的slice：参考上述代码第二行，将slice定义为一个变量；可变参数位置这样写：`[]int{1,2,3}...`也算 **有完整定义**
2. append函数，如果两个输入参数类型不一致，会在编译期间报错，错误提示见前文

如果append的可变参数写成：`1,2,3`，函数内看到的可变参数，是这几个元素组成的新的slice， 可以被推断为`[]interface{}`类型，不会报错

但像题目中的写法，函数内看到的可变参数，与函数外一致，是有类型的 ([]int)，所以会报错

## 建议将append结果赋值回去

`slice = append(slice, item)`

如上代码，建议使用append的slice，接收append的结果，两个原因：

- 如果slice底层数组容量不足、append触发扩容，且没有将append结果赋值回去时，slice相当于没有append
- 如果slice底层数组容量充足，但是长度不足 (len < cap)，修改会影响到底层数组（多见于基于数组或slice创建slice）

``` go
array := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

slice := array[:3]
_ = append(slice, 100)

fmt.Println(array) // [1 2 3 100 5 6 7 8 9 10]
```

## 避免在类型切换的单个case中编写多个类型 type switch

``` go
type s struct {
    i int
}

var i interface{} = &s{}

switch v := i.(type) {
case int:
    // 此时v是int类型
case *s, s:
    // 此时v是interface{}类型，因为case里有多个类型，包括default分支也是这样；
    // 这一行代码会在编译期间报错：不能用interface类型变量往下点
    // v.i 
}
```

结论：尽量不要在type switch的一个case里写多个类型 ^_^
