package sort

import (
	"testing"
)

func TestInsertSort(t *testing.T) {
	testWrapper(t, insertSort)
}

func TestCountingSort(t *testing.T) {
	testWrapper(t, countingSort)
}

func TestQuickSort(t *testing.T) {
	testWrapper(t, quickSort)
}

func TestMergeSort(t *testing.T) {
	testWrapper(t, mergeSort)
}

func BenchmarkInsertSort(b *testing.B) {
	benchmarkTestWrapper(b, insertSort)
}

func BenchmarkCountingSort(b *testing.B) {
	benchmarkTestWrapper(b, countingSort)
}

func BenchmarkQuickSort(b *testing.B) {
	benchmarkTestWrapper(b, quickSort)
}

func BenchmarkMergeSort(b *testing.B) {
	benchmarkTestWrapper(b, mergeSort)
}
