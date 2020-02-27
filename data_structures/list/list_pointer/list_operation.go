package listp

// Insert try to insert value into position given.
// It will return the list back after insert.
// If position is longer than list's length, it inserts at last.
func (n *Node) Insert(value, position int) *Node {
	if n.IsEmpty() {
		return &Node{Value:value}
	}

	h := &Node{Next: n}
	p := h

	for i := 1; i < position && p != nil; i++ {
		p = p.Next
	}

	insertNode := &Node{Value:value}

	if p == nil {
		p = insertNode
	} else if p.Next == nil {
		p.Next = insertNode
	} else {
		insertNode.Next = p.Next
		p.Next = insertNode
	}

	return h.Next
}

// Find try to find target in list.
// It returns first index matched,
// or -1 when list is empty or target is not exist in list.
func (n *Node) Find(target int) int {
	if n.IsEmpty() {
		return -1
	}

	p := n
	index := -1

	for count := 0; p != nil; p = p.Next {
		if p.Value == target {
			index = count
			break
		}

		count++
	}

	return index
}

// Delete try to delete target in list.
// It returns first index matched,
// or -1 when list is empty or target is not exist in list.
func (n *Node) Delete(target int) int {
	if n.IsEmpty() {
		return -1
	}

	p := n
	prev := &Node{Next:p}
	index := -1

	for count := 0; p != nil; p, prev = p.Next, prev.Next {
		if p.Value == target {
			prev.Next = p.Next
			p.Next = nil // todo: 学习go语言gc，了解是否需要这样写，才会回收被删除的节点所使用的内存空间。
			index = count
			break
		}

		count++
	}

	return index
}
