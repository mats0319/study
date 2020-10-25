package main

import "fmt"

// Conclusion: empty a slice by 'slice[:0]', it's value is not nil and can 'append'
func testNilSlice() {
    s := []int{1,2,3}
    s = s[:0]
    fmt.Println(s, s == nil, append(s, 4)) // false, can append
}

// It's ok if you append nothing to an empty slice
func testAppendNilSlice() {
    s := make([]int, 0)
    s = append(s, ([]int)(nil)...)
    fmt.Println(s, s == nil, len(s)) // false
}
