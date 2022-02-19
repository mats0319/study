package main

import "fmt"

// 1. 变量地址的改变，并不意味着变量被重新分配，此处可能要结合“函数数据段堆栈”概念理解（即函数在内存中的样子，好像是这么叫吧）
// 举例来说，我有两个结构体s1, s2，先让 s = s1，再让 s = s2；此时s变量显然没有被重新分配，而打印s变量的地址发生了变化（不能是空结构体）
func main() {
    //value := 1_000_000_000_000_000_000
    //
    //fmt.Printf("%T\n", value)
    //fmt.Printf("% x\n", value)
    //fmt.Printf("%#x\n", value)

    mswIns := &MyStringWrapper{}
    fmt.Println(mswIns.String())
}

type MyString int

func (m MyString) String() string {
    return fmt.Sprintf("MyString=%d", m) // Error: will recur forever.
}

type AnotherString struct {
    MyString
}

func (a *AnotherString) String() string {
    return a.MyString.String()
}

type MyStringWrapper struct {
    MyString
    AnotherString
}

func (m MyStringWrapper) String() string {
    return fmt.Sprintf("MyStringWrapper=%x", m.MyString)
}
