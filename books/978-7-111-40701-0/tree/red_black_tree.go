package tree

type redBlackTreeImpl struct {
	root *redBlackTreeNode
}

type redBlackTreeNode struct {
	key   int
	value int

	color redBlackTreeNodeColor

	parent *redBlackTreeNode
	left   *redBlackTreeNode
	right  *redBlackTreeNode
}

type redBlackTreeNodeColor = string

const (
	red   redBlackTreeNodeColor = "red"
	black redBlackTreeNodeColor = "black"
)

var _ balancedBST = (*redBlackTreeImpl)(nil)

func newRedBlackTree(data ...int) *redBlackTreeImpl {
	ins := &redBlackTreeImpl{}

	for i := range data {
		ins.Insert(data[i], data[i])
	}

	return ins
}

func (t *redBlackTreeImpl) Find(key int) (value int, ok bool) {
	p := t.root
	for p != nil {
		if p.key == key {
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

func (t *redBlackTreeImpl) Insert(key int, value int) {
	if t.root == nil {
		t.root = &redBlackTreeNode{
			key:   key,
			value: value,
			color: black,
		}

		return
	}

	p := t.root
	parent := p.parent
	for p != nil {
		if key == p.key {
			break
		} else if key < p.key {
			parent = p
			p = p.left
		} else { // key > p.key
			parent = p
			p = p.right
		}
	}

	// 'key' is exist, just update 'value'
	if p != nil {
		p.value = value
		return
	}

	// 'key' is not exist, insert new node
	newNode := &redBlackTreeNode{
		key:    key,
		value:  value,
		color:  red,
		parent: parent,
	}

	if parent.color == black {
		if parent.left == nil && parent.right == nil { // situation: 1.1
			if key < parent.key {
				parent.left = newNode
			} else {
				parent.right = newNode

				parent.setColor(red)
				newNode.setColor(black)
				t.rotateLeft(parent)
			}
		} else { // situation: 1.2
			parent.right = newNode

			t.insertIntoBlackNode(parent)
		}
	} else { // situation: 2
		if key < parent.key {
			parent.left = newNode
		} else {
			parent.right = newNode
		}

		t.insertIntoRedNode(newNode)
	}
}

func (t *redBlackTreeImpl) Delete(key int) {
	// todo: impl
}

// rotateRight implement: draw structure before and after rotate, summarize all changes
// in function, just mark all changed nodes and do the changes
func (t *redBlackTreeImpl) rotateRight(node *redBlackTreeNode) {
	p := node.parent // maybe nil
	l := node.left
	m := node.left.right // migrate node, maybe nil

	// maintain 'p' and its child
	if p == nil { // 'node' is root node
		t.root = l
	} else {
		if p.left == node { // 'node' is left child of parent
			p.left = l
		} else {
			p.right = l
		}
	}
	l.parent = p

	// maintain 'n' and 'l'
	l.right = node
	node.parent = l

	// maintain 'n' and 'm'
	node.left = m
	if m != nil {
		m.parent = node

		if m.color == red && node.color == red {
			t.insertIntoRedNode(m)
		}
	}

	if l.color == red && (p != nil && p.color == red) {
		t.insertIntoRedNode(l)
	}
}

func (t *redBlackTreeImpl) rotateLeft(node *redBlackTreeNode) {
	p := node.parent // maybe nil
	r := node.right
	m := node.right.left // maybe nil

	// maintain 'p' and its child
	if p == nil { // 'node' is root node
		t.root = r
	} else {
		if p.left == node { // 'node' is left child of parent
			p.left = r
		} else {
			p.right = r
		}
	}
	r.parent = p

	// maintain 'n' and 'r'
	r.left = node
	node.parent = r

	// maintain 'n' and 'm'
	node.right = m
	if m != nil {
		m.parent = node

		if m.color == red && node.color == red {
			t.insertIntoRedNode(m)
		}
	}

	if r.color == red && (p != nil && p.color == red) {
		t.insertIntoRedNode(r)
	}
}

func (n *redBlackTreeNode) setColor(c redBlackTreeNodeColor) {
	if n != nil {
		n.color = c
	}
}

// insertIntoBlackNode node is in the tree, just do balance
func (t *redBlackTreeImpl) insertIntoBlackNode(parent *redBlackTreeNode) {
    // situation: 1.2
    parent.setColor(red)
    parent.left.setColor(black)
    parent.right.setColor(black)

    if parent.parent == nil  {
        parent.setColor(black)
    } else if parent.parent.color == red {
        t.insertIntoRedNode(parent)
    }
}

// insertIntoRedNode node is in the tree, just do balance
func (t *redBlackTreeImpl) insertIntoRedNode(node *redBlackTreeNode) {
	parent := node.parent

	if parent.right == node { // situation: 2.2
		parent.setColor(red)
		node.setColor(black)
		t.rotateLeft(parent)

		parent, node = node, parent
	}

	// situation: 2.1
	parent.parent.setColor(red)
	parent.setColor(black)
	t.rotateRight(parent.parent)

	t.insertIntoBlackNode(parent)
}

var _ IBSTNode = (*redBlackTreeNode)(nil)

func (n *redBlackTreeNode) IsEmpty() bool {
	return n == nil
}

func (n *redBlackTreeNode) Key() int {
	return n.key
}

func (n *redBlackTreeNode) Left() IBSTNode {
	return n.left
}

func (n *redBlackTreeNode) Right() IBSTNode {
	return n.right
}
