package sort

import (
	"github.com/mats9693/utils/support"
	"sort"
	"testing"
)

func testWrapper(t *testing.T, sortFunc func([]int), specialValues ...int) {
	intSlice := support.GenerateRandomIntSlice(100, 100, specialValues...)
	originData := support.DeepCopyIntSlice(intSlice)

	sortFunc(intSlice)

	sortedSlice := support.DeepCopyIntSlice(originData)
	sort.Ints(sortedSlice)

	if !support.CompareOnIntSlice(sortedSlice, intSlice) {
		t.Logf("> Test insert sort failed.\n\torigin data: %v\n\texpected: %v\n\tget     : %v\n",
			originData, sortedSlice, intSlice)
		t.Fail()
	}
}

func benchmarkTestWrapper(b *testing.B, sortFunc func([]int)) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		intSlice := support.GenerateRandomIntSlice(10000, 100000)

		b.StartTimer()

		sortFunc(intSlice)
	}
}
