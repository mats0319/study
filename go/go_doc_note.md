# go官方文档阅读笔记

文档版本：

1. go语言编程规范：2021.2.10
2. effective go：go 1.17版本

## 词汇要素(lexical elements)

### 预定义标识(predeclared identifiers)

> 眼熟它们，定义类型、变量和函数的时候，尽量避开这39个名字。  
> 你可以重新定义以下的类型、变量和函数，并且能够正常编译和运行，只是执行的将是你定义的内容  
> nil是6种类型的零值：pointer / slice / channel / map / func / interface

Types:  
bool byte complex64 complex128 error float32 float64  
int int8 int16 int32 int64 rune string  
uint uint8 uint16 uint32 uint64 uintptr

Constants:  
true false iota

Zero value:  
nil

Functions:  
append cap close complex copy delete imag len  
make new panic print println real recover

### 关键字(keywords)

> 总计25个，会单独开一个文档，记录从每个关键字引出的知识点。

| keywords |             |        |           |        |
|----------|-------------|--------|-----------|--------|
| break    | default     | func   | interface | select |
| case     | defer       | go     | map       | struct |
| chan     | else        | goto   | package   | switch |
| const    | fallthrough | if     | range     | type   |
| continue | for         | import | return    | var    |

### 运算符(operators)

> 总计47个，会单独开一个文档，解释每个运算符的含义和注意事项。

| operators |        |     |      |              |     |     |     |     |
|-----------|--------|-----|------|--------------|-----|-----|-----|-----|
| +         | &      | +=  | &=   | &&           | ==  | !=  | (   | )   |
| -         | &#124; | -=  | =    | &#124;&#124; | <   | <=  | [   | ]   |
| *         | ^      | *=  | ^=   | <-           | \>  | \>= | {   | }   |
| /         | <<     | /=  | <<=  | ++           | =   | :=  | ,   | ;   |
| %         | \>>    | %=  | \>>= | --           | !   | ... | .   | :   |
|           | &^     |     | &^=  |              |     |     |     |     |

### 整型字面量(integer literals)

1. 进位制
    1. 二进制：以`0b`或`0B`开头
    1. 八进制：以`0`或`0o`或`0O`开头
    1. 十六进制：以`0x`或`0X`开头，以`a~f`或`A~F`表示`10~15`
    1. 单独的`0`表示十进制0
1. 为了可读性，可以在**整型字面量中间**和**基数前缀后**加下划线(`_`)，这些下划线对字面量的值没有影响，参考代码：
   ```go 
   var (
       intAlias1 int = 100_000
       intAlias2 int = 0x_0010
   )
   ```

### 浮点型字面量(floating-point literals)

1. 允许使用十进制和十六进制（以`0x`或`0X`开头）
1. 允许使用下划线增加可读性，规则同上
1. 浮点型字面量包含：整数部分、小数点、小数部分和指数部分；
    1. 十进制：
        1. 省略规则：整数部分和小数部分可以省略一个，小数点和指数部分可以省略一个
        1. 指数部分使用`e`表示按照10的倍数缩放
        1. 示例：`1.2e+3 // 1200`
    1. 十六进制（以`0x`或`0X`开头）：
        1. 省略规则：整数部分和小数部分可以省略一个，小数点可以省略，指数部分不可省略
        1. 指数部分使用`p`表示按照2的倍数缩放
        1. 示例：`0x1.1p+1 == 0x11p-3 // 2.125`

### 变量(variables)

变量是一段保存值的存储空间，变量值的可选范围由类型决定，变量默认初始化为对应类型的零值  
interface类型的变量也有一个明确的动态类型（运行时分配，可变，但总是和值匹配；除了预定义标识```nil```，它没有类型）

### 类型(types)

类型决定了变量的取值范围（字符串变量不能使用数字类型赋值）、操作（例如能否用“==”、内置len()/cap()函数等）和方法（为类型定义的方法）  
新类型可通过**别名声明**和**类型定义**获得：

