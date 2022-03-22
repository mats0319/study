package main

import (
	"fmt"
)

// 1. 变量地址的改变，并不意味着变量被重新分配，此处可能要结合“函数数据段堆栈”概念理解（即函数在内存中的样子，好像是这么叫吧）
// 举例来说，我有两个结构体s1, s2，先让 s = s1，再让 s = s2；此时s变量显然没有被重新分配，而打印s变量的地址发生了变化（不能是空结构体）
func main() {
	type s struct {
		i int
	}

	sNew := new(s)
	var sDef *s

	fmt.Printf("%v, %v\n", sNew.i, sDef)
}
