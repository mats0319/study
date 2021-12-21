# go官方文档阅读笔记

2021.6开始阅读，阅读版本：

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
    1. 二进制：以```0b```或```0B```开头
    1. 八进制：以```0```或```0o```或```0O```开头
    1. 十六进制：以```0x```或```0X```开头，以```a~f```或```A~F```表示10~15
    1. 单独的```0```表示十进制0
1. 为了可读性，可以在**整型字面量中间**和**基数前缀后**加下划线("_")，这些下划线对字面量的值没有影响，参考代码：
   ```go 
   var (
       intAlias1 int = 100_000
       intAlias2 int = 0x_0010
   )
   ```

### 浮点型字面量(floating-point literals)

1. 允许使用十进制和十六进制（以```0x```或```0X```开头）
1. 允许使用下划线增加可读性，规则同上
1. 浮点型字面量包含：整数部分、小数点、小数部分和指数部分；
    1. 十进制：
        1. 省略规则：整数部分和小数部分可以省略一个，小数点和指数部分可以省略一个
        1. 指数部分使用```e```表示按照10的倍数缩放
        1. 示例：```1.2e+3 // 1200```
    1. 十六进制（以```0x```或```0X```开头）：
        1. 省略规则：整数部分和小数部分可以省略一个，小数点可以省略，指数部分不可省略
        1. 指数部分使用```p```表示按照2的倍数缩放
        1. 示例：```0x1.1p+1 == 0x11p-3 // 2.125```

### 变量(variables)

变量是一段保存值的存储空间，变量值的可选范围由类型决定，变量的默认值为其类型的零值  
interface类型的变量也有一个明确的动态类型（运行时分配，可变，但总是和值匹配；除了预定义标识```nil```，它没有类型）

### 类型(types)

类型决定了值的范围（字符串变量不能使用数字类型赋值）、操作（例如能否用“==”、内置len()/cap()函数等）和方法（为类型定义的方法）  
新类型可通过**别名声明**和**类型定义**获得：

1. 仅通过别名声明获得的新类型可以继承底层类型的方法集，且受到package的限制：
    1. ```type sNew S```，```sNew```类型不继承```S```类型的方法集
    1. ```type sAlias = S```，```sAlias```类型继承```S```类型的方法集
    1. ```type sAlias2 = another_package.S```，```sAlias2```仅继承另一个包中，```S```类型的可导出方法集

类型的方法集：

1. 每一个类型都有一个关联的方法集，方法集可能为空
1. 类型的方法集决定类型实现了哪些接口
1. 方法集中的每一个方法都要有一个唯一、非空的名字
1. 接口类型的方法集就是它的接口：即使接口类型变量的实例拥有更多的方法，在用它实现接口的时候，方法集只保留接口的方法集
1. 值类型的方法集仅包含：接收者为**值类型**的方法
1. 指针类型的方法集包含：接收者为**值类型**和**指针类型**的方法
1. 拥有内嵌字段(T)的结构体类型(S)的方法集包含内嵌字段的方法集，具体规则如下：
    1. 内嵌字段为类型名(T)，S和\*S的方法集包含T类型接收者为值类型的方法
    1. 内嵌字段为类型指针(\*T)，S和\*S的方法集包含T类型接收者为值类型和指针类型的方法

字符串类型：

1. 值是一串```byte```，utf-8编码，只读
1. 其长度为字节数（非负），可使用内置函数```len()```获取，字符串常量的长度在编译时确定
1. 字符串字节可以通过索引访问，但不能取地址(&)

数组类型：

1. 是一串固定长度(numbered)的单一类型变量序列
1. 其长度总是非负的，可以用内置函数```len()```获取，长度是类型的一部分（不可变），要求为非负的int
1. 数组元素可以通过索引访问

切片类型(slice types)：

1. 是数组的一个连续段的描述符（底层是数组，slice是一个引用）
1. 零值为```nil```
1. 长度在运行期间可变，可使用内置函数```len()```获取
1. 切片元素可以通过索引访问
1. 一但初始化，就会关联一个底层数组（数组持有值，slice只是一个引用），可能存在多个切片共用同一个数组

结构体类型：

1. 结构体的嵌入式字段(embedded field)可以是**类型名**或**类型指针**，其中类型指针指向的类型，不能是接口类型或指针类型
1. 字段声明后面，可以带有一个字符串字面量的标签(tag)，可以通过反射接口在运行时获得，参考代码：
   ```go 
   struct {
       Name string `json:"name"`
   }
   ```

映射类型(map types)：

1. 对于```map[key]value```，```key```的类型可以是任何定义了```==```和```!=```的类型，或者是由定义了上述比较运算符的类型实现的接口类型
   （即，假设接口的动态类型定义有上述比较运算符，则该接口类型也可用作```map```的```key```），参考代码；
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
1. 长度为```key```的数量，运行期间可变，可使用内置函数```len()```获取