1. 仅通过别名声明获得的新类型可以继承底层类型的方法集，且受到package的限制（不受指针限制，即`= S`与`= *S`一致）：
    1. ```type sNew S```，```sNew```类型不继承```S```类型的方法集
    1. ```type sAlias = S```，```sAlias```类型继承```S```类型的方法集
    1. `type sAlias2 = another_package.S`，`sAlias2`仅继承另一个包中，`S`类型的可导出方法集

类型的方法集：

1. 每一个类型都有一个关联的方法集，方法集可能为空
1. 类型的方法集决定类型实现了哪些接口
1. 方法集中的每一个方法都要有一个唯一、非空的名字
1. 接口类型的方法集就是它的接口：即使接口类型变量的实例拥有更多的方法，在用它实现接口的时候，方法集只保留接口的方法集
1. 值类型的方法集仅包含：接收者为**值类型**的方法
1. 指针类型的方法集包含：接收者为**值类型**和**指针类型**的方法
1. 拥有内嵌字段`T`的结构体类型`S`的方法集包含内嵌字段的方法集，具体规则如下：
    1. 内嵌字段为类型名`T`，`S`和`*S`的方法集包含`T`类型接收者为值类型的方法
    1. 内嵌字段为类型指针`*T`，`S`和`*S`的方法集包含`T`类型接收者为值类型和指针类型的方法

#### 字符串类型

1. 值是一串`byte`，utf-8编码，只读
1. 其长度为字节数（非负），可使用内置函数`len()`获取，字符串常量的长度在编译时确定
1. 字符串字节可以通过索引访问，但不能取地址(&)

#### 数组类型

1. 其长度总是非负的，可以用内置函数`len()`获取，长度是类型的一部分（不可变）
2. 数组元素可以通过索引访问
3. 数组是值类型
    1. `a = b`会重新分配内存并copy所有的元素
    2. 数组作为函数的输入参数，函数内看到的会是原数组的拷贝，即，函数内对数组的修改带不出去

#### 切片类型(slice types)

1. 一但初始化，就会关联一个底层数组（数组持有值，slice只是一个引用），可能存在多个切片共用同一个数组
    1. `cap(slice)`，获取slice的容量，该值为slice关联的底层数组的元素数
2. slice持有一个底层数组的引用
    1. 把一个slice赋值给另一个slice，它们共用同一个底层数组
    2. slice作为函数的输入参数，函数内对slice底层数组元素的修改，函数外可见
3. 零值为`nil`
4. 长度在运行期间可变，可使用内置函数`len()`获取
5. slice元素可以通过索引访问

#### 结构体类型

1. 结构体的嵌入式字段(embedded field)可以是**类型名**或**类型指针**，其中类型指针指向的类型，不能是接口类型或指针类型
1. 字段声明后面，可以带有一个字符串字面量的标签(tag)，可以通过反射接口在运行时获得

#### 映射类型(map types)

1. `map`的`key`可以是任意支持判等运算符的类型(`==`)，例如整型、字符串、指针、接口等
   ```go 
   type Inter interface {
   }
   
   type imp int
   type imp2 []int
   
   func main() {
       var a imp = 1
       var i1 Inter = a
   
       var b imp2
       var i2 Inter = b
   
       var m map[Inter]int = make(map[Inter]int, 1)
       m[i1] = 100
       v, ok := m[i1]
       println(v, ok)
   
       m[i2] = 200 // panic: runtime error: hash of unhashable type imp2
       v, ok = m[i2]
       println(v, ok)
   }
   ```
2. `map`持有底层数据类型的引用，所以将`map`作为函数的输入参数时，函数内对`map`的修改，函数外可见
3. 尝试根据一个不存在的`key`获取值时，会返回对应类型的零值
   1. 可以引入第二个返回值来判断是`key`不存在，还是`key`对应的值就是该类型的零值，e.g. `v, ok := m[key]`
4. 删除一个`key`，可以使用内置函数`delete`，该函数是安全的，即使删除一个不存在的`key`也不会出错
5. 长度为`key`的数量，运行期间可变，可使用内置函数`len()`获取

#### 通道类型(channel types)

