package sort

func mergeSort(intSlice []int) {
	if len(intSlice) < 2 {
		return
	}

	middle := len(intSlice) / 2

	mergeSort(intSlice[:middle])
	mergeSort(intSlice[middle:])

	res := make([]int, 0, len(intSlice))
    // merge sorted slice
    {
        l, r :=  0, middle
        for l < middle && r < len(intSlice) {
            if intSlice[l] < intSlice[r] {
                res = append(res, intSlice[l])
                l++
            } else { // l.value >= r.value
                res = append(res, intSlice[r])
                r++
            }
        }

        res = append(res, intSlice[l:middle]...)
        res = append(res, intSlice[r:]...)
    }

    // apply modification on origin data
	for index := range intSlice {
		intSlice[index] = res[index]
	}
}
