# go运算符

> 整理自go语言官方文档，部分内容增加代码、文字解释，2021.2.10版本

The following character sequences represent operators (including assignment operators) and punctuation:  
翻译：以下字符序列表示运算符（包括赋值运算符）和标点符号。  
理解：go只有这47个运算符(?)

|operators| | | | | | | |sum:47|
|---|---|---|---|---|---|---|---|---|
|+|&|+=|&=|&&|==|!=|(|)|
|-|&#124;|-=|&#124;=|&#124;&#124;|<|<=|[|]|
|*|^|*=|^=|<-|&gt;|&gt;=|{|}|
|/|<<|/=|<<=|++|=|:=|,|;|
|%|&gt;&gt;|%=|&gt;&gt;=|--|!|...|.|:|
| |&^| |&^=| | | | | |

二元运算符优先级（数字越大，优先级越高）：

|precedence|operator|
|---|---|
|5|* / % << >> & &^|
|4|+ - &#124; ^|
|3|== != < <= > >=|
|2|&&|
|1|&#124;&#124;|

相同优先级，按照从左到右的顺序计算

## 算数运算符(arithmetic operators)

算数运算符仅应用于数值，产生与第一个操作数同类型的结果

|运算符|含义|适用类型|
|---|---|---|
|+|加|integers, floats, complex values, strings|
|-|减|integers, floats, complex values|
|*|乘|integers, floats, complex values|
|/|除|integers, floats, complex values|
|%|取余|integers|
||||
|&|按位与(bitwise AND)|integers|
|&#124;|按位或(bitwise OR)|integers|
|^|按位异或(bitwise XOR)|integers|
|&^|位清除(bit clear AND NOT)|integers|
||||
|<<|左移位|integer << unsigned integer|
|&gt;&gt;|右移位|integer >> unsigned integer|

&^运算符规则：**1 0 为 1，其他（3种情况）为 0**，参考代码：

```go 
var (
   a int = 0b_0011_1100
   b int = 0b_0110_0110
)

fmt.Printf("%b\n", a &^ b) // 1_1000

// a &^ b:
// 1. c = (按位取反)(b)
// 2. a & c
```

整型除法结果向零取整(truncated towards zero)  
对于```q = x / y```、```r = x % y```，满足```x = q*y + r  and |r| < |y|```，即：

|x|y|x/y|x%y|
|---|---|---|---|
|5|3|1|2|
|-5|3|-1|-2|
|5|-3|-1|2|
|-5|-3|1|-2|

最小的int(-128 for int8, -32768 for int16)，乘或除以-1，值不变（因为二进制补码的整型溢出）

|x|x/4|x%4|x>>2|x&3|
|---|---|---|---|---|
|11|2|3|2|3|
|-11|-2|-3|-3|1|

```go 
// 补码：-11 = -11 + 2^8 = 245 = 255 - 10 = 255 - 8 - 2，简易计算规则：按位取反，+1 
-11 = 0b_1111_0101
-11 >> 2 = 0b_1111_1101 // 补码还原，得：-3
-11 & 3:
  0b_1111_1101
& 0b_0000_0011
--------------
             1
```

x << 1 约等于 x * 2，x >> 1 约等于 x / 2；区别在于，移位运算符向下取整(truncated towards negative infinity)

整型溢出：自动忽略高位，不会报错。  
编译器优化代码时，会考虑到溢出的情况，例如编译器不会把```x < x + 1```优化为```true```

## 比较运算符(comparison operators)

|运算符|含义|
|---|---|
|==|equal|
|!=|not equal|
|<|less|
|<=|less or equal|
|\>|greater|
|\>=|greater or equal|

比较运算符的两个操作数必须可以相互赋值  
判等运算符(```==```和```!=```)要求两个操作数是可比较的，其他运算符要求两个操作数是已排序的(ordered)，具体规则如下：

1. 布尔值可比较：两个布尔值在同为ture或false时相等
1. 整型可比较、已排序
1. 浮点型可比较、已排序
1. 字符串可比较、已排序：参考代码：
    ```go 
    var a, b string
    print(a > b) // 判断a、b的第一个不同字符的ascii码
    ```
1. 指针可比较：指向同一个变量或值为nil的指针相等，指向大小为0的不同变量(distinct zero-size variables)的指针，可能相等也可能不相等
1. 通道可比较：使用同样的make函数构造的channel相等；值为nil的channel相等
1. 接口可比较：动态类型和动态变量都相等的接口相等；值为nil的接口相等
1. 非接口类型变量(s)和接口(I)可比较，要求类型s可比较且实现了接口I：当接口I的动态类型是s的类型且动态变量是s时，s和I相等
1. 结构体可比较，要求它们的字段都是可比较的：当两个结构体对应的非空字段相等时，两个结构体相等
1. 数组可比较，要求元素类型可比较且长度相同：对应位置元素相等的数组相等
1. 比较两个拥有相同动态类型，但值不可比较的接口，panic。该规则同样适用于接口数组、包含接口字段的结构体的比较
1. 切片、映射和函数不可比较，除非他们的值为nil

## 接收运算符(receive operator)

对于channel类型，形如```<-ch```的接收操作，表示从channel ch中获取值，要求channel的方向允许读（即ch不能是只写channel，chan<-类型）

## 转换(conversions)

转换分为显式转换(explicit conversion，又称强制转换)和隐式转换(implied conversion)  
强制类型转换表达式形如：```T(x)```，T类型在可能引起歧义的地方，可以用括号括起来，例如以*、<-开头的类型

强制类型转换规则```T(x)```：

1. x可分配(assignable)给T
1. 结构体，忽略标签(tags)后，对应的字段有相同的底层类型或指针指向非已定义的类型(not defined types)且指针基础类型的底层类型相同
1. x的类型和T都是整型或浮点型
1. x的类型是整型、[]byte、[]rune，T是string类型
1. x的类型是string，T是[]byte、[]rune

数字类型非常量(non-constant)之间的强制转换：有符号整型的转换可能会扩展精度（不知道这个说法准不准确，文档没看明白），参考代码：

```go 
// 目前总结出的规律是：如果int类型值为负，且二进制首位为1，则转成更高精度的uint，会在高位补1
v := uint16(0x10F0)
uint32(int8(v)) // 0xFFFFFFF0
```