1. 通过发送和接收指定类型的元素，为并发函数提供交流机制
1. channel包括三种类型（方向）：只读channel、只写channel和双向channel，仅可将双向channel强转成单向channel
1. 有缓冲的channel，可以通过```len()```函数获取当前缓冲区数据量，不受channel是否关闭的影响

### 类型相同(identical)

> **相同**类型的变量之间可以相互赋值，不需要强制类型转换。  
> 概念：类型相同可用于赋值、类型断言等场景

两个类型要么相同，要么不同，具体规则如下：

1. 已定义的类型(defined type)总是与其他任何类型都不同；相应的，类型别名总是和原类型相同
2. 否则，两个类型相同条件为：它们的底层类型字面量在结构上等价，即它们拥有相同的字面结构(literal structure)，且对应组件类型相同，具体的：
    1. 数组：元素类型相同且长度相同
    2. 切片：元素类型相同
    3. 结构体：字段序列相同，且对应字段名称、类型、标签(tag)均相同。不同包的非导出的字段总被认为是不同的
        1. `var _ struct{ i, j int } = struct{ i, j int }{}`，匿名结构体可以以此方式重新赋值
    4. 指针：基础类型相同
    5. 函数：输入、输出参数数量相同，对应参数类型相同，是否含有可变参数(形如`...int`，variadic)一致。对参数名称没有要求
    6. 接口：方法集相同。不同包的非导出方法名称总是不同的。对方法顺序没有要求
    7. 映射：key、value的类型相同
    8. 通道(channel)：元素类型相同，方向相同

## 表达式(expressions)

### 切片表达式(slice expressions)

切片表达式是从字符串、数组、数组指针或切片中，构造的一个子字符串或切片；可以指定两个或三个参数。 表达式形如`a[ low : high : max ]`：

1. `0 <= low <= high <= max <= cap(a)`
1. `max`参数可选，无论指定几个参数，最后一个参数不可省略；省略的参数采用默认值：`low`的默认值为`0`，`high`、`max`的默认值为`len(a)`
1. 通过切片表达式构造的切片，其长度为`high-low`，容量为`max-low`
1. 如果a是一个数组指针，`a[:]`是`(*a)[:]`的简写(shorthand)

### 类型断言(type assertions)

形如`v, ok := x.(T)`，要求`x`是`interface`类型。可以不带`ok`，但断言失败会`panic`

1. `T`是接口类型时，类型断言表达式的结果表示：`x`的动态类型(dynamic type)是否实现了接口T。
1. `T`不是接口类型，类型断言表达式的结果表示：`x`的动态类型是否与T相同

变量`v`的类型为`T`，根据断言是否成立，变量的值为`x`保存的值或`T`类型的零值

1. 常规打印变量类型的方式，打印接口类型变量，会输出其动态值的类型；此处可以通过**尝试调用变量`v`的非接口方法**验证其类型

### 可变参数传递(passing arguments to ... parameters)

带有可变参数的函数，形如`f(str string, v ...int) {}`

1. 可变参数只能放在参数列表的最后
1. 可变参数在函数内表现为切片类型，在前面的例子中，f函数内，v的类型为`[]int`
1. 如果传递给可变参数的是一个切片，可变参数将与传入的切片相同，且共用底层数组，参考代码：
   ```go 
   func f(str string, v ...int) {}
   
   // caller
   var s = []int{1, 2, 3}
   f("input slice：", s...) // f方法内，对v变量的修改，会体现在s变量上
   ```

## 内置函数(built-in function)

### close

关闭channel，表示不会再写该channel  
以下行为会导致panic：

1. 写一个已关闭的channel
1. 试图再次关闭一个已关闭的channel
1. 关闭一个未初始化的channel（值为nil）

### len & cap

接收各类型的变量，返回int，可以接收nil

| call   | argument type | e.g.           | result                         |
|--------|---------------|----------------|--------------------------------|
| len(v) | array         | [3]int *[3]int | array length (3 in example)    |
|        | slice         | [ ]int         | slice length                   |
|        | channel       | chan int       | channel缓冲区的元素个数，无缓冲channel结果为0 |
|        | string        | string         | string length in bytes         |
|        | map           | map[string]int | map length, number of keys     |
|||||
| cap(v) | array         | [3]int *[3]int | array length (3 in example)    |
|        | slice         | [ ]int         | slice capacity                 |
|        | channel       | chan int       | channel buffer capacity        |

