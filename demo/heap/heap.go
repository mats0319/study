package heap

import "github.com/pkg/errors"

// heap binary minimum heap
type heap struct {
    array []int
    size  int
}

func NewHeap(data ...int) *heap {
    heapIns := &heap{
        array: make([]int, 1, len(data)+1), // deprecated first index
    }

    for i := range data {
        heapIns.Push(data[i])
    }

    return heapIns
}

func (h *heap) Pop() (int, error) {
    if h.size < 1 {
        return -1, errors.New("empty heap")
    }

    popValue := h.array[1]

    h.array[1] = h.array[h.size]

    h.size--

    // todo: maintain heap, mainly about heap[1]

    return popValue, nil
}

func (h *heap) Push(v int) {
    h.array = append(h.array, v)

    h.size++

    // todo: maintain heap, mainly about new item
}

func (h *heap) maintainTop() {
    // placeholder
}
