package main

import (
	"fmt"
)

// Conclusion: if you send nothing to an optional param, it has nil value with target type slice
func testOptionalParam() {
	test()
}

func test(v ...int) {
	fmt.Println("> Node: show value of v slice, is nil or not?", v, v == nil) // true
}
