package sort

import "fmt"

func CountingSort(intSlice []int) {
	if len(intSlice) < 2 {
		return
	}

	max, min := intSlice[0], intSlice[0]
	for i := 1; i < len(intSlice); i++ {
		if intSlice[i] > max {
			max = intSlice[i]
		}

		if intSlice[i] < min {
			min = intSlice[i]
		}
	}

	length := max - min + 1

	if (min < 0 && length < max) || length > 10000 { // overflow or too long counting slice
		fmt.Println("unmatched method, too long span(between 'min' and 'max' value of slice)")
		return
	}

	count := make([]int, length)

	offset := min

	for i := range intSlice {
		count[intSlice[i]-offset]++
	}

	resIndex := 0
	for i := range intSlice {
		for count[resIndex] == 0 {
			resIndex++
		}

		count[resIndex]--

		intSlice[i] = resIndex + offset
	}

	return
}
