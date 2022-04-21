package tree

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
	ins := &avlTreeImpl{}

	for i := range data {
		ins.Insert(data[i])
	}

	return ins
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

	t.checkBalance(t.root, nil) // todo: optimize, not check all tree
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
			t.rotateLeft(node.left, node)
		}

		t.rotateRight(node, parent)
	} else if node.balanceFactor < -1 {
		if node.right.balanceFactor > 0 {
			t.rotateRight(node.right, node)
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
		if parent.left == node {
			parent.left = l
		} else {
			parent.right = l
		}
	}

	// make l.right to l.node.right
	node.left = l.right

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

	// make r.left to r.node.left
	node.right = r.left

	r.left = node
}

func (n *avlNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func big(a, b int) (res int) {
	if a > b {
		res = a
	} else {
		res = b
	}

	return
}
