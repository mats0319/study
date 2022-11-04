package main

import (
    "fmt"
    "regexp"
)

func main() {
    // RE demo
    buf := "func funcName(x int) string {"
    //解析正则表达式，如果成功返回解释器
    reg1 := regexp.MustCompile(`^func (\w+)\((.*)\)(.*){$`)
    if reg1 == nil {
        fmt.Println("regexp err")
        return
    }
    //根据规则提取关键信息
    result1 := reg1.FindStringSubmatch(buf)
    fmt.Println("result1 = ", result1)
}
