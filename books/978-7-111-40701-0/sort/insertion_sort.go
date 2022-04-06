package sort

func insertionSort(intSlice []int) {
	if len(intSlice) < 2 {
		return
	}

	insertSlice := make([]int, len(intSlice))
	for i := range intSlice {
		position := 0
		for intSlice[i] > insertSlice[position] && position < i { // find place
			position++
		}

		for index := i; index > position && index-1 >= 0; index-- { // move elements
			insertSlice[index] = insertSlice[index-1]
		}

		insertSlice[position] = intSlice[i]
	}

	for i := range intSlice {
		intSlice[i] = insertSlice[i]
	}

	return
}
