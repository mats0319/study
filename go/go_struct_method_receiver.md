# go结构体方法接收者

## 实现接口

> go可以为自定义类型实现方法，这意味着自定义类型也可以实现接口，但这不在本文的讨论范围内，本文仅讨论结构体类型。  
> go调用结构体方法有两种方式，一种是`[struct ins].method()`，一种是`[struct type].method([struct ins])`，本文统一采用第一种方法。

go语言认为：如果一个结构体以自身为接收者、实现了一个接口内的全部方法，则该结构体实现了相应接口

```go 
type I interface {
    funcPointerReceiver()
    funcValueReceiver()
}

var _ I = (*s)(nil)

type s struct {
}

func (s *s) funcPointerReceiver() {
    fmt.Println("in func pr")
}

func (s s) funcValueReceiver() {
    fmt.Println("in func vr")
}
```

上方代码中，结构体`s`实现了接口`I`，但没有完全实现 ( ^_^ ) 。

你可能已经注意到了，我们在结构体`s`的方法中，分别使用了**指针类型接收者**和**值类型接收者**，为了方便叙述，我们做出以下约定：

1. 指针类型接收者方法：形如`func (s *s) f() { }`，方法接收者为结构体的指针类型，以下简称为**指针方法**
2. 值类型接收者方法：形如`func (s s) f() { }`，方法接收者为结构体的值类型，以下简称为**值方法**
3. 结构体的指针类型实例：形如`var sp *s`，类型为结构体指针的变量，以下简称为**指针变量**
4. 结构体的值类型实例：形如`var sv s`，类型为结构体的变量，以下简称为**值变量**
5. 接口类型实例（值为结构体的指针类型）：形如`var iIns I = &s{}`，变量类型为接口类型，值为结构体的指针类型变量，以下简称为**接口类型的指针变量**
6. 接口类型实例（值为结构体的值类型）：形如`var iIns I = s{}`，变量类型为接口类型，值为结构体的值类型变量，以下简称为**接口类型的值变量**

回到刚才的问题，为什么我们说结构体`s`并没有完全实现接口`I`呢？  
因为此时结构体`s`的**值变量**不被认为实现了接口`I`（`var sv I = s{}` 报错）

一个接口类型实例的值可以是指针变量或值变量、一个方法可以是指针方法或值方法，简单排列组合，一共有4种**接口类型实例调用方法**的情况：

1. 接口类型的**指针变量**尝试调用**指针方法**
2. 接口类型的**指针变量**尝试调用**值方法**
3. 接口类型的**值变量**尝试调用**指针方法**
4. 接口类型的**值变量**尝试调用**值方法**

以上情况中，仅情况**3**不被允许。

原因是指针方法内部可以修改接收者（下节展开），而如果通过值变量调用，传进去的是值变量的副本，方法内的修改无法应用到调用者（方法的接收者），  
故go语言直接在语言层面上不允许这种行为（该描述仅针对**接口类型**的值变量，如果使用值变量本身，是可以调用指针方法的，再下节展开）。

因为值变量无法调用指针方法，又因为go不允许函数重载，导致一个接口函数对应的指针方法与值方法不能共存，所以只要有指针方法，就不认为值变量实现了接口。

## 方法内修改接收者

方法的接收者可以简单理解为函数的输入参数，所以它和其他输入参数有一样的问题：函数内的修改能否带出来？

结论：go语言是值传递，所以如果变量本身是一个地址，那么函数内的修改能带出去。

一个有意思的例子：

```go 
slice := make([]int, 0, 3)
fmt.Println(slice) // []

f := func (s []int) {
  s = append(s, 10)
  fmt.Println(s) // [10]
}

f(slice)
fmt.Println(slice) // []

fmt.Println(slice[:cap(slice)]) // [10 0 0]
```

在函数内通过索引修改`slice`的值（即修改`slice`持有引用的底层数组），函数外是可见的；而函数内对`slice len`的修改则带不出来。  
因为`slice`的实现，其值通过指针指向一个数组，即使经过值传递的拷贝，指针仍旧指向该数组（只是指针本身的地址变了）；  
而`slice`的`len`，因为是一个`int`，在`slice`作为输入参数时被拷贝，函数内修改的`len`是一份临时的拷贝，无法带出来。

## 值变量调用指针方法

有别于接口类型的值变量，（结构体类型的）值变量可以调用指针方法，是因为go编译器会把**可寻址的值类型变量**重写为：`(&s).[method name]`

这样一来，对于值变量调用指针方法，就像**函数有一个`*int`类型的输入参数，我们传一个`int`类型变量的地址进去，函数内的修改能带出来**一样了。

