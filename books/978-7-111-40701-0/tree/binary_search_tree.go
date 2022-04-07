package tree

import "sort"

type binarySearchTreeImpl struct {
	root *bstNode
}

type bstNode struct {
	value int

	left  *bstNode
	right *bstNode
}

var _ binarySearchTree = (*binarySearchTreeImpl)(nil)

func newBST(data []int) *binarySearchTreeImpl {
	sort.Ints(data)

	return &binarySearchTreeImpl{
		root: buildBST(data),
	}
}

func (t *binarySearchTreeImpl) Find(v int) bool {
	isExist := false

	p := t.root
	for p != nil {
		if v == p.value {
			isExist = true
			break
		} else if v < p.value {
			p = p.left
		} else { // v > p.value
			p = p.right
		}
	}

	return isExist
}

func buildBST(data []int) *bstNode {
	if len(data) < 1 {
		return nil
	} else if len(data) == 1 {
		return &bstNode{value: data[0]}
	}

	length := len(data)

	middle := length/2 + length%2

	root := &bstNode{value: data[middle]}

	root.left = buildBST(data[:middle])

	root.right = buildBST(data[middle+1:])

	return root
}
