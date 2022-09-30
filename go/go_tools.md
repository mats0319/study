# Go tools

[reference](https://go.dev/doc/cmd)

go提供了一套程序来构建和处理go代码，它们通常通过`go`命令调用，例如`go fmt`

## vet

检查go代码，报告可疑结构(suspicious constructs)，例如参数与格式字符串不一致的printf调用

vet使用的启发式方法不保证所有的报告都是真正的问题，但它可以找到编译器未捕获的错误

vet不一定能找到所有的错误，所以不要依赖它判断程序正确性

需要在go module目录执行：`go vet [path]`

我的使用方式：`go vet -json ./... *> go_vet_report.txt`

1. 使用json格式的错误报告，主要是因为默认格式下，如果没有检查到可疑结构，vet不会有任何输出，初学时容易误解成命令没有执行
2. 写文件是一个好习惯
3. 为什么使用`*>`？powershell在执行命令时，会报如下错误，所以使用`*>`将所有输出都重定向到文件

```txt 
go : # github.com/mats9693/unnamed_plan/services/shared/const
所在位置 行:1 字符: 1
+ go vet -json ./... *> a.txt
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : NotSpecified: (# github.com/ma...es/shared/const:String) [], RemoteException
    + FullyQualifiedErrorId : NativeCommandError
 
{}
# github.com/mats9693/unnamed_plan/services/shared/http/response
{}
```

## fmt

格式化代码

`go fmt [flags] [path ...]`
or
`gofmt [flags] [path ...]`

flags:

```txt 
-d
	Do not print reformatted sources to standard output.
	If a file's formatting is different than gofmt's, print diffs
	to standard output.
-e
	Print all (including spurious) errors.
-l
	Do not print reformatted sources to standard output.
	If a file's formatting is different from gofmt's, print its name
	to standard output.
-r rule
	Apply the rewrite rule to the source before reformatting.
-s
	Try to simplify code (after applying the rewrite rule, if any).
-w
	Do not print reformatted sources to standard output.
	If a file's formatting is different from gofmt's, overwrite it
	with gofmt's version. If an error occurred during overwriting,
	the original file is restored from an automatic backup.
	
-cpuprofile filename
	Write cpu profile to the specified file.
```

-r的参数rule的格式：`pattern -> replacement`

我的使用方式：`gofmt -w -s -l .`，修改源文件、简化代码、列举修改了哪些文件

## 问题

1. `gofmt`和`go fmt`结果不一样，`go fmt`功能甚至没有`-w`参数，是什么原因呢
