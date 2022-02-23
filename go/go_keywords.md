# go关键字（draft）

全部关键字列举，标识`x`表示该关键字没有什么可展开的知识点

| keywords |               |        |           | `x` means no expand |
|----------|---------------|--------|-----------|---------------------|
| break    | default       | func   | interface | select              |  
| case ×   | defer         | go     | map       | struct              |  
| chan     | else  ×       | goto   | package × | switch              |  
| const    | fallthrough × | if     | range     | type                |  
| continue | for           | import | return    | var                 |  

## break

1. 可以用在for循环、switch/select的case中
2. 可以通过带一个参数（标签），退出并跳过标签标识的for/switch/select语句，一个**有效的标签**应满足以下条件：
    1. 标签只能出现在for/switch/select语句的上一行（中间可以有注释和空行）
    2. 标签紧挨着的语句的语句列表(statement list)包含break语句

```go 
ALL:
    for range [3]struct{}{} {
        for range [3]struct{}{} {
            break ALL
        }
    }

// break ALL退出后，不会再执行该循环，而是从这里继续执行
```

## chan

1. 可以用来声明和定义`channel`，`channel`包含读写、只读、只写三种类型，  
   1. 只读和只写多见于函数输入参数位置，为了控制该函数权限
2. 可以使用内置函数`close`关闭`channel`
3. 读写不同状态`channel`的结果：

|            | read             | write                           | note                                                                  |
|------------|------------------|---------------------------------|-----------------------------------------------------------------------|
| 仅声明(值为nil) | 阻塞               | 阻塞                              | 整个进程没有活动`goroutine`时，将触发死锁，错误：`all goroutines are asleep - deadlock!` |
| 初始化-无缓冲    | 阻塞               | 阻塞                              |                                                                       |
| 初始化-有缓冲    | 正常返回             | 缓冲满时阻塞                          | 无数据时，见上一行                                                             |
| 已关闭-无缓冲    | 返回`channel`类型的零值 | `panic: send on closed channel` | `_, ok := <- ch`，`ok`为`false`                                         |
| 已关闭-有缓冲    | 正常返回             | `panic: send on closed channel` | 无数据时，见上一行                                                             |

## const

1. 可使用iota，iota表示的自然数序列会在每个const语句块中初始化，在多行const中，iota的值仅与常量的位置有关 参考代码：

```go 
const (
    c = 'x'   // 'x', iota = 0 is covered
    a = iota  // 1
    b         // 2
    s = "xxx" // "xxx", iota = 3 is covered
    d         // 4
)

const (
    i = iota // 0, iota re-init
    j        // 1
)
```

## continue

1. `continue`语句终止执行最近for循环的后续语句  
2. `continue`语句可以带标签，标签只能出现在for语句的上一行，且该语句的语句列表包含continue语句， 带标签的continue语句表示跳过标签表示的for循环的后续语句，执行下一次循环。参考代码：

```go 
func f() {
ContinueLabel:
	for range [3]struct{}{} {
		fmt.Println("a")
		for i := 0; i < 3; i++ {
			fmt.Println(i)
			continue ContinueLabel
			fmt.Println("b")
		}
	}

	fmt.Println("finish")
	
	// output: a 0 a 0 a 0 finish
}

func f() {
	for range [3]struct{}{} {
		fmt.Println("a")
ContinueLabel:
		for i := 0; i < 3; i++ {
			fmt.Println(i)
			continue ContinueLabel // 本例与不带标签的continue相同
			fmt.Println("b")
		}
	}

	fmt.Println("finish")
	
	// output: a 0 1 2 a 0 1 2 a 0 1 2 finish
}
```

## default

可使用于switch/select，与case同级，详情参考对应关键字

## defer

1. `defer`函数的参数会被立刻求值
2. `defer`函数在`return`语句为返回值赋值之后、函数退出之前执行
   1. `defer`函数可能改变当前函数的返回值
3. 如果`defer`的函数为`nil`，则执行时panic
4. `defer`函数自己的返回值会被丢弃
5. `defer`语句会压栈（注册）一个函数：
   1. 如果`defer`表达式为一串链式调用，则之前的调用函数会立刻执行，仅最后一个函数会被压栈
   2. 后注册的`defer`函数先执行
6. 若显式调用`os.Exit()`，则注册的`defer`函数不会执行

参考代码：

```go 
func f() {
	defer print(1)
	defer print(2)

	return
	defer print(3) // 没有注册，所以不会执行
}

// output: 21
```

```go 
type chainCall struct {
	n int
}

func (c *chainCall) set(value int) *chainCall {
	c.n = value
	print("set value: ", value)

	return c
}

func f() {
    cc := &chainCall{}

	a := 1
	b := 2
	defer cc.set(1).set(2).set(a+b) // 这里压栈的函数是：cc.set(3)

	a += 10
	cc.set(a)
	
	// output: 1 2 11 3
}
```

## for

三种形式：

1. 单一条件(single condition)，类似其他语言的while循环，形如`for a < b { }`
2. for子句(for clause)，类似其他语言的for循环，形如`for i := 0; i < 10; i++ { }`
3. range子句(range clause)，与range关键字联合使用，形如`for i, v := range arr { }`

在for语句中声明的变量（在for子句和range子句形式中，使用短变量声明），作用域是for语句体，每次迭代会被重复使用

## func

1. 函数可以返回多个值
2. 函数输出参数可以命名，然后当做常规变量使用，就像输入参数一样
    1. 具名的输出参数，会在函数开始时，被初始化为对应类型的零值
    2. 若函数有多个输出参数，要么全部命名、要么全部不命名，不允许只为部分输出参数命名
    3. 拥有具名输出参数的函数，`return`可以不带参数

