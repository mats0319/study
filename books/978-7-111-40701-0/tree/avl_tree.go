package tree

import "sort"

type avlTreeImpl struct {
	root *avlNode
}

type avlNode struct {
	value int

	height        int
	balanceFactor int

	left  *avlNode
	right *avlNode
}

var _ avlTree = (*avlTreeImpl)(nil)

func newAVLTree(data ...int) *avlTreeImpl {
	sort.Ints(data)

	return &avlTreeImpl{
		root: buildAVL(data),
	}
}

func (t *avlTreeImpl) Find(v int) bool {
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

func (t *avlTreeImpl) Insert(v int) {
	if t.root == nil {
		t.root = &avlNode{value: v}
		return
	}

	p := t.root
	for {
		if v < p.value {
			if p.left == nil {
				p.left = &avlNode{value: v}
				break
			}

			p = p.left
		} else {
			if p.right == nil {
				p.right = &avlNode{value: v}
				break
			}

			p = p.right
		}
	}

	t.checkBalance(t.root, nil)
}

func (t *avlTreeImpl) Delete(v int) {
	if t.root == nil {
		return
	}

	p := t.root
	isLeftChild := false
	var pre *avlNode
	for !p.isLeaf() {
		if v == p.value {
			break
		} else if v < p.value {
			pre = p
			p = p.left
			isLeftChild = true
		} else { // v > p.value
			pre = p
			p = p.right
			isLeftChild = false
		}
	}

	// 'v' is not exist
	if v != p.value {
		return
	}

	// find 'v', recurse move its higher child to current node(value only) and del the last leaf node
	pBackup := p
	preBackup := pre
	for !p.isLeaf() {
		leftHeight := -1
		if p.left != nil {
			leftHeight = p.left.height
		}

		rightHeight := -1
		if p.right != nil {
			rightHeight = p.right.height
		}

		pre = p

		if leftHeight > rightHeight {
			p.value = p.left.value
			p = p.left
			isLeftChild = true
		} else {
			p.value = p.right.value
			p = p.right
			isLeftChild = false
		}
	}

	if t.root == p {
		t.root = nil
		return
	}

	if isLeftChild {
		pre.left = nil
	} else {
		pre.right = nil
	}

	t.checkBalance(pBackup, preBackup)
}

func (t *avlTreeImpl) checkBalance(node *avlNode, parent *avlNode) {
	if node == nil {
		return
	}

	leftHeight := -1
	if node.left != nil {
		t.checkBalance(node.left, node)
		leftHeight = node.left.height
	}

	rightHeight := -1
	if node.right != nil {
		t.checkBalance(node.right, node)
		rightHeight = node.right.height
	}

	node.height = big(leftHeight, rightHeight) + 1

	node.balanceFactor = leftHeight - rightHeight

	if node.balanceFactor < -1 || node.balanceFactor > 1 {
		t.doBalance(node, parent)
	}
}

func (t *avlTreeImpl) doBalance(node *avlNode, parent *avlNode) {
	if node.balanceFactor > 1 {
		if node.left.balanceFactor < 0 {
			t.rotateLeft(node.left, parent)
		}

		t.rotateRight(node, parent)
	} else if node.balanceFactor < -1 {
		if node.right.balanceFactor > 0 {
			t.rotateRight(node.right, parent)
		}

		t.rotateLeft(node, parent)
	}
}

func (t *avlTreeImpl) rotateRight(node *avlNode, parent *avlNode) {
	l := node.left

	// maintain relation with parent
	if parent == nil {
		t.root = l
	} else {
		if parent == node {
			parent.left = l
		} else {
			parent.right = l
		}
	}

	// if 'l' has right child, move it to 'node'.left
	if l.right != nil {
		r := l.right

		node.left = r
	}

	// maintain relation between 'node' and 'l'
	l.right = node
}

func (t *avlTreeImpl) rotateLeft(node *avlNode, parent *avlNode) {
	r := node.right

	// maintain relation with parent
	if parent == nil {
		t.root = r
	} else {
		if parent.left == node {
			parent.left = r
		} else {
			parent.right = r
		}
	}

	// if 'r' has left child, move it to 'node'.right
	if r.left != nil {
		l := r.left

		node.right = l
	}

	// maintain relation between 'node' and 'r'
	r.left = node
}

func (n *avlNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func buildAVL(data []int) *avlNode {
	if len(data) < 1 {
		return nil
	} else if len(data) == 1 {
		return &avlNode{value: data[0]}
	}

	length := len(data)

	middle := length/2 + length%2

	root := &avlNode{value: data[middle]}

	root.left = buildAVL(data[:middle])

	root.right = buildAVL(data[middle+1:])

	return root
}

func big(a, b int) (res int) {
	if a > b {
		res = a
	} else {
		res = b
	}

	return
}