切片的容量是它底层数组的元素数

无论何时，都有```0 <= len(v) <= cap(v)```

参数为数组和数组指针（不包含非常量函数调用）时，len() / cap()函数结果为常量

## packages

1. go程序通过把包连接在一起而构建
2. 一个包可以有一个或多个源文件，它们共同声明属于包的常量、变量、类型和函数

每一个源文件包含（构建参数等，例如`// +build `不在讨论范围内）：

1. 包声明(package clause)
    1. 包声明不能是空白标识符(`_`)
    1. 一个包的源文件都要在一个路径中、一个路径中只能有一个包的源文件
1. 导入声明(import declarations)，可能为空
    1. 导入声明指出，当前源文件依赖导入的包的功能，且可以访问导入的包的可导出标识符(exported identifiers)
    2. 导入声明通过路径导入包
    3. 默认使用包名访问其可导出的标识符，也可以自定义包的访问名，特例与参考代码：
        1. 使用`.`作为自定义访问名，表示可以像调用当前包内标识符一样，调用目标包的可导出标识符
        2. 使用`_`作为自定义访问名，表示只调用目标包的init函数，不导入任何标识符
        ```go 
        import (
            "github.com/mats9693/document" // 常规调用格式
            . "github.com/mats9693/utils/uuid" // 假设包含A()方法
        ) 
       
        func f() {
            A() // 正常需要使用 uuid.A()
        }
        ```
    4. 导入声明声明了当前包和导入的包之间的依赖关系，不允许出现以下行为：
        1. 一个包直接或间接地导入它自己
        1. 导入一个包但是没有用到
1. 顶级声明(top-level declarations)，包括常量、变量、类型和函数，可能为空

## 程序初始化和执行(program initialization and execution)

### 零值(zero value)

变量分配存储空间时，如果没有明确的初始化，则会初始化为对应类型的零值

### 包初始化(package initialization)

1. 首先初始化没有依赖的变量，然后重复初始化**依赖已经初始化完成的变量**
2. 如果初始化步骤结束时，仍有变量未初始化，即存在变量循环依赖，编译报错：`initialization loop`
3. 多值表达式会同时初始化，参考代码：
    ```go 
    var (
        x    = a
        a, b = f() // a、b同时初始化，且在x初始化之前
    )
    ```
4. 在初始化环节，空白变量(blank variables)与其他变量统一处理
5. `init`函数可以定义很多个，但它们只在初始化阶段被调用，而不能在其他任何阶段调用

包的初始化过程：

1. 初始化包级(package-level)变量
2. 按照源代码顺序调用（多个）`init`函数
3. 如果当前包有多个源文件，则根据它们到达编译器的顺序初始化
    1. 不要依赖同一个包内的`init`函数的执行顺序
4. 如果当前包有依赖(import)，则优先初始化依赖的包
    1. 有依赖关系的包之间，可以依赖`init`函数的执行顺序
5. 如果一个包被多个包依赖，被依赖的包也只会初始化一次

通过构建，保证不会出现包的循环依赖

包的初始化——包括变量初始化和`init`函数调用——是单线程的，即使有的`init`函数开启新的线程，初始化程序也会等待该`init`函数返回后，再执行下一步

### 程序执行(program execution)

一个完整的程序从单一的非导入的(unimported)main包开始连接、创建

main包要求包名为main，并且声明一个没有输入输出参数的main函数

程序从初始化main包、调用main函数开始执行，当main函数返回时，程序退出，不会等待其他goroutine完成

1. 在GMP模型中，`goroutine0`返回即表示程序结束

## 注意事项(system considerations)

### unsafe包

内置的unsafe包已经由编译器实现，可以通过导入“unsafe”访问，更方便底层编程，包括违反类型系统的操作(operations that violate the type system)

使用unsafe的包必须手动审查类型安全，且可能无法移植

提供的函数详情参考源码unsafe包，等用到的时候再来看源码注释和文档。

