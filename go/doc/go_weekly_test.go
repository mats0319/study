package doc

import (
	"fmt"
	"testing"
)

func TestGoWeekly(t *testing.T) {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	slice := array[2:5]
	slice2 := append(slice, 100)

	fmt.Println(array, len(slice), cap(slice), len(slice2), cap(slice2))
}
