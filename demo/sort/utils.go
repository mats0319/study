package sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func testWrapper(t *testing.T, sortFunc func([]int)) {
	intSlice := generateRandomIntSlice(100, 100)
	originData := deepCopyIntSlice(intSlice)

	sortFunc(intSlice)

	sortedSlice := deepCopyIntSlice(originData)
	sort.Ints(sortedSlice)

	if !compareOnIntSlice(sortedSlice, intSlice) {
		t.Logf("> Test insert sort failed.\n\torigin data: %v\n\texpected: %v\n\tget     : %v\n",
			originData, sortedSlice, intSlice)
		t.Fail()
	}
}

func benchmarkTestWrapper(b *testing.B, sortFunc func([]int)) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		intSlice := generateRandomIntSlice(10000, 1000)

		b.StartTimer()

		sortFunc(intSlice)
	}
}

// generateRandomIntSlice generate random int slice, you can set 'max length' and 'max value' of slice
//   @param maxLength: max length of slice, min length of slice is 10
//   @param maxValue: max value of slice element, in fact, slice[i] is random in the area [-'max value', 'max value')
func generateRandomIntSlice(maxLength int, maxValue int) []int {
	if maxLength <= 10 {
		maxLength = 11
	}

	intSlice := make([]int, rand.Intn(maxLength-10)+10) // length: [10, 'max length')
	for i := range intSlice {
		intSlice[i] = rand.Intn(maxValue*2) - maxValue // item value: [-'max value', 'max value')
	}

	return intSlice
}

func deepCopyIntSlice(data []int) []int {
	res := make([]int, len(data))
	for i := range data {
		res[i] = data[i]
	}

	return res
}

func compareOnIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	isEqual := true
	for i := range a {
		if a[i] != b[i] {
			isEqual = false
			break
		}
	}

	return isEqual
}
