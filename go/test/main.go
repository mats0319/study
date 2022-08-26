package main

import (
	"fmt"
	"reflect"
)

func main() {
	type s struct {
		i int
	}

	var i interface{} = &s{}

	switch v := i.(type) {
	case *s, s:
		fmt.Println(reflect.TypeOf(v))
		//v.i
	}
}
