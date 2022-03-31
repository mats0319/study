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

func TestHeapSort_2(t *testing.T) {
	testWrapper(t, heapSortOnLinkedList)
}

func TestRadixSort_LSD_1(t *testing.T) {
	testWrapper(t, radixSort, -100, -10, -1, 0, 1, 10, 100)
}

func TestRadixSort_LSD_2(t *testing.T) {
	testWrapper(t, radixSortLSD_2, -100, -10, -1, 0, 1, 10, 100)
}

func TestRadixSort_MSD_1(t *testing.T) {
	testWrapper(t, radixSortMSD_1, -100, -10, -1, 0, 1, 10, 100)
}

func TestRadixSort_MSD_2(t *testing.T) {
	testWrapper(t, radixSortMSD_2, -100, -10, -1, 0, 1, 10, 100)
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

func BenchmarkHeapSort_2(b *testing.B) {
	benchmarkTestWrapper(b, heapSortOnLinkedList)
}

func BenchmarkRadixSort_LSD_1(b *testing.B) {
	benchmarkTestWrapper(b, radixSort)
}

func BenchmarkRadixSort_LSD_2(b *testing.B) {
	benchmarkTestWrapper(b, radixSortLSD_2)
}

func BenchmarkRadixSort_MSD_1(b *testing.B) {
	benchmarkTestWrapper(b, radixSortMSD_1)
}

func BenchmarkRadixSort_MSD_2(b *testing.B) {
	benchmarkTestWrapper(b, radixSortMSD_2)
}