## 短变量声明

```go 
a := f()
a, b := g()
```

短变量声明语法的目的，是复用变量，举例来说，在一长串的`if-else`中，只需要使用一个`err`，并且它在每个地方都表示我们想要它表示的含义

对于短变量声明中存在相同名称变量的情况（如上示例代码），满足以下全部条件时，短变量声明不会创建新的变量

1. 两次短变量声明在同一个作用域内
2. 两次短变量声明都是**有效**的：
    1. 相应结果可以赋值给对应类型的变量（即上例中`g`函数的第一个结果，可以赋值给a变量）
    2. 第二次短变量声明中，至少有一个新的变量（如上例中的`b`变量，空白标识符不算新的变量）

## 空白标识符

1. 与短变量声明结合
    1. 用于多返回值函数，忽略某个不关心的值 e.g. `_, err := os.Stat(path)`
    2. 类型断言 e.g. `_, ok := val.(Interface)`，判断变量`val`是否实现了`Interface`接口
2. （开发过程中）忽略报错：**未使用的**导入或变量
    1. 引入一个包但未使用、定义一个变量但未使用，都是一种go语言不允许的浪费，且可能预示着更大的错误
    2. 开发过程中可能会遇到这种情况，你当然可以删除这些未使用的内容，也可以使用空白标识符
3. 用于引入一个包，并且只为了该包的副作用（side effects，此处指包的init函数）
    1. e.g. `import _ "net/http/pprof`
4. 接口检查(interface checks)
    1. e.g. `var _ Interface = (*Struct)(nil)`，检查`Struct`类型是否实现了`Interface`接口；  
       习惯上，仅在没有静态转换的代码中使用这种声明，即，假设包含以下代码，则不需要该声明：
    ```go 
    func NewStruct() Interface {
        return &Struct{} // 这里已经包含了接口检查
    }
    ```

## 数据(data)

### 构造函数(constructor)与初始化

> 与C不同，Go语言可以返回局部变量的地址

内置函数new：

1. 创建一个指定类型的变量，赋值为对应类型的**零值**，返回指向该变量的**指针**

内置函数make：

1. 只能用于创建`slice`、`map`、`channel`，返回**已初始化**的**指定类型**
    1. `make([]int, 1, 3)`，返回一个长度为`1`、容量为`3`的`[]int`

自定义构造函数：

```go 
// from package os
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := new(File)
	f.fd = fd
	f.name = name
	f.dirinfo = nil
	f.nepipe = 0
	return f
}
```

结构体的初始化，可以使用复合字面量(composite literal)：

1. `&File{fd, name, nil, 0}`，匿名复合字面量的顺序按照结构体字段顺序，且需要与结构体字段一一对应
2. `&File{fd: fd, name: name}`，具名复合字面量可以只初始化部分字段且不必考虑字段顺序，其余字段将被初始化为对应类型的零值
3. `&File{}`，作为上一条的特例，当复合字面量没有字段时，所有字段被初始化为零值，该写法与`new(File)`等价

### 二维slice(two-dimensional slices)

1. 因为slice长度可变，所以二维slice天然可以保存没有对齐的矩阵，示例代码：
    ```go 
    type LinesOfText [][]byte     // A slice of byte slices.
    text := LinesOfText{
        []byte("Now is the time"),
        []byte("for all good gophers"),
        []byte("to bring some fun to the party."),
    }
    ```

### print

> 仅列举标准打印函数部分格式符号，详情参考fmt包文档

`%v`：

1. 默认打印格式
2. 打印`map`时，输出结果按照`key`的字典序排列
3. 打印结构体时，`%+v`会打印字段名
4. 任意类型，`%#v`按照完整的go语法打印结果

`%q`：

1. 打印`string`、`[]byte`类型时，带引号
2. `%#q`，尽量使用反引号代替引号
3. 打印**数字类型**或`rune`类型时，输出一个单引号`rune`常量

`%x`：可用于`string`、`[]byte`、**数字类型**，生成一个长的十六进制字符串（字母小写，使用`X`表示字母大写）

`%T`：打印变量类型