通道类型(channel types)：

1. 通过发送和接收指定类型的元素，为并发函数提供交流机制
1. channel包括三种类型（方向）：只读channel、只写channel和双向channel，仅可将双向channel强转成单向channel
1. 有缓冲的channel，可以通过```len()```函数获取当前缓冲区数据量，不受channel是否关闭的影响

## 类型和值的属性(properties of types and values)

### 类型标识符(type identity) / 类型相同(identical)

> 本节主要讲类型相同，这里的相同指可以直接把一个类型的值直接赋给另一个类型的变量，不需要强制类型转换。  
> 概念：类型相同可用于赋值、类型断言等场景

两个类型要么相同，要么不同，具体规则如下：

1. 已定义的类型(defined type)总是与其他任何类型都不同；相应的，类型别名总是和原类型相同
1. 否则，两个类型相同条件为：它们的底层类型字面量在结构上等价，即它们拥有相同的字面结构(literal structure)，且对应组件类型相同，具体的：
    1. 数组：元素类型相同且长度相同
    1. 切片：元素类型相同
    1. 结构体：字段序列相同，且对应字段名称、标签(tag)均相同。不同包的非导出的字段名称总是不同的
    1. 指针：基础类型相同
    1. 函数：输入、输出参数数量相同，对应参数类型相同，是否含有可变参数(形如```...int```，variadic)一致。对参数名称没有要求
    1. 接口：方法集相同。不同包的非导出方法名称总是不同的。对方法顺序没有要求
    1. 映射：key、value的类型相同
    1. 通道(channel)：元素类型相同，方向相同

## 表达式(expressions)

### 切片表达式(slice expressions)

切片表达式是从字符串、数组、数组指针或切片中，构造的一个子字符串或切片；可以指定两个或三个参数。 表达式形如```a[ low : high : max ]```：

1. ```0 <= low <= high <= max <= cap(a)```
1. max参数可选，无论指定几个参数，最后一个参数不可省略；省略的参数采用默认值：low的默认值为```0```，high、max的默认值为```len(a)```
1. 通过切片表达式构造的切片，其长度为```high-low```，容量为```max-low```
1. 如果a是一个数组指针，```a[:]```是```(*a)[:]```的简写(shorthand)，参考代码：
   ```go 
   var ap *[3]int = &[3]int{}
   se := ap[:] // typeof(se): []int
   ```

### 类型断言(type assertions)

形如```v, ok := x.(T)```，要求x是interface类型。可以不带ok，但断言失败会panic

T的类型：

1. T是接口类型，类型断言表达式的结果表示：x的动态类型(dynamic type)是否实现了接口T。
1. T不是接口类型，类型断言表达式的结果表示：x的动态类型是否与T相同

如果类型断言成立，则v的值为x保存的值，类型为T
(T为接口类型时，可能打印的类型结果是实例的类型，但v本身仍是接口类型，可以通过尝试调用实例的非接口方法进行验证)， 否则，v的类型为T，值为T类型的零值

### 可变参数传递(passing arguments to ... parameters)

带有可变参数的函数，形如```f(str string, v ...int) {}```

1. 可变参数只能放在参数列表的最后
1. 可变参数在函数内表现为切片类型，在前面的例子中，f函数内，v的类型为```[]int```
1. 如果传递给可变参数的是一个切片，可变参数将与传入的切片相同，且共用底层数组，参考代码：
   ```go 
   func f(str string, v ...int) {}
   
   // caller
   var s = []int{1, 2, 3}
   f("input slice：", s...) // f方法内，对v变量的修改，会体现在s变量上
   ```

## 语句(statements)

### 标签语句(labeled statements)

标签可以用于goto、break和continue语句

### 表达式语句(expression statements)

以下内置函数不能出现在表达式语句中：  
append cap complex imag len make new real  
unsafe.Alignof unsafe.Offsetof unsafe.Sizeof

### 赋值(assignments)

左操作数要求可寻址、是```map```的```key```或者是空白标识符（```_```仅可用于```=```赋值，不可用于```+=```等赋值操作）

元组赋值(tuple assignment)将多值操作的各个元素分配给变量列表，有两种形式：

1. 右操作数是**一个**多值表达式，形如```x, y = f()```，例如函数调用、映射操作、类型断言等
1. 左右操作数数量相等，形如```x, y = 1, 2```

元组赋值可以使用空白标识符忽略右操作数：```x, _ = f()```

元组赋值分成两个步骤：

1. 计算左操作数中的索引表达式、间接指针
2. 计算右操作数中的表达式
3. 从左到右依次赋值

即，元组赋值中，不适合既为一个变量赋值，又使用这个变量，参考代码：

