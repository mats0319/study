package main

import "fmt"

func main() {
	var r []rune = []rune{'a', 'b'}

	fmt.Printf("%p, %p", r[0], r[1])
}
