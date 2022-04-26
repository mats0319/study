# 数据竞争

数据竞争(data race)，产生条件：

1. 两个goroutine同时访问相同变量
2. 至少有一个是写操作

[参考资料：go内存模型](https://go.dev/ref/mem/)

数据竞争，即我们通常所说的**并发问题**

## 使用

为了解决数据竞争问题，go内置了一个数据竞争检测工具(data race detector)，在go命令中添加`-race`标志即可启用

```cmd 
go test -race mypkg    // to test the package
go run -race mysrc.go  // to run the source file
go build -race mycmd   // to build the command
go install -race mypkg // to install the package
```

race detector只能找出运行时的数据竞争，所以未执行的代码中隐藏的数据竞争不会被检测出来，解决办法如下：

1. 编写覆盖率100%的测试（理想情况）
2. 测试阶段使用`go build -race`编译出的二进制文件

启用数据竞争检测工具，会导致构建出的二进制文件更大、运行更慢，故仅建议在测试环境下启用该工具

## 数据检测工具配置

可以通过`GORACE`环境变量配置检测工具的部分行为，格式为`GORACE="[option1]=[value1] [option2]=[value2]"`

1. log_path：输出位置，默认输出到控制台，配置输出路径后，将直接在后面拼接`.[pid]`，形成文件，保存检测到的数据竞争行为
    1. 举例：`log_path='C:\User\Desktop'`，会在user文件夹生成`Desktop.[pid]`文件；  
       而`log_path='C:\User\Desktop\'`，会在desktop文件夹生成`.[pid]`文件
    2. 如果没有检测到数据竞争行为，则不会创建文件

## 报告结构

当检测工具发现一个数据竞争时，它会打印一个报告，报告包含冲突访问的堆栈跟踪信息(stack traces for conflicting accesses)，  
以及创建冲突访问的goroutine的堆栈信息，例如下面的例子：

```txt 
WARNING: DATA RACE
Read by goroutine 185:
  net.(*pollServer).AddFD()
      src/net/fd_unix.go:89 +0x398
  net.(*pollServer).WaitWrite()
      src/net/fd_unix.go:247 +0x45
  net.(*netFD).Write()
      src/net/fd_unix.go:540 +0x4d4
  net.(*conn).Write()
      src/net/net.go:129 +0x101
  net.func·060()
      src/net/timeout_test.go:603 +0xaf

Previous write by goroutine 184:
  net.setWriteDeadline()
      src/net/sockopt_posix.go:135 +0xdf
  net.setDeadline()
      src/net/sockopt_posix.go:144 +0x9c
  net.(*conn).SetDeadline()
      src/net/net.go:161 +0xe3
  net.func·061()
      src/net/timeout_test.go:616 +0x3ed

Goroutine 185 (running) created at:
  net.func·061()
      src/net/timeout_test.go:609 +0x288

Goroutine 184 (running) created at:
  net.TestProlongTimeout()
      src/net/timeout_test.go:618 +0x298
  testing.tRunner()
      src/testing/testing.go:301 +0xe8
```
