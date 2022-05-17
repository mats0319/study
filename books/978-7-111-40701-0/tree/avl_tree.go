package tree

type avlTreeImpl struct {
	root *avlTreeNode
}

type avlTreeNode struct {
	key   int
	value int

	height        int
	balanceFactor int

	left  *avlTreeNode
	right *avlTreeNode
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
		t.root = &avlTreeNode{
			key:   key,
			value: value,
		}

		return
	}

	p := t.root
	needCheckBalance := true
	for {
		if key == p.key {
			p.value = value
			needCheckBalance = false

			break
		} else if key < p.key {
			if p.left == nil {
				p.left = &avlTreeNode{
					key:   key,
					value: value,
				}

				break
			}

			p = p.left
		} else { // key > p.key
			if p.right == nil {
				p.right = &avlTreeNode{
					key:   key,
					value: value,
				}

				break
			}

			p = p.right
		}
	}

	if needCheckBalance {
		t.checkBalance(t.root, nil)
	}
}

func (t *avlTreeImpl) Delete(key int) {
	// find 'key'
	p := t.root // t.root = nil is ok
	var pre *avlTreeNode
	for p != nil && key != p.key {
		pre = p

		if key < p.key {
			p = p.left
		} else { // key > p.key
			p = p.right
		}
	}

	// 'key' is not exist
	if p == nil {
		return
	}

	// 'key' exist, and tree has only one node, just empty the tree and no necessary to check balance
	if t.root == p {
		t.root = nil
		return
	}

	// 'key' exist, recurse move value of its higher child to current node until leaf node, del the leaf node
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
		} else {
			p.value = p.right.value
			p = p.right
		}
	}

	if pre.left == p {
		pre.left = nil
	} else {
		pre.right = nil
	}

	t.checkBalance(t.root, nil)
}

func (t *avlTreeImpl) checkBalance(node *avlTreeNode, parent *avlTreeNode) {
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

func (t *avlTreeImpl) doBalance(node *avlTreeNode, parent *avlTreeNode) {
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

func (t *avlTreeImpl) rotateRight(node *avlTreeNode, parent *avlTreeNode) {
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

	// make l.right to node.right
	node.left = l.right

	l.right = node
}

func (t *avlTreeImpl) rotateLeft(node *avlTreeNode, parent *avlTreeNode) {
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

	// make r.left to node.left
	node.right = r.left

	r.left = node
}

func (n *avlTreeNode) isLeaf() bool {
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
