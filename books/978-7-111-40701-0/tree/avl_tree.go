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
		t.keepBalance(t.root, nil)
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

	// 'key' exist, find next node of 'p' and remove it
	var child *avlTreeNode
	if p.left != nil && p.right != nil {
		pre = p
		del := p.right
		for del.left != nil {
			pre = del
			del = del.left
		}

		p.key = del.key
		p.value = del.value

		p = del
		child = del.right
	} else { // 'p' is leaf node or has only one child
		if p.left != nil {
			child = p.left
		} else if p.right != nil {
			child = p.right
		}
	}

	if pre.left == p {
		pre.left = child
	} else {
		pre.right = child
	}

	t.keepBalance(t.root, nil)
}

func (t *avlTreeImpl) keepBalance(node *avlTreeNode, parent *avlTreeNode) {
	if node == nil {
		return
	}

	leftHeight := -1
	if node.left != nil {
		t.keepBalance(node.left, node)
		leftHeight = node.left.height
	}

	rightHeight := -1
	if node.right != nil {
		t.keepBalance(node.right, node)
		rightHeight = node.right.height
	}

	node.height = big(leftHeight, rightHeight) + 1

	node.balanceFactor = leftHeight - rightHeight

	if node.balanceFactor < -1 || node.balanceFactor > 1 {
		t.doBalance(node, parent)
		calcHAndB(parent)
	}
}

func calcHAndB(node *avlTreeNode) {
	if node == nil {
		return
	}

	leftHeight := -1
	if node.left != nil {
		calcHAndB(node.left)
		leftHeight = node.left.height
	}

	rightHeight := -1
	if node.right != nil {
		calcHAndB(node.right)
		rightHeight = node.right.height
	}

	node.height = big(leftHeight, rightHeight) + 1

	node.balanceFactor = leftHeight - rightHeight
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

func big(a, b int) (res int) {
	if a > b {
		res = a
	} else {
		res = b
	}

	return
}

var _ IBSTNode = (*avlTreeNode)(nil)

func (n *avlTreeNode) IsEmpty() bool {
	return n == nil
}

func (n *avlTreeNode) Key() int {
	if n == nil {
		return -1
	}

	return n.key
}

func (n *avlTreeNode) Left() IBSTNode {
	if n == nil {
		return nil
	}

	return n.left
}

func (n *avlTreeNode) Right() IBSTNode {
	if n == nil {
		return nil
	}

	return n.right
}