```go 
 var (
     a int = 1
     b map[int]int = make(map[int]int)
 )

 b[1] = 10

 a, b[a] = 2, 20 // 赋值前，首先将b[a]转化为b[1]，然后分别为a、b[1]赋值
 fmt.Println(len(b), b) // 1, {1: 20}
```

赋值运算中，要求右操作数类型能够赋值给对应的左操作数，具体的：

1. 空白标识符可以接收任何类型的值
1. 将未定义类型的常量(untyped constant)赋值给接口类型变量或空白运算符时，常量首先被隐式转换成默认类型

### for语句(for statements)

三种形式：

1. 单一条件(single condition)，类似其他语言的while循环，形如```for a < b { }```
1. for子句(for clause)，类似其他语言的for循环，形如```for i := 0; i < 10; i++ { }```
1. range子句(range clause)，与range关键字联合使用，形如```for i, v := range arr { }```

其中与range关键字联合使用的形式，可用于范围表达式(range expression)的类型以及可以获得的值如下：

| type           | e.g.                    | 1st value | 2nd value |
|----------------|-------------------------|-----------|-----------|
| array or slice | [3]int *[3]int [ ]int   | index     | value     |
| string         | string                  | index     | rune      |
| map            | map[int]string          | key       | value     |
| channel        | (chan int) (<-chan int) | value     |           |

对于channel，range会持续从中读数据，直到channel关闭；如果channel为nil，range会一直阻塞

在for语句中声明的变量（在for子句和range子句形式中，使用短变量声明），作用域是for语句体，每次迭代会被重复使用

### go语句(go statements)

go语句在相同的地址空间内，使用独立的并发线程执行函数调用，形如```go func() { }()```  
表达式要求是函数或方法，不能用括号括起来，不支持部分内置函数（参考表达式语句一节）

### select语句(select statements)

select与switch结构上相似，但不支持fallthrough关键字  
select语句要求每一个case都是通信操作(communication operations)

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

	fmt.Printf("side effects   : %v\n", sideEffects)
	fmt.Printf("selected times : %v\n", selected)
	fmt.Printf("remain quantity: %v\n", remain)
	
        // side effects   : [1000 1000 336]
        // selected times : [328 336 336]
        // remain quantity: [328 336 0]
}

func addOne(slice []int, index int) int {
	slice[index]++
	return slice[index]
}
```

### return语句(return statement)

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

### break语句(break statement)

break语句终止当前函数中，与break语句最接近的for/switch/select语句  
break语句可以带标签，标签只能出现在for/switch/select语句的上一行，且该语句的语句列表(statement list)包含break语句

### continue语句(continue statement)

continue语句终止执行最近for循环的后续语句  
continue语句可以带标签，标签只能出现在for语句的上一行，且该语句的语句列表包含continue语句， 带标签的continue语句表示跳过标签表示的for循环的后续语句，执行下一次循环。参考代码：

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

### goto语句(goto statement)

从接触编程开始，几乎所有的意见都是：不要用goto  
可一方面，go语言保留了goto关键字，作为仅有的25个关键字之一，在go源码中也有所体现；  
另一方面，go github有issue提议：取消使用goto，重写对应源码。

个人倾向于不使用goto，所以本节仅简单描述goto语句规则，为了能看懂使用goto的代码。

goto语句需要一个标签，然后将控制转移到对应标签的位置，要求标签与goto语句在同一函数内。

不建议使用goto跳过变量声明

不能使用goto跳转到另一个局部代码块中（例如不能在一个for循环外面，跳转到里面）

### defer语句(defer statement)

部分内置函数不能defer，参考表达式语句一节

defer函数的参数会被立刻求值  
defer函数在return语句为返回值赋值之后、函数返回到它的调用者之前执行，参考代码：

```go 
func f() (n int) {
	defer func() {
		n++
	}()

	return 1
}