```go 
// 值变量调用指针方法，示例代码
type s struct {
    name string
}

func (s *s) funcPointerReceiver() {
    s.name = "in func pr"
}

var sv s = s{}
sv.funcPointerReceiver() // go编译器将当前行重写为“(&sv).funcPointerReceiver()”
// sv.name = "in func pr"
```

类比：函数有一个`*int`类型的输入参数，函数内的修改可以带出来

```go 
func carryOutModify(ip *int) {
    *ip = 1
}

var intIns int
carryOutModify(&intIns) // intIns = 1
```

不可寻址的情况举例

```go  
func newS() s {
    return s{name: "default name"}
}

func main() {
    newS().funcPointerReceiver() // 编译不通过，提示不能对“newS()”取地址
}
```

## 总结：接收者类型定义

本节将结合一些应用场景，讨论为方法定义哪种类型的接收者更合适。

1. 结构体很大：例如结构体中包含一个1k的字符串，使用指针方法可以避免无意义地拷贝这1k内容
2. 方法内需要修改接收者：虽然值方法也可以通过**返回+重新赋值**实现修改效果，但参考上一条：定义指针方法可以避免可能存在的无意义的拷贝
3. …… 如果只考虑结构体类型的话，似乎什么场景都可以使用指针方法？别急，来看看这个：

```go 
type s struct {
}

func (si *s) String() string { // ATTENTION on method receiver <-
    fmt.Println("in self-define 'String' method")
    return ""
}

fmt.Println("> value: ", s{})
fmt.Println("> ptr: ", &s{})
```

背景知识：go语言标准打印函数，在未指定输出格式时，优先检查待打印类型是否实现了`stringer`接口；  
若类型实现了接口，则调用类型的`String()`方法作为结果并打印。

如果我们使用**指针类型接收者**实现接口，则打印**值类型变量**时，将不会执行`String()`函数，因为此时值类型变量没有实现接口

再看看这个：

```go 
// 两个问题：
// 1，为什么Unmarshal不是预期的结果：
//  解：其实marshal的结果就不对了，因为结构体的值类型变量不被认为实现了指针方法，进而不被认为实现了接口，
//     所以json.marshal(d)调用的不是我们为JsonData结构体实现的marshal方法，而是json包的marshal方法
//    （见json.marshal函数注释，如果类型实现了Marshaler接口，会优先调用类型的方法）
//     又因为结构体的name字段是非导出的，所以json.marshal的结果是：{}，长度为2，内容为空
// 2，两种方法分别是怎么修改代码（加一个‘&’，删一个‘*’）：见注释-solution
type JsonData struct {
    name string
}

func (c *JsonData) MarshalJSON() (bytes []byte, err error) { // solution 1: 删除本行的‘*’
    return json.Marshal(c.name)
}

func (c *JsonData) UnmarshalJSON(bytes []byte) error {
    return json.Unmarshal(bytes, &c.name)
}

func TestJsonData(t *testing.T) {
    d := JsonData{name: "test"}
    bytes, err := json.Marshal(d) // solution 2: 改成 'json.Marshal(&d)'
    if err != nil {
        fmt.Println("1", err)
        return
    }
    //fmt.Println(len(bytes), string(bytes)) // "test", string, length = 6

    var d2 JsonData
    err = json.Unmarshal(bytes, &d2)
    if err != nil {
        fmt.Println("2", err)
    }
    if d != d2 {
        fmt.Println("3", d, d2)
    }

    return
}
```

这个例子告诉我们：如果你想要实现别人的接口，通常要注意更多内容，例如：结构体字段的可见性、对方接口的特殊情况（例如本例中的marshal规则）等。  
值得一提的是，就像本例中使用值变量的错误一样，使用指针方法可以通过编译，但会在运行时报错；这时，除了提供文档描述、建议使用指针变量以外，还可以使用值方法。

总结

1. 结构体功能性代码：
    1. 显式检查结构体是否实现了接口：```var _ I = (*s)(nil)```
2. 指针方法可以将方法内对接收者的修改带出去（无论使用值变量还是指针变量调用）
3. 编译器会重写可寻址的值变量调用指针方法，但这不代表值变量实现了指针方法
4. 通常情况下，我们都可以选择**定义指针方法+使用指针变量**，这也符合go语言的使用习惯；  
   但在实现别人的接口、且总体环境复杂的情况下，可以考虑使用值方法来降低运行时错误的概率。
5. 如果你看到这里，什么都没记住，那么就记住下面这句话吧：全部定义指针方法、全部使用指针变量、在调用方法前，保证变量非空（防止调用值方法时解引用`*`错误）
