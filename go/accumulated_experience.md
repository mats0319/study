# 经验积累

## declared but not used

在go1.18版本以前，在函数字面量（闭包）内，变量作为左值出现，则视为该变量**已使用**

```go 
var count int

func () {
   count = 10
}()
```

查找一些资料，原因大致是gc报告变量未使用的规则问题

参考资料：

1. https://segmentfault.com/a/1190000041047545
2. https://github.com/golang/go/issues/49214
