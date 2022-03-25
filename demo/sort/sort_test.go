package sort

import (
	"testing"
)

func TestInsertSort(t *testing.T) {
	testWrapper(t, insertionSort)
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

func TestHeapSort(t *testing.T) {
	testWrapper(t, heapSort)
}

func BenchmarkInsertSort(b *testing.B) {
	benchmarkTestWrapper(b, insertionSort)
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

func BenchmarkHeapSort(b *testing.B) {
	benchmarkTestWrapper(b, heapSort)
}
