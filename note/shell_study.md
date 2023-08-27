# shell 学习笔记

> 不是系统性的学习shell，仅针对我需要的内容

## 后台执行

`([cmd] &)`

shell命令通常是阻塞的，即它会等待一条命令完成再执行下一条，如果不想等待命令完成，或者命令是阻塞的（例如监听一个端口），
可以使用`&`后台执行、取消阻塞，

但后台执行的命令是挂载到当前命令行窗口的，后台执行的命令，会在关闭shell窗口时停止，如果想让关闭shell窗口不影响后台执行的程序，可以使用`([cmd] &)`

使用括号包裹整条语句，会将命令挂载到系统守护进程上，这样即使关闭shell窗口，也不会影响后台执行的程序

## 查询进程

`ps -aux | grep [filter_str]`

ps：显示进程状态

`|`：管道，将前一条命令的输出，用作下一条命令的输入

grep：根据过滤字符串过滤结果

简写：`pgrep [filter_str]`，ps+grep，返回结果为pid数组（不包含grep命令的pid）

## 结束进程

`kill [pid]`

`kill $(ps -aux | grep [filter_str])`：向符合条件的进程发送信号

kill命令的含义是向进程发送信号，`-9`表示无条件退出，但由进程自行决定是否退出。也因此，kill无法终止系统进程和守护进程

简写：`pkill [filter_str]`，pgrep+kill

## 为文件添加执行权限

`chmod +x [file name]`

## 变量/函数

定义变量的等号前后不能有空格，例如`name="mario"`  
通过`$`符号使用变量，例如`path="$(name)"`

根据编程规范，变量整体应使用双引号包裹，如上例

定义函数要在使用之前，输入参数不需要提现在函数定义上，如：

```shell 
function upgrade_service() {
  cd ./"$1" || exit
  chmod +x ./"$2"
  (./"$2" &)
  cd ..
}

upgrade_service "service_core" "unnamed_plan_service_core"
```
