package sort

import (
	"testing"
)

func TestInsertionSort(t *testing.T) {
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

func TestHeapSort_ImplOnLinkedList(t *testing.T) {
	testWrapper(t, heapSortOnLinkedList)
}

func TestRadixSort(t *testing.T) {
	testWrapper(t, radixSort, -100, -10, -1, 0, 1, 10, 100)
}

func TestRadixSort_Method2(t *testing.T) {
	testWrapper(t, radixSort_2, -100, -10, -1, 0, 1, 10, 100)
}

func BenchmarkInsertionSort(b *testing.B) {
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

func BenchmarkHeapSort_ImplOnLinkedList(b *testing.B) {
	benchmarkTestWrapper(b, heapSortOnLinkedList)
}

func BenchmarkRadixSort(b *testing.B) {
	benchmarkTestWrapper(b, radixSort)
}

func BenchmarkRadixSort_Method2(b *testing.B) {
	benchmarkTestWrapper(b, radixSort_2)
}