`stringer`接口：

1. 未指定打印类型时，若待打印变量实现了`stringer`接口，则调用类型的`String()`方法作为结果并打印
    1. `stringer`接口：`type stringer interface { String() string }`
    2. 如果想要让类型的**值变量**和**指针变量**打印结果一样，需要使用**值类型接收者**实现`stringer`接口

若接口的实现中调用了打印函数，可能导致**打印函数 -> `String()`函数**的循环调用；满足以下全部条件时，产生循环：

1. 打印函数的格式是字符串，如：`fmt.Sprintf("%x", m)`，包括但不限于`%s`、`%x`、`%v`
2. 给的值，其类型实现了`stringer`接口
    ```go 
    type MyString int

    func (m MyString) String() string {
        return fmt.Sprintf("MyString=%x", m) // Error: will recur forever, stack overflow
    }
    ```

相应的，避免循环调用，可以从这几方面入手：

1. 明确指定非字符串类型的打印格式，如：`fmt.Sprintf("%d", m)`
2. 给的值，其类型未实现`stringer`接口，可以使用强制类型转换，如`fmt.Sprintf("MyString=%x", string(m))`

## 初始化

### 常量

1. 只有数字、字符、字符串和布尔类型可以定义为常量
2. 常量在编译期间创建，即使是函数内部定义的常量

常量枚举器：iota

1. iota表示从0开始，逐渐增加的自然数
2. iota可以用作表达式的一部分，且表达式可以隐性重复
    ```go 
    const (
        _           = iota // ignore first value by assigning to blank identifier
        KB ByteSize = 1 << (10 * iota)
        MB
        GB
        TB
    )
    ```
3. 一个文件中，定义了几个`const ()`，iota在每一个`const ()`中，都是从0开始

## 嵌入(embedding)

1. 接口可以包含另一个接口
2. 结构体可以包含另一个结构体，但情况更复杂：
    1. 内嵌的是结构体值类型还是指针类型
    2. 内嵌的结构体字段是否具名

一个结构体（以下称为内嵌结构体）匿名嵌入另一个结构体（以下称为外层结构体），这种方式称为**组合**

1. 组合将**内嵌结构体的方法集**中，外层结构体没有的方法，加入外层结构体的方法集
2. 调用到外层结构体的方法时，如果该方法来自内嵌结构体的方法集，则方法的接收者为内嵌结构体，即，此时方法内无法访问外层结构体中的其他字段
3. 外层结构体可以通过类型名（不带包名）访问内嵌结构体

内嵌结构体可能带来命名冲突问题，例如**内嵌结构体类型名相同**和**不同的内嵌结构体存在同名方法**，相关规则：

1. 不允许在同一嵌套层级嵌入同名结构体
2. 若同一层级内，不同的内嵌结构体存在同名方法，则在调用时，不允许通过外层结构体直接调用方法，而是要指明内嵌结构体的名字

## 并发(concurrency)

> 并发编程是一个很大的话题，本节只是列举一些go语言在这方面的亮点

### 通过通信共享内存(share by communicating)

1. 不通过共享内存通信，而是通过通信共享内存
    1. 通过`channel`传递数据，同一时刻只会有一个goroutine访问，设计上避免了数据竞争

如何理解这个模型：

1. 假设有一个典型的单线程程序在运行，它本身不需要同步化处理
2. 假设再运行一个相同的实例，它也不需要同步化处理
3. 现在让它们通信，如果通信是同步的，也不需要同步化处理
    1. go channel是同步的，详见channel阻塞机制

### Goroutines

起名叫goroutine，是因为现有的术语表达的含义都不准确

goroutine有一个简单的模型：它是一个在同一地址空间中，与其他goroutine并发执行的函数

goroutine被多路复用到多个系统级线程上(OS threads)

在函数或者方法前面使用`go`关键字，该函数或方法将在一个新的goroutine中执行，（函数）调用结束后，goroutine退出

## Errors

### recover

```go 
func safelyDo(work *Work) {
    defer func() {
        if err := recover(); err != nil {
            log.Println("work failed:", err)
        }
    }()
    do(work)
}
```
