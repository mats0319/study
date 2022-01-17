package main

import "fmt"

func main() {
    testOnIntAndIntPointer()
}

func testOnIntAndIntPointer() {
    a, b := funcWrapper(1)
    fmt.Printf("%p, %d, %p, %v\n", a, *a, &b, b)

    a, p1 := funcWrapper(2)
    fmt.Printf("%p, %d, %p, %v\n", a, *a, &b, b)

    p2, b := funcWrapper(10)
    fmt.Printf("%p, %d, %p, %v\n", a, *a, &b, b)

    _, _ = p1, p2
}

func funcWrapper(i int) (*int, int) {
    return &i, i
}
