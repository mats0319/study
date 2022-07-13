package heap

import "github.com/pkg/errors"

type binaryMinimumHeapOnArray struct {
	data []int
	size int // valid data amount, can use as index directly
}

func NewBinaryMinimumHeapOnArray(data ...int) Heap {
	heapIns := &binaryMinimumHeapOnArray{
		data: make([]int, 1, len(data)+1), // deprecated first element
	}

	for i := range data {
		heapIns.Push(data[i])
	}

	return heapIns
}

func (h *binaryMinimumHeapOnArray) Push(v int) {
	h.data = append(h.data, v)

	h.size++

	h.shiftUp()
}

func (h *binaryMinimumHeapOnArray) Pop() (int, error) {
	if h.size < 1 {
		return 0, errors.New("empty heap")
	}

	popValue := h.data[1]

	h.data[1] = h.data[h.size]

	h.size--

	h.shiftDown()

	return popValue, nil
}

// shiftUp shift last element up to its position
func (h *binaryMinimumHeapOnArray) shiftUp() {
	index := h.size

	for index > 1 && h.data[index] < h.data[index/2] { // loop until root or 'parent' <= 'index'
		h.data[index], h.data[index/2] = h.data[index/2], h.data[index]

		index /= 2
	}
}

// shiftDown shift first element down, in each step, we swap 'index' with its smaller child
func (h *binaryMinimumHeapOnArray) shiftDown() {
	for index := 1; index*2 <= h.size; { // loop until leaf node
		value := h.data[index]

		left := value
		if left > h.data[index*2] {
			left = h.data[index*2]
		}

		right := value
		if index*2+1 <= h.size && right > h.data[index*2+1] { // if exist right child and 'index' > 'right child'
			right = h.data[index*2+1]
		}

		// 'index' is in its position: less than or equal to 'children'(include child not exist situation)
		if left == value && right == value {
			break
		}

		// need more swap, swap 'index' with its smaller child
		if left <= right {
			h.data[index], h.data[index*2] = h.data[index*2], h.data[index]
			index *= 2
		} else {
			h.data[index], h.data[index*2+1] = h.data[index*2+1], h.data[index]
			index = index*2 + 1
		}
	}
}
