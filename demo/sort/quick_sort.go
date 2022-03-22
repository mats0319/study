package sort

func quickSort(intSlice []int) {
	if len(intSlice) < 2 {
		return
	}

	target := intSlice[0] // select first element as 'target'

	// re-assign element as: [smaller] 'target' [bigger]
	small, big := 0, len(intSlice)-1
	for i := 1; i < len(intSlice) && small != big; {
		if intSlice[i] < target {
			intSlice[small] = intSlice[i]
			small++
			i++
		} else { // if data >= 'target'
			intSlice[i], intSlice[big] = intSlice[big], intSlice[i]
			big--
		}
	}

	intSlice[small] = target // small == big now

	// recursion left side and right side of 'target'
	quickSort(intSlice[:small])
	quickSort(intSlice[small+1:])
}
