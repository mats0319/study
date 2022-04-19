# makefile

[reference](https://www.ruanyifeng.com/blog/2015/02/make.html)

编译(compile)：把代码变成可执行文件 构建(build)：先编译这个，还是先编译那个（即编译的安排）

make是一种构建工具，makefile则是make的构建规则（文件）

当你在makefile中，编写了构建`a.txt`的规则，就可以通过`make a.txt`构建`a.txt`文件了

makefile默认文件名为`makefile`，也可以通过命令行指定其他的构建规则文件：`make -f rules.txt` / `make -file=rules.txt`

## 安装(windows)

使用chocolatey安装高版本的make，**以管理员身份打开powershell**，执行命令：

```powershell 
// 安装chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

choco install make
```

## makefile文件格式

makefile文件由一系列规则(rules)构成，每条规则的格式如下：

```text 
<target> : <prerequisites> 
[tab]  <commands>
```

1. target：构建目标，`make [target]`即表示执行target对应的脚本
2. prerequisites：前置条件，可选
3. commands：命令，需要以`tab`开头，可选
4. 前置条件与命令不可同时为空
5. 每条规则定义了一个操作，包含：操作名、操作的前置条件以及操作是如何执行的

### 目标(target)

target一般是文件名，也可以是一个操作的名字：

```makefile 
clean:
     rm *.exe
```

若当前目录存在名为`clean`的文件，则`make clean`不会执行任何操作，因为make因为`clean`文件已存在，  
为了避免这种情况，可以显式定义`clean`为一种操作：

```makefile 
.phony: clean
```

此时`make clean`不会再检查对应文件

`.phony`叫做**内置目标名**，当它们出现在target位置上时，会表示一些特殊的含义，更多内置目标名，  
参考[手册](https://www.gnu.org/software/make/manual/html_node/Special-Targets.html#Special-Targets)

如果make命令没有带一个target，会默认执行makefile中的第一个target

### 前置条件(prerequisites)

前置条件通常是一组文件名，之间用空格分隔，它规定了target什么时候需要重新构建：（满足以下任意一条，target需要重新构建）

1. 任一前置文件不存在
2. 任一前置文件有更新（前置文件的last-modification时间比target的时间戳新）

```makefile 
result.txt: source.txt
    cp source.txt result.txt
    
source.txt:
    echo "this is the source" > source.txt
```

举个例子，上方代码中，连续调用两次`make result.txt`：

1. 第一次会构建`source.txt`和`result.txt`文件
2. 第二次因为`source.txt`存在且无更新，不会执行任何命令

### 命令(commands)

命令解释了如何构建，由一行或多行脚本组成，每行命令开头都要有`tab`，或者通过内置目标`.recipeprefix`声明新的符号

每一行命令都在一个单独的命令行窗口中运行，若想要在一个命令行窗口中执行全部脚本，可以设置内置目标`.oneshell`

### 语法

注释：makefile仅支持使用`#`注释一行

回声(echoing)：make默认在执行命令每条前，打印一次

通配符：支持`*`/`?`/`...`通配符

模式匹配：主要使用`%`，可以简化相同格式的重复构建写法

```makefile 
%.o: %.c
```

假设当前目录下有`f1.c`/`f2.c`两个文件，上面的写法等价于：

```makefile 
f1.o: f1.c
f2.o: f2.c
```

自定义变量：通过等号定义变量，通过`$()`使用变量

```makefile 
txt = Hello World
test:
    @echo $(txt)
```

当通过变量定义变量时，[参考](https://stackoverflow.com/questions/448910/what-is-the-difference-between-the-gnu-makefile-variable-assignments-a)

make提供一系列的内置变量，[参考](https://www.gnu.org/software/make/manual/html_node/Implicit-Variables.html)
