# go关键字（draft）

### 全部关键字列举

|keywords|||||
|---|---|---|---|---|
|break|default|func|interface|select|  
|case|defer|go|map|struct|  
|chan|else|goto|package|switch|  
|const|fallthrough|if|range|type|  
|continue|for|import|return|var|  

### 从关键字引出的知识点

#### break

1. 可以用在for循环、switch/select的case中
1. 可以通过带一个参数（标签），退出并跳过标签标识的for/switch/select语句

```go 
ALL:
    for range [3]struct{}{} {
        for range [3]struct{}{} {
            break ALL
        }
    }

// break ALL退出后，不会再执行该循环，而是从这里继续执行
```

#### case

可以用在switch/select中，详情参考对应关键字

#### chan

1. 可以用来声明和定义channel，channel包含读写、只读、只写三种类型，  
   只读和只写多见于函数参数，为了控制该函数仅能读或写该channel
2. 使用不同状态channel的结果：

| |只读|只写|备注|
|---|---|---|---|
|仅声明|死锁|死锁| |
|初始化（无缓冲）|阻塞|阻塞| |
|初始化（有缓冲）|无数据时阻塞|缓冲满时阻塞| |
|已关闭（无缓冲）|返回channel类型的零值[，false]|panic: send on closed channel| |
|已关闭（有缓冲）|有数据时正常返回|panic: send on closed channel|无数据时，同上|

参考代码：

```go 
func nilChan() {
	var ch chan int

	ch <- 1
	<- ch
}

func noBufferedChan() {
	var ch = make(chan int)
	//var ch = make(chan int, 0)

	<- ch
	ch <- 1
}

func bufferedChan() {
	var ch = make(chan int, 1)

	ch <- 1
	ch <- 1
	<- ch
}

func noBufferedClosedChan()  {
	var ch = make(chan int)
	close(ch)

	v, ok := <- ch
	fmt.Println(v, ok)

	ch <- 1
}

func bufferedClosedChan() {
	var ch = make(chan int, 1)
	ch <- 1
	close(ch)

	v, ok := <- ch
	fmt.Println(v, ok)

	ch <- 1
}
```

#### const

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

#### continue

仅适用于for循环内

#### default

可使用于switch/select，与case同级，详情参考对应关键字

#### defer

1. defer需要注册才能执行，后注册的defer先执行，defer会在函数退出前执行，参考代码：

```go 
func f() {
    defer print(1)
	defer print(2)

	return
	defer print(3) // 没有注册，所以不会执行
}

// output: 21
```

1. 每个defer关键字只会把一个函数“压栈”，并且会立刻计算函数的输入参数，参考代码：

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

1. defer可以结合内置函数recover，处理panic，详情参考recover函数声明与注释

#### else

可用于```if...else...```句式

#### fallthrough

可用于switch的case中，表示无条件执行下一个case/default的语句块  
不可用于switch的最后一个case；不可用于类型switch(type switch)

#### for

循环

#### func

声明和定义函数

#### go

开线程执行函数

#### goto

无条件跳转到指定位置（需要标签）

#### if

条件判断

#### import

导入package

#### interface

接口

#### map

映射

#### package

包声明

#### range

可用于```for...range...```句式

#### return

函数退出

#### select

1. 每个case都要求是io操作
1. 随机选择一个可执行的case
1. 所有case阻塞则执行default，无default则select语句阻塞

#### struct

定义结构体

#### switch

条件判断

#### type

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

#### var

定义变量
