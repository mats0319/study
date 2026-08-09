# postgreSQL官方文档阅读笔记

> [文档版本：v18.3](https://www.postgresql.org/docs/current/index.html)
> 2026.5.18：主要阅读前两大章，快速阅读后续章节

## SQL

> SQL通常不区分大小写，除了用反引号括起来的内容
> 注意SQL中的单引号和分号
> SQL输入由一系列命令组成，一个命令由一系列词元(token)组成并以分号终止

建表：`CREATE TABLE weather (
city            varchar(80),
temp_lo         int,           -- low temperature
temp_hi         int,           -- high temperature
prcp            real,          -- precipitation
date            date
);`

插入数据：`insert into weather (city, temp_lo, temp_hi, prcp, date)
values ('San Francisco', 43, 57, 0.0, '1994-11-29')`

- 列名可以省略，但许多开发者认为，显式提供列名比依赖隐含顺序更好
- 列可以改变顺序，也可以省略

查询数据：`select * from weather [where xxx];`

- 查询语句包含：**要返回的列**、**从哪些表中检索**和**可选的限定条件**
- 列可以包含表达式，例如：`(temp_hi+temp_lo)/2 as temp_avg`，其中`as`表示对新的列重命名，是可选的
- 去重：`select distinct`
- 多表查询可以在列名前添加表名区分，例如：`select * from weather w join cities c on w.city = c.name`
- 聚合函数：从多行输入中计算出单个结果

### 关键词与标识符

- SQL关键词和标识符必须以字母(`a-z`)或下划线(`_`)开头，后续字符可以是字母、下划线、数字(`0~9`)或美元符号
  （SQL标准禁用美元符号，所以使用它会降低应用可移植性）
- 标识符的长度超出`namedatalen-1`的部分将被截断，默认情况下，该值为`64`，可以在`src/include/pg_config_manual.h`修改
- 关键词和未加引号的标识符不区分大小写
- 定界标识符：使用双引号包裹的任意字符序列（除了编码为0的字符）
    - 定界标识符必然是标识符：`select`是关键词，而`"select"`则是标识符、可以用来表示列名或者表名
    - 定界标识符区分大小写
    - 定界标识符包含任意序列：可以以此构造原本不可能的表名或列名，例如包含空格的名称
    - 折叠规则：SQL规定一个不加引号的标识符应该折叠为大写，而PG将其折叠为小写
        - 举个例子，SQL规则中，`foo`等价于`"FOO"`；而PG中，`foo`等价于`"foo"`

### 常量

- 字符串常量：由单引号(`''`)括起来的任意字符序列
    - 单引号转义：`''`表示输入一个单引号（在字符串常量中），例如`'Mario''s home'`
    - 两个字符串常量中间仅包含空白和至少一个换行符时，它们会被连在一起且被视为一个字符串常量，例如：
    ```
    select 'foo'
    'bar';
    等价于：（这和dart语言一致）
    select 'foobar';
    ```
- 位串常量：二进制或16进制字符串，需要在字符串前添加`B`/`X`，例如：`B'1001'`/`X'0F'`（不区分大小写）

### 注释

单行注释：`--`

### 表

表由行和列组成，一个表的列的数量和顺序相对固定，而行的数量更容易变化，它反映了某一时刻储存了多少数据。

- 一个表可能出现多行完全相同的记录
- 查询结果默认是无序的
- SQL通常不为行提供唯一标识

表的约束：通常来说，列的类型定义就是表的第一层约束，但这往往不能满足使用，此时可以为表定义更细致的约束

- 检查约束：要求一列/一行的值必须满足一个布尔表达式，举个例子：
  ```text
  CREATE TABLE products (
  product_no integer,
  name text,
  price numeric CHECK (price > 0),
  CHECK (price > 0)
  );
  ```
    - 列约束：定义在一列的末尾，表示对该列数据应用该约束
    - 表约束：定义在单独的一行，表示对表的所有行应用该约束
    - 列约束可以写成表约束，反过来则不行
- 非空约束：指定列不能有空值
- 唯一约束：指定列或组列不能有重复值
    - 列：单独的一列
    - 组列：多个列划分成一组。组列的唯一约束只要求不能有**组列中所有列的值均相同**的两行记录
    - 唯一约束会在列/组列上创建一个唯一B-树索引
    - 两个空值不被视为相等
- 主键：等价于**唯一约束 + 非空约束**
- 外键：约束指定列（或组列）中的值必须匹配出现在另一个表中的某些行的值，即为指定列设置了一个（动态的）取值范围

数据操作：

- 插入数据`INSERT INTO products (product_no, name, price) VALUES (1, 'Cheese', 9.99);`
    - 插入数据的最小单位是行，即不能只插入几个列、而是要插入完整的一行
    - 一条插入指令可以插入多行记录
- 更新数据`UPDATE products SET price = 10 WHERE price = 5;`
    - 因为行没有唯一标识，所以并不能直接指定更新哪一行，而是要指定一个条件、满足条件的行会被更新（条件为空则更新全部行）
    - 新的列值可以是常量，也可以是标量(scalar)表达式`UPDATE products SET price = price * 1.10;`
    - 一条更新指令可以更新多个列
- 删除数据`DELETE FROM products WHERE price = 10;`
    - 删除数据的最小单位是行
    - 与更新相似，删除也需要指定一个条件，不指定则会删除全部行

### 查询

`[WITH with_queries] SELECT select_list FROM table_expression [sort_specification]`

- 可以只有`select`子句，此时查询变成计算器，例如：`SELECT 3 * 4;`/`SELECT random();`

#### from子句

`FROM table_reference [, table_reference [, ...]]`

- 连接类型：交叉连接和限定连接
    - 交叉连接：`T1 CROSS JOIN T2`，两张表的笛卡尔积
        - `FROM T1 CROSS JOIN T2`等效于`FROM T1 INNER JOIN T2 ON TRUE`也等效于`FROM T1,T2`
    - 限定连接：
      ```text
      T1 { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN T2 ON boolean_expression
      T1 { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN T2 USING ( join column list )
      T1 NATURAL { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN T2
      ```
        - 内连接：inner join，对于T1的每一行，如果T2有满足连接条件的行，则在生成的连接表中生成一行
        - 左外连接：left outer join，对于T1的每一行，如果T2有满足连接条件的行则在生成的连接表中生成一行，没有则生成一行T2的列为空的行
        - 右外连接：right outer join，对于T2的每一行，如果T1有满足连接条件的行则在生成的连接表中生成一行，没有则生成一行T1的列为空的行
        - 全外连接：full outer join = left outer join + right outer join
    - `on`子句是最常用的连接条件表现形式，它接收一个布尔值表达式，如果`on`运算为真，则在生成的连接表中生成一行
        - `on`可以指定与连接不直接相关的条件，不过建议将连接无关条件放在`where`子句（因为`on`的条件在连接前处理，可能导致结果不同）
    - `USING`接收一个共享列名列表，在生成的连接表中废除冗余列
        - 生成的连接表的列构造规则：`on`先添加T1的列、再添加T2的列；`USING`先添加共享列，然后添加T1、T2剩余列
    - `NATURAL`是`USING`的一个特例，它会将所有的共享列加入`USING`
- 表的别名：`FROM table_reference [AS] alias`，其中`as`是可选的
    - 定义别名后，原本的名字不可用
    - 通常情况下，定义别名是为了方便人类阅读，但把一个表连接到它自身时必须使用别名：
      `SELECT * FROM people AS mother JOIN people AS child ON mother.id = child.mother_id;`
    - 如果同时包含别名和连接，可以用括号显式指定运算符优先级：`SELECT * FROM (my_table AS a CROSS JOIN my_table) AS b ...`
- 子查询：`FROM (SELECT * FROM table1) AS alias_name`
    - SQL规定子查询必须提供表别名，PG则允许省略

#### where子句

`WHERE search_condition`，其中的查询条件为任意布尔值表达式

#### select子句

`SELECT a, b, c FROM ...`，`select`可以选择展示`from`子句结果的哪些列，`*`表示展示全部

#### order by子句

该子句在`select`子句结束后执行，对最终结果进行排序

- 这里指的是逻辑上的执行顺序，定义了结果应该是什么样的；实际执行中往往不会按照该顺序执行

```text
SELECT select_list
    FROM table_expression
    ORDER BY sort_expression1 [ASC | DESC] [NULLS { FIRST | LAST }]
             [, sort_expression2 [ASC | DESC] [NULLS { FIRST | LAST }] ...]
```

- 排序表达式可以是表达式，例如：`SELECT a, b FROM table1 ORDER BY a + b, c;`
- 升序表示将较小的值排在前面，而*较小*由`<`运算符定义
    - 空值被认为大于任何非空值

#### limit / offset 子句

```text
SELECT select_list
    FROM table_expression
    [ ORDER BY ... ]
    [ LIMIT { count | ALL } ]
    [ OFFSET start ]
```

这两个子句允许你只查询一部分数据

### 数据类型

数据库支持很多的类型，还可以创建新的类型，但我们通常会在应用层编写处理而不是数据库

类型，参考go语言的概念，一个类型定义了一个取值范围、一组运算符和一组方法，例如：

- 几何类型支持平移、缩放、旋转等操作
- 文本搜索类型更方便全文检索
- json类型自带格式校验与操作函数
- ......

### 函数和操作符

#### like

`string LIKE pattern [ESCAPE escape-character]`
`string NOT LIKE pattern [ESCAPE escape-character]`

- 如果`pattern`中不包含百分号(`%`)或下划线(`_`)，`like`的行为与等号(`=`)相似（在索引、尾随空格等的处理上不同）
- `%`：匹配任意数量的字符
- `_`：匹配单个字符

#### case

```text
CASE WHEN condition THEN result
     [WHEN ...]
     [ELSE result]
END
或
CASE expression
    WHEN value THEN result
    [WHEN ...]
    [ELSE result]
END
```

可以用来实现分组排序：

```text
order by
	case v
		when 100 then 1
		when 0 then 2
		else 3
	end
```

## 总结

主要阅读完前两大章以及快速阅读后续章节，我产生了一个想法：数据库能做很多事情，但是我们似乎很少用数据库去做这些事情？  
通过请教其他人，我了解到这是软件开发的普遍现象，并不只是*我以为*大家都不用

举几个例子：

1. 外键。假设现在需要保证数据之间的关系性（例如创建订单要求付款人一定要在系统中注册过），虽然外键就是被设计出来做这个事情的，
   但似乎人们更倾向于在应用层自行编写相关约束。
   一个原因是微服务可能分数据库，导致根本无法添加外键；
   一个原因是外键作为强一致性，在允许弱一致性的场景中会将一些警告当作错误抛出：
    - 迁移数据时批量导入顺序错误
    - 创建订单和用户注册发向集群的两个数据库，且它们中途没有同步数据
    - 删除用户和创建订单并发来到数据库，数据库可能选择先删除用户
2. 类型和函数。pg原生支持很多类型和函数，例如：枚举、网络地址类型、几何类型、自定义类型；求平均值、三角函数、自定义函数等，
   人们似乎也很少使用。
   一个原因是历史问题，存储过程可以调用java代码，可能导致项目与数据库耦合严重，后来渐渐就不用了；
   一个原因是更换数据库是否方便，pg有很多非标准SQL的特化功能，如果使用了这些功能，在更换数据库的时候会很麻烦
3. 触发器。触发器可以解决软删除后新增记录可能导致唯一性约束冲突的问题（即在包含唯一性约束的列上，新增记录使用与软删除记录相同的值，
   会导致唯一性约束冲突，但在业务逻辑上不应该冲突；触发器可以在删除前/后将该条记录移动到其他表以解决该问题）
   这可能导致在特殊场景中执行结果不符合预期，并且很难想到问题的原因。
   举个例子，在数据库集群中定义有插入触发器，数据同步的策略是从日志同步。
   此时向数据库A发布插入命令，A将包含插入命令、触发器执行两条日志。
   数据库B同步A的触发器执行日志时，会执行一次触发器行为；同步A的插入命令时还会执行一次触发器行为
   （此处假设集群节点设置完全一致，包括触发器的定义）。
   如果触发器行为不是幂等的，将导致未知结果（例如触发器行为是扣费，此时就会多次扣费），并且很难想到问题所在。
