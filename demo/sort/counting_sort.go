package sort

func countingSort(intSlice []int) {
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

	length := max - min + 1 // calc counting slice length

	counting := make([]int, length)

	offset := min

	for i := range intSlice { // build counting slice
		counting[intSlice[i]-offset]++
	}

	countingIndex := 0
	for i := range intSlice {
		for counting[countingIndex] == 0 { // skip empty and used elements in counting slice
			countingIndex++
		}

		counting[countingIndex]--

		intSlice[i] = countingIndex + offset // recover value according to 'counting index' and 'offset'
	}
}
