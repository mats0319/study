# kotlin

## 变量首先要可以变

val: 常量
var: 变量

kotlin很重视变量是否可修改，包括集合里默认创建的都是只读的实例、想要创建可编辑的实例需要更复杂的初始化函数。

## 基本类型

整型: Byte Short Int Long (var year: Int = 2026)

- 长度分别为8/16/32/64 bits
- 支持`0x0F`/`0b01`，不支持8进制数字类型字面量
  无符号整型: UByte UShort UInt ULong (`var score: UInt = 100u`)
  浮点数: Float Double (`var value: Float = 10.5f; var price: Double = 20.99`)
  布尔类型: Boolean (`var isAdmin: Boolean = true`)
  字符类型: Char (`var c: Char = '.'`)
  字符串: String (`var str: String = "Hello !"`)

变量在使用前必须初始化，尝试打印一个没有初始化的变量会触发编译期错误

## 集合

- Lists：列表
    - 只读列表：`var readOnlyList: List<Int> = listOf(1,2,3)`
    - 可变列表：`var mutableList: MutableList<Int> = mutableListOf(1,2,3)`
    - 可用方法：`count（计算长度）`、`in（检查元素是否存在，var v = "str" in strList`、`add/remove`、`first/last`
    - 可以将一个可变的列表赋值给只读列表，集合和映射也是
- Sets：集合（无序，唯一）
    - 只读：`setOf()`/可变：`mutableSetOf()`
    - 可用方法：`count`、`in`、`add/remove`
- Maps：映射（k-v）
    - 只读：`mapOf()`/可变：`mutableMapOf()`
    - 初始化：`var m: Map<Int, String> = mapOf(10 to "a", 20 to "b")`，可变map也可以正常赋值：`map["key"] = value`
    - 尝试通过不在map里的key获取value，将得到`null`
    - 可用方法：`remove`、`count`、`countainsKey(检查key是否存在)`、`keys/values(获得key或value的集合(MutableSet<T>))`、
      `in(检查k或v在不在map里，var v = 1 in map.keys / var v = "a" in map.values)`

## 控制流

条件表达式：if/when（我们使用的其他语言很少有把when当switch用的，所以虽然kotlin推荐使用when，还是决定尽量不用）

范围：`..`，接受开区间、倒序、设置步长

- `1..4`=`1,2,3,4`, `1..<4`=`1,2,3`
- `4 downTo 1`=`4,3,2,1`
- `1..5 step 2`=`1,3,5`

循环：for/while/do-while，很遗憾，这里没有办法兼顾其他语言的经验了

- for：遍历，只能用来遍历范围、集合（list、set、map），`for (value in collection) {}`
- while：循环，条件表达式为真时，执行代码块

## 函数

`fun [func name](): [return type] {}`

函数最多只能返回一个值

默认参数：`fun f(pageSize: Int = 10) {}`

- 一个函数可以有多个默认参数
- 如果默认参数出现在参数列表中间、调用时又忽略了它，那么你应为后续所有参数命名
  `fun f(v1: Int = 0, v2: Int) {}  ->  f(v2=1)`
- 可变数量参数：`fun f(vararg v: String) {}`->`f(v=arrayOf("a","b","c"))`
    - 如果传的是数组（与函数声明中的变量类型不符），则需要在调用时为变量命名
    - 如果该可变数量参数后面还有其他参数，其他参数需要在调用时命名

函数体如果只有一行表达式，则函数可以简写成`fun sum(x:Int,y:Int):Int=x+y`，但还是为了兼顾其他语言使用习惯，通常不用

lambda表达式，略。（时间紧、简单学习；lambda总是能扩写成可读性更好的标准函数，所以也不想学）

## 类

```kotlin
class C(var Field2: String = "") { // 类头，可省略，需要在创建类实例时赋值
    var Field1: Int = 0
}
```

数据类，专门用来存储数据，自带一些数据处理成员函数

- 形如：`data class User(var id: Int, var name: String)`
- 成员函数：`toString()`、`equals()/==`、`copy()（深度拷贝，不影响原有实例）`

继承：单继承，被继承的类需要标记为open/abstract

```kotlin
class C : ParentClass() {
    override var category: String = "C"
}

abstract class ParentClass {
    abstract var category: String
}
```

实现接口：可以实现多个接口（实现类要带括号，接口不带）

```kotlin
interface PaymentMethod {
    // Functions are inheritable by default
    fun initiatePayment(amount: Double): String
}

class CreditCardPayment(val cardNumber: String, val cardHolderName: String, val expiryDate: String) : PaymentMethod {
    override fun initiatePayment(amount: Double): String {
        // Simulate processing payment with credit card
        return "Payment of $$amount initiated using Credit Card ending in ${cardNumber.takeLast(4)}."
    }
}

fun main() {
    val paymentMethod = CreditCardPayment("1234 5678 9012 3456", "John Doe", "12/25")
    println(paymentMethod.initiatePayment(100.0))
    // Payment of $100.0 initiated using Credit Card ending in 3456.
}
```

## 空安全

可空类型：`var str: String? = null`/`str?.length（返回null，不会出错）`/`str?.length ?: 0（自定义为空时的返回）`

## 扩展函数

不改变已有代码，为其增加功能

```kotlin
fun String.bold(): String = "<b>$this</b>"

fun main() {
    // "hello" is the receiver
    println("hello".bold())
    // <b>hello</b>
}
```

## 对象

```kotlin
object DoAuth {
    fun takeParams(username: String, password: String) {
        println("input Auth parameters = $username:$password")
    }
}

fun main() {
    // The object is created when the takeParams() function is called
    DoAuth.takeParams("coding_ninja", "N1njaC0ding!")
    // input Auth parameters = coding_ninja:N1njaC0ding!
}
```

## 第三方库

https://klibs.io

## 语言指南

`package my.code`package定义应放在文件顶部（第一行）

`println()/readln()`打印到控制台、从控制台读取输入

### 关键词

硬性关键词，无论什么时候都不能用做标识符：（28）
as break class continue do
else false for fun if
in interface is null object
package return super this throw
true try typealias typeof val
var when while

### 可见性

private protected internal public(default)

### 类型检查

```kotlin
var v: Any = "a string"
if (v is String) {} // 应用在基本类型上，if代码块内部v已经隐式转换成String类型了（例如可以.length）

var v: Animal = Dog()
if (v is Dog) {} // 应用在自定义class的子类上，if代码块内部v已经隐式转换成Dog类型，可以访问其属性和方法

var v: Any = "a string"
if (v !is String) {return}
// 如果程序执行到这里，v也会被隐式转换成String类型

if (v is String && v.length > 0) {} // 我的理解是：is直接会在作用域内将变量修改为对应类型，包括if表达式、代码块

if (v is T1 || v is T2) {} // 此时会将v隐式转换成T1、T2的最近公共父类
```

强制类型转换：`v as String`，通常看上去像是is的语法糖（类型转换），还可以**强制**转换，例如把父类实例强转成子类类型

```kotlin
var v:Animal = Dog()
var v2 = v as? Dog
```

### 类型别名

`typealias FilesMap<K> = MutableMap<K, MutableList<File>>`

类型别名只是给开发者看的，编译器会将其与原类型统一处理，所以如果函数需要一个原类型变量，可以传一个别名类型变量；反过来也是

## 编写android app界面

- 入口：`oncreate()`
- 布局：`setContent()`
- UI：标有`@Composable`注解的函数


