# 经验积累

## declared but not used before 1.18

在go1.18版本以前，在函数字面量（闭包）内，变量作为左值出现，则视为该变量**已使用**

```go 
var count int

func () {
   count = 10 // 不会提示'count'变量已定义但未使用
}()
```

查找一些资料，原因大致是gc报告变量未使用的规则问题

参考资料：

1. https://segmentfault.com/a/1190000041047545
2. https://github.com/golang/go/issues/49214

## 建议将append结果赋值给同一个变量

`slice = append(slice, item)`

如上代码，建议使用append的slice，接收append的结果

因为append在slice剩余容量足够的情况下，会修改到slice持有的底层数组，导致非预期的结果发生

举例来说，`slice2 = append(slice1, item)`，看上去slice1没有作为左值出现，实际上slice1可能被修改

参考资料：

1. https://gist.github.com/mats9693/19a29266ebbef13ad2432124a8c4772c

## type switch `case a, b:`

```go 
type s struct {
    i int
}

var i interface{} = &s{}

switch v := i.(type) {
case *s, s:
    // v.i // 此时v是interface{}类型，因为case里有多个类型，包括default分支也是这样；这一行代码会在编译期间报错
}
```

结论：尽量不要在type switch的一个case里写多个类型 ^_^
