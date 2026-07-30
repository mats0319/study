# Rust程序设计语言阅读笔记 The Rust Programming Language

[Rust程序设计语言](https://kaisery.github.io/trpl-zh-cn/title-page.html)

rust语言特性：

- 编译器会检查出很多你可能以为是**多管闲事**的错误，但这是有道理的
- 编译器有很强的类型推断能力
- rust是基于表达式的语言(expression-based)

rust语法与我熟悉的go有所不同：

- 语句末尾应带有分号
- 链式调用换行时，点号在下一行
- 几乎全部使用蛇形命名法
- rust range:`let slice1 = &s[0..5]`，像是把go的冒号放平：`slice2:=array[1:5]`
- rust没有空值
- rust实现一个目标的方式可能有很多种，例如`match`和`if let`、字符串拼接加号(`+`)和`format!`宏
    - rust提供给开发者更多的权限，可以更精准的控制性能，但是学习成本也更高

## cargo

> [The Cargo Book](https://doc.rust-lang.org/cargo/)

包管理工具

`cargo new [project name]` 创建一个指定命名的目录和项目，初始化git并生成`.gitignore`（如果已经在一个git仓库中则不会创建）
`cargo build` 编译项目，默认不执行优化、使用debug模式，可以使用`--release`执行优化
`cargo run` 编译项目并运行
`cargo check` 检查当前项目是否可编译
`cargo doc --open` 本地构建所有依赖的文档，并在浏览器打开

## 常见编程概念

### 变量与可变性 Variables and Mutability

#### 变量

`let mut x;`

变量默认是不可变的，只能赋值一次（声明时赋值或后续赋值）。

变量遮蔽(shadow)：使用`let`复用变量名，本质上是定义一个新的变量

```rust
fn f() {
    let x = "100";
    let x: u32 = x.trim().parse().expect("invalid number");
}
```

- 可以改变不可变变量
- 可以（显式）改变变量类型

#### 常量

`const day: u32 = 24*60*60;`

常量与不可变变量的区别：

- 常量永远不可变（不能使用mut）
- 必须声明类型
- 常量可以定义在全局作用域
- 只能用常量表达式赋值

### 数据类型 Data Types

#### 标量 Scalar

包含整型、浮点型、布尔类型和字符类型

整型和浮点型都支持基本数学运算：加减乘除、取余

##### 整型

| 长度      | 有符号   | 无符号   |
|---------|-------|-------|
| 8 bit   | i8    | u8    |
| 16 bit  | i16   | u16   |
| 32 bit  | i32   | u32   |
| 64 bit  | i64   | u64   |
| 128 bit | i128  | u128  |
| 架构相关    | isize | usize |

允许使用十进制、二进制、十六进制、八进制、字节（仅限`u8`）的整型字面量，允许使用`_`作为视觉分隔符

整型溢出：

- debug模式会检查并panic
- release模式不会检查，自动舍去超出范围的位

除法向0舍入

##### 浮点型 Float-Point Numbers

分为`f32`/`f64`两种，都是有符号的

##### 布尔类型

使用`bool`表示，有`true`/`false`两个值

##### 字符类型

`let c = 'z';`

使用`char`表示，字符类型字面量由单引号(`''`)包裹，类型大小为4个字节

#### 复合类型 Compound

包含元组(tuple)和数组(array)

##### 元组 Tuple

元组是一种将多个不同类型的值组合成一个复合类型的通用方式。

元组长度固定：声明后不能改变

```rust
fn f() {
    let tup: (i32, f64, u8) = (500, 6.4, 1);

    let (x, y, z) = tup; // 使用模式匹配(pattern matching)解构(destructure)
    let v = tup.0; // 使用点号(`.`)加索引访问
}
```

单元(unit)：不带任何值的元组，值和类型都写作`()`，常用于表示空值或空返回类型，如果一个表达式没有返回，它会隐式返回一个单元

##### 数组 array

```rust 
fn f() {
    let a: [i32; 5] = [1, 2, 3, 4, 5];
    let a = [3; 5]; // 等价于：[3,3,3,3,3]

    let v1 = a[0]; // 使用索引访问数组元素
}
```

数组：

- 每个元素类型相同
- 长度不可变
- 数据分配在栈(stack)上
- 更灵活的选择：vector

### 函数 Functions

- 语句 Statements：执行一些操作，没有返回，例如`let y = 1;`
- 表达式 Expressions：计算并产生一个值，例如`{ 3 + 1 }`(`let y = { 3 + 1 }; // y = 4`)
    - 表达式的结尾**没有分号(`;`)**

```rust
fn five() -> i32 {
    // 不能有分号，有分号就变成了语句，又因为语句没有返回，所以函数隐式返回单元（unit，空白元组），
    // 与函数声明中的`i32`返回值类型不匹配
    5
}
```

代码块的值就是其中最后一个表达式的值

多值返回：返回一个元组(tuple)

### 注释 Comments

行注释：`//`

### 控制流 Control Flow

#### if

```rust
fn f() {
    let number = 3;

    if number < 5 { // 条件必须是一个bool值
        println!("condition was true");
    } else { // 也可以跟else if
        println!("condition was false");
    }

    // if是一个表达式（有返回值），所以可以像下面这样写：
    // 理解：实际执行时，会根据条件，将整个if替换成一个代码块（例如在本例中，条件为true，整个if替换为`{5}`）；
    // 同时，因为要使用if表达式的值，所以要求if的两个代码块的值的类型相同（例如本例中的`{5}`和`{6}`）
    let x = if number != 0 { 5 } else { 6 };
}
```

#### 循环

包含`loop`、`while`、`for`

##### loop

一直执行，直到手动退出

```rust
fn f() {
    let mut counter = 0;

    let result = loop {
        counter += 1;

        if counter >= 10 {
            break counter; // break可以退出循环，break表达式也可以带值
        }
    };

    'loop_label: loop {
        loop {
            break 'loop_label; // 在嵌套循环中，break和continue默认作用于最内层循环，可以通过循环标签(loop label)指定作用层级
        }
    }
}
```

##### while

每次执行前判断，条件为`ture`时执行

```rust
fn f() {
    let mut number = 3;

    while number > 0 { // `loop {}`就像是`while true {}`
        number -= 1;
    }
}
```

##### for

遍历一个数组/集合

```rust
fn f() {
    let a = [10, 20, 30, 40, 50];

    for element in a {
        println!("the value is: {element}");
    }

    for number in (1..4).rev() { // 打印：3! 2! 1!
        println!("{number}!");
    }
}
```

## 所有权 Ownership

所有权(ownership)是rust管理内存的一组规则。所有的应用程序都必须管理它运行时使用的计算机内存。
一些语言使用gc(garbage collection)；一些语言要求程序员亲自分配和释放内存，rust则是通过所有权系统管理内存。
编译器会根据一系列规则进行检查，能编译就说明程序能正常管理它运行时使用的内存。

### 堆栈 Heap & Stack

堆/栈都是代码在运行时可供使用的内存

对比：

- 栈是有组织的结构，后进先出、数据大小已知且固定
- 堆：
    - 存入数据：内存分配器找到一块足够大的空间，标记为已使用，返回一个指向该空间的指针
    - 访问数据：通过指针
- 栈的使用比堆更快，因为存入时不需要找空间、访问时也不需要通过指针

所有权系统做的事情：（管理堆数据）

- 跟踪哪部分代码正在使用堆上的哪些数据
- 减少堆上的重复数据的数量
- 清理堆上不再使用的数据

### 所有权规则

- rust中的每个值都有一个所有者(owner)
- 一个值在任一时刻有且只有一个所有者
- 所有者离开作用域时，值被丢弃

#### 变量作用域

从声明位置直到代码块结束

#### 内存与分配

其他语言中可能有浅拷贝（只拷贝指针的值，实际上共用底层数据），rust面对类似情况，会直接将原始值无效：

```rust
fn f() {
    let s1 = String::from("hello");
    // 此时s1、s2都能使用，它们在堆上使用各自的底层数据（产生堆上复制）
    // 使用`.clone()`函数表示，你知道该操作可能很慢、很消耗资源，并且对此有预期
    let s2 = s1.clone();
    // 触发移动(move)，因为String类型没有实现Copy trait。从这里开始s1失效
    // 如果一个类型实现了Copy trait，那么在赋值之后，旧的变量依旧有效；
    // rust不允许自身或其任何部分实现了Drop trait的类型使用Copy trait(?)
    let s2 = s1;

    println!("{s1}, world!"); // &#10006;

    let x = 5;
    let y = x; // 此时x、y均有效，因为整型大小确定，它们都被分配在栈上，与上面堆的例子不同
}
```

#### 所有权与函数

```rust
fn main() {
    let s = String::from("hello");  // s 进入作用域

    takes_ownership(s);             // s 的值移动到函数里，所以到这里不再有效

    let x = 5;                      // x 进入作用域

    makes_copy(x);                  // x 应该移动函数里，但 i32 是 Copy 的，

    println!("{}", x);              // 所以在后面可继续使用 x

} // 这里，x 先移出了作用域，然后是 s。但因为 s 的值已被移走，没有特殊之处

fn takes_ownership(some_string: String) { // some_string 进入作用域
    println!("{some_string}");
} // 这里，some_string 移出作用域并调用 `drop` 方法。占用的内存被释放

fn makes_copy(some_integer: i32) { // some_integer 进入作用域
    println!("{some_integer}");
} // 这里，some_integer 移出作用域。没有特殊之处
```

#### 返回值与作用域

```rust
fn main() {
    let s1 = gives_ownership();        // gives_ownership 将它的返回值传递给 s1

    let s2 = String::from("hello");    // s2 进入作用域

    let s3 = takes_and_gives_back(s2); // s2 被传入 takes_and_gives_back, 它的返回值又传递给 s3
} // 此处，s3 移出作用域并被丢弃。s2 被 move，所以无事发生。s1 移出作用域并被丢弃

fn gives_ownership() -> String {       // gives_ownership 将会把返回值传入
    // 调用它的函数

    let some_string = String::from("yours"); // some_string 进入作用域

    some_string                        // 返回 some_string 并将其移至调用函数
}

// 该函数将传入字符串并返回该值
fn takes_and_gives_back(a_string: String) -> String {
    // a_string 进入作用域

    a_string  // 返回 a_string 并移出给调用的函数
}
```

### 引用与借用 Reference & Borrow

向函数传值的引用，不会改变值的所有权

```rust
fn main() {
    let s1 = String::from("hello");

    let len = calculate_length(&s1);

    println!("The length of '{s1}' is {len}.");
}

fn calculate_length(s: &String) -> usize { // 引用默认不可变
    s.len()
} // 这里，s 离开了作用域。但因为它并不拥有引用值的所有权，所以什么也不会发生

fn mutable_reference(s: &mut String) -> usize {
    s.push_str(", world!")
}
```

创建一个引用的过程成为借用 borrowing

```rust
fn f() {
    let s = String::from("hello");

    // 引用的创建可能导致数据竞争 data race，所以不允许同时存在**包含可变引用的复数(>=2)个引用**
    let r1 = &s; // &#10004; 
    let r2 = &s; // &#10004; 不可变引用可以随意创建，就像多个只读线程不会触发数据竞争
    let r3 = &mut s; // &#10006; 之前有引用了，不能创建同一个值的可变引用，因为一读一写也会触发数据竞争

    println!("{r1},{r2},{r3}");

    let r4 = &mut s; // 引用的作用域是**从声明位置到最后一次使用的位置**，所以这里可以继续创建s的引用
}
```

编译器保证一个引用永远不会变成悬垂引用

- 悬垂指针 dangle pointer：指向一块已经释放的内存的指针
- 引用必须总是有效的

### Slice

slice可以引用集合中一段连续的元素序列，是一种引用，不拥有所有权

```rust 
fn f() {
    let s = String::from("hello world");

    let hello = &s[0..5]; // 等价于`[..5]`，省略开始位置，表示从头开始；与`0`含义一致
    let slice1 = &s[6..]; // 省略结束位置，表示直到末尾
    let slice2 = &s[..]; // 同时省略开始和结束位置，表示完整的s
}
```

## 结构体 Struct

```rust
struct User {
    id: i32,
    user_name: String,
    last_login: i64,
}

impl User { // 为结构体定义方法
    // 等价于`self: &Self`，其中`Self`是结构体别名，可以改成`User`。
    // self也可以改，但是改了就不能用`user1.print()`了，只能使用`User::print(&user1)`
    fn print(&self) {
        println!("id: {}", self.id);
    }
}

fn main() {
    let mut user1 = User {
        id: 100,
        user_name: String::from("mario"),
        last_login: 1700000,
    };

    user1.id = 200; // 只有mut struct才可以修改

    let user2 = User { id: 200, ..user1 }; // 使用一个变量的值构造新的对象
    println!("{}", user1.id); // &#10006;
    println!("{}", user1.user_name); // &#10004; 移动使用后的值（这个值已经给user2了）

    user2.print(); // 调用结构体方法
}

/* 元组结构体 tuple struct */

struct Color(i32, i32, i32);
struct Point(i32, i32, i32); // 不同的元组结构体是不同的类型，哪怕它们拥有相同的字段

fn main() {
    let black = Color(0, 0, 0);
    let origin = Point(0, 0, 0);
    let Point(x, y, z) = origin; // 解构元组结构体需要显式提供结构体名
}

/* 类单元结构体 unit-like struct */
// 前情提要：unit即空元组`()`
// 此类结构体常用于*想在一个类型上实现trait，又不需要这个类型存储任何数据*的场景

struct AlwaysEqual;

fn main() {
    let subject = AlwaysEqual;
}
```

### 所有权

如果结构体字段均为拥有所有权的类型（例如`String`拥有所有权，`&str`没有），则结构体有效期间，其数据也是有效的。
结构体可以包含没有所有权的字段（例如引用），但此时需要使用生命周期 life time

## 枚举和模式匹配 enumerate

```rust
#[derive(Debug)]
enum IPAddressType {
    IPv4(String), // 枚举也可以绑定值和方法
    IPv6,
}

impl IPAddressType {
    fn to_string(&self) {
        println!("{:?}", self);
    }
}

fn main() {
    let v4 = IPAddressType::IPv4(String::from("127.0.0.1"));
    v4.to_string()
}
```

rust没有空值，类似的有`Option<T>`枚举

```rust
enum Option<T> {
    None,
    Some(T),
}

fn f() {
    // `Option<T>`类型非常常用，所以不需要显式引入、也不需要`Option::`前缀就可以使用`None`和`Some`
    let some_number = Some(5); // Option<i32>
    let some_char = Some('e'); // Option<char>

    let absent_number: Option<i32> = None;
}
```

### match

参考go的switch，match可以将一个值与一系列模式相比较

```rust 
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
    SomeOthers,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        } // 分支之间用逗号`,`分隔
        Coin::Nickel => 5, // 每个分支包含一个模式、分隔符`=>`和一些代码
        Coin::Dime => 10, // 如果一个分支的模式可以匹配match的值，则执行对应分支的代码，执行结果将作为match表达式的值
        Coin::Quarter => 25,
        // match必须覆盖所有分支，或者包含default分支
        // 如果不想使用default模式匹配到的值，可以使用`_`占位符
        _ => 0,
    }
}
```

## 包、Crates与模块

- 包（Packages）：Cargo 的一个功能，它允许你构建、测试和分享 crate。
    - 一系列crate的组合（至少一个，library crate至多一个）
- Crates：一个模块树，可以产生一个库(library crate)或可执行文件(binary crate)。
    - crate是Rust编译器每次处理的最小代码单位
    - binary crate有`main`函数，library crate没有；通常我们说的rust库是library crate，这与其他语言中library含义一致
- 模块（Modules）和 use：允许你控制作用域和路径的私有性。
- 路径（path）：一个为例如结构体、函数或模块等项命名的方式。

module工作方式：

- 从 crate 根节点开始: 当编译一个 crate, 编译器首先在 crate 根文件（通常，对于一个库 crate 而言是 src/lib.rs，
  对于一个二进制 crate 而言是 src/main.rs）中寻找需要被编译的代码。
- 声明模块: 在 crate 根文件中，你可以声明一个新模块；比如，用 mod garden; 声明了一个叫做 garden 的模块。
  编译器会在下列路径中寻找模块代码：
    - 内联，用大括号替换 mod garden 后跟的分号
    - 在文件 src/garden.rs
    - 在文件 src/garden/mod.rs
- 声明子模块: 在除了 crate 根节点以外的任何文件中，你可以定义子模块。比如，你可能在 src/garden.rs 中声明 mod vegetables;。
  编译器会在以父模块命名的目录中寻找子模块代码：
    - 内联，直接在 mod vegetables 后方不是一个分号而是一个大括号
    - 在文件 src/garden/vegetables.rs
    - 在文件 src/garden/vegetables/mod.rs
- 模块中的代码路径: 一旦一个模块是你 crate 的一部分，你可以在隐私规则允许的前提下，从同一个 crate 内的任意地方，通过代码路径引用该模块的代码。
  举例而言，一个 garden vegetables 模块下的 Asparagus 类型可以通过 `crate::garden::vegetables::Asparagus` 访问。
- 私有 vs 公用: 一个模块里的代码默认对其父模块私有。为了使一个模块公用，应当在声明时使用 `pub mod` 替代 `mod`。
  为了使一个公用模块内部的成员公用，应当在声明前使用pub。
- use 关键字: 在一个作用域内，use关键字创建了一个项的快捷方式，用来减少长路径的重复。
  在任何可以引用 `crate::garden::vegetables::Asparagus` 的作用域，你可以通过 `use crate::garden::vegetables::Asparagus;`
  创建一个快捷方式，然后你就可以在作用域中只写 Asparagus 来使用该类型。

所有标识符默认是父模块私有的，可以通过`pub`定义成公有的。

- 公有结构体的每个字段默认私有，需要额外设置为共有
- 公有枚举的每个变体默认公有

use：根据以上规则，每次调用一个模块以外的标识符，都要写形如：`crate::garden::vegetables::Asparagus`的代码，这太长了且没有必要，
所以通常会使用：`use crate::garden::vegetables;`，然后在整个作用域内直接写：`let v = vegetables::Asparagus`

- 虽然可以use到函数，但是人们习惯只use到一个结构体/枚举等对象
- `pub use`：重导出 re-export，只使用use导入的内容对外部模块也是不可见的，使用pub可以让外部对其可见

## 常见集合

包含Vector、String和Hash Map，它们都存储在堆上，适用所有权规则

### 向量 Vector

```rust
fn f() {
    let v: Vec<i32> = Vec::new();
    let mut v = vec![1, 2, 3]; // 带着值初始化

    v.push(10);

    /* 读vec的两种方法 */

    let i: &i32 = &v[1]; // 索引不存在：panic
    let i: Option<&i32> = v.get(1); // get返回Option<&T>类型，如果索引不存在，返回`None`

    /* 所有权规则 */

    let mut v = vec![1, 2, 3, 4, 5];

    let first = &v[0];

    // 因为push可能因空间不足而重新分配，所以v.push会持有一个可变引用
    // 如果后面还使用了不可变引用，那么在v.push这一行执行期间，就同时持有了可变和不可变的引用，这是不行的
    // v.push(6);

    println!("The first element is: {first}");

    /* 遍历 */

    for i in &mut v {
        *i += 50; // 也是因为所有权规则，不能在for里面删除vec元素
    }

    /* 存放一组类型不同的值：使用枚举 */

    enum SpreadsheetCell {
        Int(i32),
        Float(f64),
        Text(String),
    }

    let row = vec![
        SpreadsheetCell::Int(3),
        SpreadsheetCell::Text(String::from("blue")),
        SpreadsheetCell::Float(10.12),
    ];
}
```

### 字符串 String

在实现上，String是一个带有一些额外保证、限制和功能的`vec<u8>`

```rust
fn main() {
    let s = String::new();
    // to_string()可用于任何实现了Display trait的类型，例如字符串字面量
    let s = "string literal".to_string();
    let mut s = String::from("string literal"); // 与上一行等价

    s.push_str("push str");

    // let s = s + &s; // 编译失败，String的加号(+)会获取加号左值的所有权
    let s = format!("{s}-{s}"); // format!宏不会获得值的所有权，建议用这种方式

    /* 索引 */

    // String不支持根据索引获取字符，因为一个UTF-8字符可能占据多个字节，此时根据索引获取的结果没有意义，rust直接禁止这种行为
    let hello = "Здравствуйте";
    // let answer = &hello[0]; // error: str不能使用整型索引
    let s = &hello[0..4]; // 可以使用range索引，但是如果range没算好、对一个字节的一部分进行索引，则会panic

    // 操作字符串可以交给rust遍历：（可以按照字符 char 或字节 byte 遍历）
    for c in "Зд".chars() { // or `xx.bytes()`
        println!("{c}"); // print: З д / 208 151 208 180
    }
}
```

### 哈希映射 Hash map

```rust
use std::collections::{hash_map, HashMap};

fn f() {
    let mut scores = HashMap::new();

    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Yellow"), 50);

    /* 访问 */

    // get 返回`Option<&V>`类型
    // copied() 返回`Option<V>`类型
    // unwrap_or() 返回对应key的value，如果key不存在则返回默认值
    let v = scores.get("Blue").copied().unwrap_or(0);

    for (key, value) in &scores {
        println!("{key}: {value}"); // 打印不保证顺序
    }

    /* 所有权 */

    // 将持有所有权的值插入map时，会同步转移所有权
    let field_name = String::from("Favorite color");
    let field_value = String::from("Blue");

    let mut map = HashMap::new();
    map.insert(field_name, field_value);
    // let v = field_name; // &#10006 使用移动后的值

    /* 更新map */

    // 向已存在的key插入，默认覆盖历史value
    // scores.insert(String::from("Blue"), 10); // warning：重复的key
    // scores.insert(String::from("Blue"), 25);

    // 如果key不存在则插入、存在则无行为
    // entry().or_insert() 返回目标值的可变引用 &mut （不管是新创建的值，还是已经有的值）
    scores.entry(String::from("Blue")).or_insert(100);
    *scores.entry(String::from("Blue")).or_insert(200) += 1; // 根据旧值更新
}
```

## 错误处理 error

主动调用`panic!()`宏或执行rust明确禁止的行为（例如数组访问索引越界）会触发panic并停止程序

```rust
fn read_username_from_file() -> Result<String, io::Error> {
    let username_file_result = File::open("hello.txt");

    let mut username_file = match username_file_result {
        Ok(file) => file,
        Err(e) => return Err(e),
    };

    let mut username = String::new();

    match username_file.read_to_string(&mut username) {
        Ok(_) => Ok(username),
        Err(e) => Err(e),
    }
}

// 与上一个函数功能相同
fn read_username_from_file_2() -> Result<String, io::Error> {
    // `?`：在 `Result<T,E>`/`Option<T>`/其他实现了`FromResidual`类型的值 后面使用（以下内容以`Result<T,E>`为例）
    // - 如果值是`OK`，表达式返回OK中的值，程序继续执行
    // - 如果值是`Err`，Err将作为整个函数的返回值提前返回，
    //   Err会经过from函数，将错误类型转换为函数返回类型中定义的错误类型
    // 使用限制：需要返回类型中定义的错误类型能兼容表达式的错误类型
    let mut username = String::new();
    File::open("hello.txt")?.read_to_string(&mut username)?;
    Ok(username)
}
```

虽然文档中列出了一些建议使用`panic!`的场景，不过我认为这些都可以使用`Result<T,E>`；
文档指出，程序处理不了的状态应使用`panic!`，但我感觉文档中的例子都可以处理，不知道是不是还没有从go转换过来思路的原因

## 泛型、trait和生命周期 Generic Trait and Lifetime

### 泛型

rust的泛型没有运行时开销，因为在编译期间，rust会将泛型单态化 monomorphization

以下是一些常见写法：

```rust
fn largest<T: std::cmp::PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];

    for item in list {
        if item > largest {
            largest = item;
        }
    }

    largest
}

struct Point<T> {
    x: T,
    y: T,
}

enum Option<T> {
    Some(T),
    None,
}

impl<T> Point<T> {
    fn x(&self) -> &T {
        &self.x
    }
}

impl Point<f32> {
    fn distance_from_origin(&self) -> f32 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}
```

### trait

trait类似其他语言中的接口 interface

```rust
pub trait Summary {
    fn summarize_author(&self) -> String;

    fn summarize(&self) -> String {
        format!("(Read more from {}...)", self.summarize_author()) // trait的默认实现中可以使用其他trait
    }
}

pub struct NewsArticle {
    pub headline: String,
    pub location: String,
    pub author: String,
    pub content: String,
}

impl Summary for NewsArticle {
    fn summarize_author(&self) -> String {
        format!("{}", self.author)
    }
}

// 将trait作为函数输入参数
// pub fn notify(item1: &impl Summary, item2: &impl Summary) {
// pub fn notify<T: Summary>(item1: &T, item2: &T) {
// pub fn notify(item: &(impl Summary + Display)) {

// 当输入类型复杂或多样时，可以使用`where`
// fn some_function<T: Display + Clone, U: Clone + Debug>(t: &T, u: &U) -> i32 {
// fn some_function<T, U>(t: &T, u: &U) -> i32
// where
//     T: Display + Clone,
//     U: Clone + Debug,
// {

// 将trait作为输出参数也一样，只是如果有多个分支，每个分支返回的类型要一样(? 18章)
```

### 生命周期 lifetime

生命周期通常和值的作用域保持一致，生命周期注解并不改变生命周期，而是描述多个引用的生命周期之间的关系

生命周期注解是函数签名的一部分，它只在函数签名中有所体现

```rust
// 这个函数编译不过，因为编译器不知道返回值借用的是x还是y(?没看懂，可能是x的生命周期在else之前，检查工具不管这种情况吗)
fn longest(x: &str, y: &str) -> &str {
    if x.len() > y.len() { x } else { y }
}

// 添加生命周期注解后可以编译
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() { x } else { y }
}

fn f() {
    let s: &'static str = "with static lifetime"; // static生命周期表示存活于整个程序期间
}
```

### 同时使用泛型、trait、生命周期

```rust
use std::fmt::Display;

fn longest_with_an_announcement<'a, T>(
    x: &'a str,
    y: &'a str,
    ann: T,
) -> &'a str
where
    T: Display,
{
    println!("Announcement! {ann}");
    if x.len() > y.len() { x } else { y }
}
```

## 测试 Test

我写go的过程中总结出了类似的经验：在代码旁边编写单元测试保证功能正确执行、在单独路径编写集成测试验证程序对外提供的接口

`cargo new [folder name] --lib`
`cargo test`
`cargo test [func name]` // 测试单个函数，函数名可以只是片段
`cargo test -- --test-threads=1` // 取消并发，只使用一个线程执行测试
`cargo test -- --ignored` // 只运行正常测试期间会被忽略的函数
`cargo test -- --include-ignored` // 运行全部测试

```rust
pub fn add(left: u64, right: u64) -> u64 {
    left + right
}

#[cfg(test)] // 告诉编译器：只有在执行`cargo test`的时候才编译和运行这部分代码
mod tests {
    use super::*;

    #[test]
    // #[should_panic(expected = "less than or equal to 100")]
    #[ignore] // 默认不执行
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}
```

## 迭代器与闭包 Iterators and Closures

闭包可以简单的理解成不具名函数，同时闭包可以直接访问定义位置的其他值

迭代器：遍历一个序列中的每个元素，对其执行某些操作

- 遍历由程序控制，程序员只需要编写想要执行的操作
- 可以执行的操作包括提前终止遍历

迭代器和for循环做几乎一样的事情，使用迭代器几乎不会有性能损失，比起写for循环，迭代器不用创建新的变量并复制数据

使用迭代器可以使代码量更少、重点更突出，但是要求程序员对所有权有足够的了解

## cargo高级功能

略（我们只是写一个rust加解密工具，不需要把rust代码发布给rust收录，也不需要使用工作空间来组织代码结构等）

## 智能指针 Smart Pointer

```rust
enum List {
    Cons(i32, Box<List>),
    Nil,
}

use crate::List::{Cons, Nil};

fn main() {
    let list = Cons(1, Box::new(Cons(2, Box::new(Cons(3, Box::new(Nil))))));
}
```

## 并发 Concurrency

```rust 
use std::sync::mpsc;
use std::thread;

fn main() {
    let (tx, rx) = mpsc::channel();

    thread::spawn(move || {
        let val = String::from("hi");
        tx.send(val).unwrap();
    });

    let received = rx.recv().unwrap();
    println!("Got: {received}");
}
```

## 异步与多线程 Async Await Future and Stream

## 面向对象编程特性 Object-Oriented Programming

```rust
pub trait Draw {
    fn draw(&self);
}

// 以下两种定义都支持任意实现trait的类型

pub struct Screen {
  pub components: Vec<Box<dyn Draw>>, // vec元素之间的类型可以不同
}

pub struct Screen<T: Draw> {
  pub components: Vec<T>, // 编译期间会单态化，例如可能编译成`Vec<Button>`或`Vec<Image>`等，元素类型相同
}
```

## 模式与匹配 pattern

所有可以使用模式的位置：
- match / if let
- let 
- while let
- for
- 函数参数
