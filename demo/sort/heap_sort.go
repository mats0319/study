package sort

import (
	"fmt"
	"github.com/mats9693/study/demo/heap"
)

func heapSort(intSlice []int) {
	heapIns := heap.NewBinaryMinimumHeap(intSlice...)

	res := make([]int, 0, len(intSlice))
	for range intSlice {
		v, err := heapIns.Pop()
		if err != nil {
			fmt.Println("> heap pop failed, error: ", err)
			quickSort(intSlice)
			return
		}

		res = append(res, v)
	}

	for i := range intSlice {
		intSlice[i] = res[i]
	}
}