## go

go语句在相同的地址空间内，使用独立的并发线程执行函数调用，形如`go func() { }()`

## goto

从接触编程开始，几乎所有的意见都是：不要用goto  
可一方面，go语言保留了goto关键字，作为仅有的25个关键字之一，在go源码中也有所体现；  
另一方面，go github有issue提议：取消使用goto，重写对应源码。

个人倾向于不使用goto，所以本节仅简单描述goto语句规则，为了能看懂使用goto的代码。

goto语句需要一个标签，然后将控制转移到对应标签的位置，要求标签与goto语句在同一函数内。

不建议使用goto跳过变量声明

不能使用goto跳转到另一个局部代码块中（例如不能在一个for循环外面，跳转到里面）

## if

条件判断

## import

导入package

## interface

接口

## map

映射

## range

可用于```for...range...```句式，遍历`array`、`slice`、`string`、`map`、`channel`类型

| type           | e.g.                    | 1st value | 2nd value |
|----------------|-------------------------|-----------|-----------|
| array or slice | [3]int *[3]int [ ]int   | index     | value     |
| string         | string                  | index     | rune      |
| map            | map[int]string          | key       | value     |
| channel        | (chan int) (<-chan int) | value     |           |

## return

实现约束：如果在局部作用域内，定义了与返回值同名的变量，则该作用域内，不能使用隐式返回。参考代码：

```go 
func f() (err error) {
	{
		err := errors.New("new error")
		if err != nil {
			return // 编译报错：err is shadowed during return
		}
	}

	return
}
```

## select

1. select与switch结构上相似，但不支持fallthrough关键字  
2. select语句要求每一个`case`都是通信操作(communication operations)

select语句执行过程：

1. 计算所有case的表达式，所有副作用(side effects)正常执行；接收语句左侧的赋值语句（包括短变量声明）不会执行。
1. 选择分支：
   1. 如果有可继续的通信，则通过统一的伪随机(uniform pseudo-random)选择一个
   1. 如果有default case，执行default case
   1. 如果不存在以上情况，则select语句阻塞，直到有至少一个可继续的通信（空白的select语句会一直阻塞）
1. 执行准备：
   1. 除非选择的是default case，否则就要执行对应的通信操作
   1. 如果选中case的表达式为赋值语句，则赋值语句正常执行
1. 执行选中case的语句块

以下代码包含“副作用正常执行”、“随机选择一个case执行”以及“全部计算case表达式、只执行选中case的赋值”的演示，参考代码：

```go 
func main() {
	var (
		times = 1000
		c1    = make(chan int, times)
		c2    = make(chan int, times)
		c3    = make(chan int, times)
	)

	sideEffects := make([]int, 3)
	selected := make([]int, 3)
	{
		count := 1
		c3 <- count
		for i := 0; i < times; i++ {
			select {
			case c1 <- addOne(sideEffects, 0):
				selected[0]++
			case c2 <- addOne(sideEffects, 1):
				selected[1]++
			case sideEffects[2] = <-c3:
				selected[2]++
				count++
				c3 <- count
			}
		}

		<-c3
	}

	remain := make([]int, 3)
	{
		chs := []chan int{c1, c2, c3}
		for i := 0; i < len(remain); i++ {
			remain[i] = len(chs[i])
		}
	}

	fmt.Printf("side effects  : %v\n", sideEffects)
	fmt.Printf("selected times: %v\n", selected)
	fmt.Printf("remain amount : %v\n", remain)
	
        // side effects  : [1000 1000 336] // mainly two first values
        // selected times: [328 336 336]
        // remain amount : [328 336 0] // mainly two first values, compare with last line
}

func addOne(slice []int, index int) int {
	slice[index]++
	return slice[index]
}
```

## struct

定义结构体

## switch

1. 类似C语言的switch：`switch [expression] {}`
2. 把一连串`if-else`写成switch：`switch {case [expression]: // do sth}`
    1. 第一个表达式值为`true`的`case`会被执行
3. 类型转换(type switch)
    ```go 
    // from official doc
    var t interface{}
    t = functionOfSomeType()
    switch t := t.(type) {
        default:
            fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
        case bool:
            fmt.Printf("boolean %t\n", t)             // t has type bool
        case int:
            fmt.Printf("integer %d\n", t)             // t has type int
        case *bool:
            fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
        case *int:
            fmt.Printf("pointer to integer %d\n", *t) // t has type *int
    }
    ```

## type

定义类型  
新类型不继承底层类型的方法集，类型别名可以，但直接继承的方法集受到package的限制（可导出的方法可以正常调用非导出的方法），参考代码：

```go 
type S struct {
}

func (s *S) ExportedFunc() {
}

func (s *S) nonExportedFunc() {
}

type sNew S
type sAlias = S // 类型别名声明(alias declarations)
type sAlias2 = another_package.S // 假设another_package包，定义有同样的结构体S与方法集

func main() {
    var sn = &sNew{}
    var sl = &sAlias{}
    var sl2 = &sAlias2{}
    
    sn.ExportedFunc()       // wrong
    sn.nonExportedFunc()    // wrong
    sl.ExportedFunc()       // right
    sl.nonExportedFunc()    // right
    sl2.ExportedFunc()      // right，即使该可导出方法中调用了不可导出的方法，也能正确执行
    sl2.nonExportedFunc()   // wrong
}
```

## var

定义变量
