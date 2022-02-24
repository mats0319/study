package sort

func InsertSort(intSlice []int) {
	if len(intSlice) < 2 {
		return
	}

	insertSlice := make([]int, len(intSlice))
	for i := range insertSlice {
		position := 0
		for insertSlice[position] <= intSlice[i] && position < i {
			position++
		}

		for index := i; index > position && index > 1; index-- {
			insertSlice[index] = insertSlice[index-1]
		}

		insertSlice[position] = intSlice[i]
	}

	for i := range intSlice {
		intSlice[i] = insertSlice[i]
	}

	return
}
