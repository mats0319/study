package test

import (
	"fmt"
	"testing"
)

func TestByteSliceSideEffect(t *testing.T) {
	str := "test byte slice"
	byteSlice := []byte(str)
	modifyByteSlice(byteSlice)

	// pass []byte into a func, modify in func will Not bring back outside, need to re-set to var
	fmt.Println(len(byteSlice) == len(str))
}

func modifyByteSlice(bytes []byte) {
	bytes = append(bytes, " after modify"...)
}

func TestMapSideEffect(t *testing.T) {
	m := make(map[string]struct{})
	modifyMap(m)

	// pass 'map' into a func, modify in func will bring back outside
	fmt.Println(len(m) == 1)
}

func modifyMap(m map[string]struct{}) {
	m["1"] = struct{}{}
}