// f() = 2
// 以上代码等价于：
func f() (n int) {
	n = 1
	n++
	return
}
```

如果defer的函数为nil，则执行时panic  
defer函数可能改变当前函数的返回值  
defer函数自己的返回值会被丢弃

defer语句会**压栈（注册）一个**函数：

1. 如果defer表达式为一串链式调用，则之前的调用函数会立刻执行，仅最后一个函数会被压栈
1. 后注册的defer函数先执行

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

若显式调用os.Exit()，则注册的defer函数不会执行，参考代码：

```go 
func main() {
    defer println(1) // no print

    os.Exit(0)
}
```

## 内置函数(built-in functions)

内置函数没有go标准类型，所以只能出现在调用表达式中，不能用作函数类型变量。  
理解：内置函数虽然叫函数，但没有标准类型的它们更像是一种语法糖。

### close

关闭channel，表示不会再写该channel  
以下行为会导致panic：

1. 写一个已关闭的channel
1. 试图再次关闭一个已关闭的channel
1. 关闭一个未初始化的channel（值为nil）

关闭一个channel，且channel的缓存已清空之后，读channel的行为会返回channel类型的零值，多值返回的第二个参数为false

### len & cap

接收各类型的变量，返回int，可以接收nil

| call   | argument type            | result                         |
|--------|--------------------------|--------------------------------|
| len(v) | array ( [3]int *[3]int ) | array length (3 in example)    |
|        | slice ( [ ]int )         | slice length                   |
|        | channel ( chan int )     | channel缓冲区的元素个数，无缓冲channel结果为0 |
|        | string                   | string length in bytes         |
|        | map ( map[string]int )   | map length, number of keys     |
||||
| cap(v) | array ( [3]int *[3]int ) | array length (3 in example)    |
|        | slice ( [ ]int )         | slice capacity                 |
|        | channel ( chan int )     | channel buffer capacity        |

切片的容量是它底层数组的元素数

无论何时，都有```0 <= len(v) <= cap(v)```

参数为数组和数组指针（不包含非常量函数调用）时，len() / cap()函数结果为常量

## packages

go程序通过把包连接在一起而构建  
一个包可以有一个或多个源文件，它们共同声明属于包的常量、变量、类型和函数  
每一个源文件包含（构建参数等，例如```// +build ```不在讨论范围内）：

1. 包声明(package clause)
    1. 包声明不能是空白标识符("_")
    1. 一个包的源文件都要在一个路径中、一个路径中只能有一个包的源文件
1. 导入声明(import declarations)，可能为空
    1. 导入声明指出，当前源文件依赖导入的包的功能，且可以访问导入的包的可导出标识符(exported identifiers)
    1. 导入声明通过路径导入包
    1. 默认使用包名访问其可导出的标识符，也可以自定义包的访问名，特例与参考代码：
        1. 使用“.”作为自定义访问名，表示可以像调用当前包内标识符一样，调用目标包的可导出标识符
        1. 使用“_”作为自定义访问名，表示只调用目标包的init函数，不导入任何标识符
        ```go 
        import (
            "github.com/mats9693/document" // 常规调用格式
            . "github.com/mats9693/utils/uuid" // 假设包含A()方法
        ) 
       
        func f() {
            A() // 正常需要使用 uuid.A()
        }
        ```
    1. 导入声明声明了当前包和导入的包之间的依赖关系，不允许出现以下行为：
        1. 一个包直接或间接地导入它自己
        1. 导入一个包但是没有用到
1. 顶级声明(top-level declarations)，包括常量、变量、类型和函数，可能为空

## 程序初始化和执行(program initialization and execution)

### 零值(zero value)

变量分配存储空间时，如果没有明确的初始化，则会初始化为对应类型的零值

### 包初始化(package initialization)

首先初始化没有依赖的变量，然后重复初始化“依赖已经初始化完成的变量”  
如果初始化步骤结束时，仍有变量未初始化，即存在变量循环依赖，编译报错：initialization loop  
多值表达式会同时初始化，参考代码：

```go 
var (
	x    = a
	a, b = f() // a、b同时初始化，且在x初始化之前
)
```

在初始化环节，空白变量(blank variables)与其他变量统一处理

init函数可以定义很多个，但它们只在初始化阶段被调用，而不能在其他任何阶段调用

包的初始化过程：

1. 初始化包级(package-level)变量
1. 按照源代码顺序调用（多个）init函数
1. 如果当前包有多个源文件，则根据它们到达编译器的顺序初始化（不要依赖同一个包内的init函数的执行顺序）
1. 如果当前包有依赖(import)，则优先初始化依赖的包（有依赖关系的包之间，可以依赖init函数的执行顺序）
1. 如果一个包被多个包依赖，被依赖的包也只会初始化一次

通过构建，保证不会出现包的循环依赖

包的初始化——包括变量初始化和init函数调用——是单线程的，即使有的init函数开启新的线程，初始化程序也会等待该init函数返回后，再执行下一步

### 程序执行(program execution)

一个完整的程序从单一的非导入的(unimported)main包开始连接、创建  
main包要求包名为main，并且声明一个没有输入输出参数的main函数  
程序从初始化main包、调用main函数开始执行，当main函数返回时，程序退出，不会等待其他goroutine完成

## errors

“errors are values”

## 注意事项(system considerations)

### unsafe包

内置的unsafe包已经由编译器实现，可以通过导入“unsafe”访问，更方便底层编程，包括违反类型系统的操作(operations that violate the type system)

使用unsafe的包必须手动审查类型安全，且可能无法移植

提供的函数详情参考源码unsafe包，等用到的时候再来看源码注释和文档。
