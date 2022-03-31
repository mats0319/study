package heap

import "github.com/pkg/errors"

type binaryMinimumHeapOnLinkedList struct {
	pre *listNode // pre.left is the root of heap
}

type listNode struct {
	value int

	parent *listNode
	left   *listNode
	right  *listNode
	prev   *listNode
	next   *listNode
}

func NewBinaryMinimumHeapOnLinkedList(data ...int) Heap {
	heapIns := &binaryMinimumHeapOnLinkedList{
		pre: &listNode{},
	}

	for i := range data {
		heapIns.Push(data[i])
	}

	return heapIns
}

func (h *binaryMinimumHeapOnLinkedList) Push(v int) {
	if h.isNil() {
		h.pre.left = &listNode{value: v}
		return
	}

	root := h.root()

	left, leftDepth := root.getLeftLeafAndDepth()

	_, rightDepth := root.getRightLeafAndDepth()

	next := &listNode{value: v}

	if leftDepth == rightDepth { // full tree, append to left.left
		left.left = next

		next.parent = left
	} else { // not full, append to left.next
		isLeftChild := true // if 'left' is its parent's left child
		for left.next != nil {
			left = left.next
			isLeftChild = !isLeftChild
		}

		// link 'next' with its parent, doubly
		parentOfNext := left.parent
		if !isLeftChild {
			parentOfNext = parentOfNext.next
			parentOfNext.left = next
		} else {
			parentOfNext.right = next
		}

		next.parent = parentOfNext

		// link 'next' with its prev, doubly
		left.next = next
		next.prev = left
	}

	h.shiftUp(next)
}

func (h *binaryMinimumHeapOnLinkedList) Pop() (int, error) {
	if h.isNil() {
		return 0, errors.New("empty heap")
	}

	root := h.root()

	if root.isLeaf() {
		h.pre.left = nil
		return root.value, nil
	}

	popValue := root.value

	left, leftDepth := root.getLeftLeafAndDepth()

	right, rightDepth := root.getRightLeafAndDepth()

	if leftDepth == rightDepth { // full tree
		left = right

		root.value = left.value // 'left' is last node now

		left.parent.right = nil
		left.prev.next = nil
	} else { // not full
		isLeftChild := true // if 'left' is its parent's left child
		for left.next != nil {
			left = left.next
			isLeftChild = !isLeftChild
		}

		root.value = left.value // 'left' is last node now

		if isLeftChild {
			left.parent.left = nil
		} else {
			left.parent.right = nil
		}

		if left.prev != nil {
			left.prev.next = nil
		}
	}

	h.shiftDown()

	return popValue, nil
}

// shiftUp shift up value of given node
func (h *binaryMinimumHeapOnLinkedList) shiftUp(node *listNode) {
	for node != h.root() && node.value < node.parent.value {
		node.value, node.parent.value = node.parent.value, node.value

		node = node.parent
	}
}

// shiftDown shift root element down, in each step, we swap root and its smaller child
func (h *binaryMinimumHeapOnLinkedList) shiftDown() {
	for node := h.root(); !node.isLeaf(); {
		leftValue := node.value
		if node.left != nil && leftValue > node.left.value {
			leftValue = node.left.value
		}

		rightValue := node.value
		if node.right != nil && rightValue > node.right.value {
			rightValue = node.right.value
		}

		if node.value == leftValue && node.value == rightValue {
			break
		}

		if leftValue < rightValue {
			node.value, node.left.value = node.left.value, node.value
			node = node.left
		} else {
			node.value, node.right.value = node.right.value, node.value
			node = node.right
		}
	}
}

func (h *binaryMinimumHeapOnLinkedList) isNil() bool {
	return h.pre.left == nil
}

func (h *binaryMinimumHeapOnLinkedList) root() *listNode {
	return h.pre.left
}

func (n *listNode) isLeaf() bool { // no situation: has right child but not have left child
	return n.left == nil
}

func (n *listNode) getLeftLeafAndDepth() (*listNode, int) {
	left := n
	depth := 0
	for ; left.left != nil; left = left.left {
		depth++
	}

	return left, depth
}

func (n *listNode) getRightLeafAndDepth() (*listNode, int) {
	right := n
	depth := 0
	for ; right.right != nil; right = right.right {
		depth++
	}

	return right, depth
}
