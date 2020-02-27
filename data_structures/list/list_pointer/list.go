package listp

import "fmt"

type Node struct {
	Value int
	Next *Node
}

func CreateListWithHeader(value int) *Node {
	return &Node{Next:&Node{Value:value}}
}

func CreateList(value int) *Node {
	return CreateListWithHeader(value).offHeader()
}

func (n *Node) offHeader() *Node {
	return n.Next
}

func (n *Node) IsEmpty() bool {
	return n == nil
}

// IsEqual judge if two list has same value only, not reference address.
func IsEqual(a, b *Node) bool {
	aCopy, bCopy := a, b

	for aCopy != nil && bCopy != nil  {
		if aCopy.Value != bCopy.Value {
			return false
		}

		aCopy = aCopy.Next
		bCopy = bCopy.Next
	}

	return aCopy == nil && bCopy == nil
}

func MakeList(values ...int) *Node {
	h := &Node{}
	p := h

	for i := range values {
		p.Next = &Node{Value:values[i]}
		p = p.Next
	}

	return h.Next
}

func (n *Node) PrintList() {
	if n.IsEmpty() {
		fmt.Println("输入列表为空")
		return
	}

	p := n
	for p.Next != nil {
		fmt.Printf("%d -> ", p.Value)
		p = p.Next
	}

	fmt.Println(p.Value, ".")

	return
}
