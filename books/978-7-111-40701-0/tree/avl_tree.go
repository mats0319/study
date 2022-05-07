package tree

type avlTreeImpl struct {
	root *avlNode
}

type avlNode struct {
	key   int
	value int

	height        int
	balanceFactor int

	left  *avlNode
	right *avlNode
}

var _ balancedBST = (*avlTreeImpl)(nil)

func newAVLTree(data ...int) *avlTreeImpl {
	ins := &avlTreeImpl{}

	for i := range data {
		ins.Insert(data[i], data[i])
	}

	return ins
}

func (t *avlTreeImpl) Find(key int) (value int, ok bool) {
	p := t.root
	for p != nil {
		if key == p.key {
			value = p.value
			ok = true
			break
		} else if key < p.key {
			p = p.left
		} else { // key > p.key
			p = p.right
		}
	}

	return
}

func (t *avlTreeImpl) Insert(key int, value int) {
	if t.root == nil {
		t.root = &avlNode{
			key:   key,
			value: value,
		}

		return
	}

	p := t.root
	for {
		if key == p.key {
			p.value = value
			break
		} else if key < p.key {
			if p.left == nil {
				p.left = &avlNode{
					key:   key,
					value: value,
				}

				break
			}

			p = p.left
		} else { // key > p.key
			if p.right == nil {
				p.right = &avlNode{
					key:   key,
					value: value,
				}

				break
			}

			p = p.right
		}
	}

	t.checkBalance(t.root, nil)
}

func (t *avlTreeImpl) Delete(key int) {
	if t.root == nil {
		return
	}

	p := t.root
	isLeftChild := false
	var pre *avlNode
	for !p.isLeaf() {
		if key == p.key {
			break
		} else if key < p.key {
			pre = p
			p = p.left
			isLeftChild = true
		} else { // key > p.key
			pre = p
			p = p.right
			isLeftChild = false
		}
	}

	// 'key' is not exist
	if key != p.key {
		return
	}

	// 'key' exist, recurse move its higher child to current node(value only) and del the last leaf node
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

	t.checkBalance(t.root, nil)
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
