# kotlin 学习笔记

参考资料：（上次学习时间：2026.8）

- [kotlin tour (beginner + intermediate)](https://kotlinlang.org/docs/kotlin-tour-hello-world.html)
- [kotlin语言指南](https://kotlinlang.org/docs/types-overview.html)

## 定义和打印变量

- val：只读变量，赋初始值后不可更改（建议默认使用该类型，除非确实有修改需要 (?)）
- var：可变变量

字符串模板：访问变量或表达式，将其转为字符串（常用于打印）

```kotlin
val v: Int = 10  // 变量声明
println("value: $v") // 使用时若变量未初始化，则会触发编译错误
println("value: ${v + 1}")
```

## 类型

### 基本类型

- 整型：Byte Short Int Long
    - 长度分别为8 16 32 64 bits
    - 整型字面量支持：十进制子面量 (`100`)、十六进制 (`0xFF`)、二进制 (`0b11`)
- 无符号整型：UByte UShort UInt ULong
- 浮点型：Float Double
- 布尔类型：Boolean
- 字符：Char
- 字符串：String
    - 多行字符串：使用三个引号包围，形如`"""xxx"""`，参考go反引号 (` `` `)，保留内部格式（换行和缩进）、内部不需要转义

整型除法总是舍弃小数部分，如果需要保留，至少将一个操作数设置为浮点型

位运算是以函数的形式存在，例如左移 (`shl()`)、按位与 (`and()`)

### 集合 Collection

- List：普通集合
    - `List是有序的`：文档里的有序并不是内置了排序功能，它只是表达List可以按照顺序索引，相应的Set则不能使用索引访问
- Set：无重复元素集合
- Map：映射，键值对集合

```kotlin
// list
val readOnlyList: List<Int> = listOf(1, 2, 3) // kotlin可以推断集合类型，也可以为变量显式声明类型
val mutableList: MutableList<Int> = mutableListOf(1, 2, 3)

val listCasting: List<Int> = listOf(1, 2, 3) // 可以将可变列表赋值给一个不可变的变量
println(listCasting.count()) // list扩展函数：`count`、`add/remove`、`first/last`

println(10 in listCasting) // in 操作符

// set
val readOnlySet: Set<Int> = setOf(1, 2, 3)
val mutableSet: MutableSet<Int> = mutableSet(1, 2, 3)
println(mutableSet.count()) // set扩展函数：`count`、`add/remove`

// map
val readOnlyMap: Map<Int, String> = mapOf(1 to "yi", 2 to "er")
val multableMap: MutableMap<Int, String> = mutableMapOf(1 to "yi", 2 to "er")
println(multableMap.count()) // map扩展函数：`count`、`remove`、`containsKey`、`keys/values`

multableMap[9] = "jiu" // 可变map也可以正常赋值
println(multableMap[100]) // `null`

println(10 in multableMap.keys) // 在key中查找可以不带`.keys`
println("yi" in multableMap.values)
```

### 类型别名

可能是kotlin也知道自己整出来的类型名很长，它允许定义类型别名。

在程序看来，类型原名和别名是一样的，所以它们之间可以随意交换使用

```kotlin
typealias UserIndex = Map<Long, User>
typealias FileTable<K> = MutableMap<K, MutableList<File>>
```

## 控制流 Control Flow

### 条件表达式 if when

```kotlin
val d: Int
if (ok) {
    d = 1
} else {
    d = 2
}

// if用作表达式，参考rust
// 语句块如果只有一行，`{}`可以省略
val d: Int = if (ok) 1 else 2

// when，参考go的switch、rust的match
// when即可以用作语句（无返回），也可以用作表达式（有返回）
// when用作表达式时建议显式提供变量类型；用作表达式时必须考虑到所有情况或者提供默认分支
// when和subject一起使用，还会检查是否已覆盖所有可能情况，参考rust match enum/Result/Option
val res: String = when (d) { // 报错：when有多种类型返回值，建议把res类型修改为Any
    1 -> "yi"
    2 -> println("er")
    else -> println("unkonwn")
}

// when不带参数，每一个分支都要是一个bool表达式
val res = when {
    d > 0 -> d
    else -> 0
}
```

### 范围

```kotlin
1..4 // 1,2,3,4
1..<4 // 1,2,3
4 downTo 1 // 4,3,2,1
1..5 step 2 // 1,3,5
```

### 循环 for while

```kotlin
// for 遍历一个集合/范围
for (counter in 1..5) {
    print(counter) // 12345
}

// while
while (counter < 10) {
    counter++
}

// do while
do {
    counter++
} while (counter < 10)
```

## 函数 function

```kotlin
fun f(x: Int = 10, Y: Int): Int { // 默认参数
    return x + y
}

f(x = 10, y = 20) // 命名参数：你可以为输入参数命名，这样你可以以任意顺序提供参数，或者跳过默认参数

fun f() {} // 没有返回值的函数，返回值类型为`Unit`

fun sum(x: Int, y: Int): Int = x + y // 如果函数体只有一行，可以这样写

// lambda函数：略

// 函数类型，不同的输入、输出参数使函数拥有不同类型，例如定义一个闭包：
val f: (Int, Int) -> Int = fun(x: Int, y: Int): Int {
    return x + y
}

// 扩展函数：不改变当前代码，为一个类型扩展功能
fun String.bold(): String = "<b>$this</b>" // 扩展函数中使用`this`表示receiver

fun main() {
    // "hello" is the receiver
    println("hello".bold()) // <b>hello</b>  
}
```

## 类 class

```kotlin
class Customer(val id: Int, var email: String) { // 类头，也是构造函数
    var price: Int = 0

    fun cost(money: Int) {
        price -= money
    }
}

// 数据类，用于持有数据，可以方便的打印、判等、复制
data class User(val name: String, val id: Int)

/* 继承 */

// kotlin的类只支持单继承，所有类最终都继承自Any。接口可以多继承
fun main() {
    val k = Kotlin("Mario's")
    println(k.program("practice"))
    println(k.programDemo("practice"))
}

// 抽象类：不能创建实例，默认可以被继承
// 普通属性/方法：需要在抽象类中使用'open'修饰符，属性和方法才能在子类中重写
// 抽象属性/方法：继承必须实现抽象类的*全部*抽象属性和方法
abstract class ProgramLanguage(val note: String) {
    open val id: Int = 0
    abstract val name: String

    abstract fun hello()

    open fun program(str: String): String {
        return "$note $str"
    }
}

interface ProgramDemo {
    // 接口不应该带有属性，因为kotlin最终会编译成java字节码运行在java虚拟机上，而JVM规定接口不能持有数据
    // 接口在一些情况下会声明属性，但是这个属性最终还是由实现接口的类持有
    fun programDemo(str: String): String
}

// 普通类添加'open'修饰符变成开放类，开放类可被继承
open class Kotlin(note: String) : ProgramLanguage(note), ProgramDemo { // 先写继承、再写接口
    override val id: Int = 1
    override val name: String = "Kotlin"

    override fun hello() = println("Hello, Kotlin")

    override fun program(str: String): String {
        return "$note $name $str"
    }

    override fun programDemo(str: String): String {
        return "$note $name $str (face to object language)"
    }
}

/* 对象 Object：只有一个实例的类 singleton */

object Kotlin {
    val language = "Kotlin"

    fun program(str: String): String {
        return "$language $str"
    }
}

println(Kotlin.program("practice"))

// 与类相似，对象也有数据对象 data object，只是数据对象没有复制方法，因为对象只能有一个实例

/* 枚举类 */

// 常规枚举
enum class Language {
    Java, Kotlin
}

// 带数值的枚举
enum class Language(val n: String) {
    Java("Java"),
    Kotlin("Kotlin"); // 枚举类可以带有方法，如果带有方法，枚举常量与方法中见要用分号(`;`)分隔

    fun isKotlin(): Boolean {
        return this.n == Kotlin.toString()
    }
}
```

## 空安全 null safety

kotlin会在编译期间检查null相关

```kotlin
var nullable: String? = null // 可空类型

nullable?.toString() // 安全调用，如果nullable为null，则表达式的值为null，不会抛出错误

nullable?.toString() ?: "" // elvis操作符，`?:`表示如果左操作数是null，则表达式的值为右操作数
```

```kotlin
val str: Any = "str"
println(str is String) // true，判断变量是否是指定类型，'!is'表示对判断结果取反
// is可以用来检查子类型

val s = str as String // 强制类型转换，失败会panic；使用'as?'转换失败时，表达式的值为null
```

## 包和导入 package and import

第三方软件包：https://klibs.io

```kotlin
package com.mario.kotlin

import com.mario.kotlin as K
```

## 关键词

硬性关键词，无论什么时候都不能用做标识符：（28）  
as break class continue do else false   
for fun if in interface is null   
object package return super this throw true  
try typealias typeof val var when while
