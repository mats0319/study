package tree

import "fmt"

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

func printRedBlackTree(tree *redBlackTreeImpl) string {
	if tree == nil {
		return ""
	}

	nodeList := []*redBlackTreeNode{tree.root}
	res := ""
	for len(nodeList) > 0 {
		nextLayer := make([]*redBlackTreeNode, 0, 2*len(nodeList))

		for len(nodeList) > 0 {
			node := nodeList[0]
			nodeList = nodeList[1:]

			if node == nil {
				res += "null "
			} else {
				if node.color == black {
					res += fmt.Sprintf("%d ", node.key)
				} else { // node.color == red
					res += fmt.Sprintf("%d_R ", node.key)
				}

				nextLayer = append(nextLayer, node.left, node.right)
			}
		}

		res += "\n"
		nodeList = nextLayer
	}

	return res
}

func newRedBlackTree(data ...int) *redBlackTreeImpl {
	ins := &redBlackTreeImpl{}

	for i := range data {
		ins.Insert(data[i], data[i])
		fmt.Println(data[i], ": \n", printRedBlackTree(ins))
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

// Insert todo: code according to doc, no test, no optimize
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
	isLeftChild := false
	for p != nil {
		if key == p.key {
			p.value = value
			break
		} else if key < p.key {
			isLeftChild = true
			parent = p
			p = p.left
		} else { // key > p.key
			isLeftChild = false
			parent = p
			p = p.right
		}
	}

	// 'key' is exist
	if p != nil {
		return
	}

	// 'key' is not exist, insert new node
	newNode := &redBlackTreeNode{
		key:    key,
		value:  value,
		color:  red,
		parent: parent,
	}

	if isLeftChild {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	// parent.color is 'black'
	if parent.color == black {
		return
	}

	// 4. parent.color is 'red'
	pp := parent.parent
	isLeftChildP := false
	hasBrotherNode := false
	if pp.left == p && pp.right == nil {
		hasBrotherNode = true
		isLeftChildP = true
	} else if pp.left == nil && pp.right == p {
		hasBrotherNode = true
		isLeftChildP = false
	}

	// 4.2. if 'parent' do not have brother node, tree now degenerate to a linked list
	if !hasBrotherNode {
		if isLeftChildP {
			if !isLeftChild { // LR situation
				t.rotateLeft(parent)
				parent = pp.left // according to 'isLeftChildP'
			}

			pp.setColor(red)
			parent.setColor(black)
			t.rotateRight(pp)
		} else {
			if isLeftChild {
				t.rotateRight(parent)
				parent = pp.right // according to 'isLeftChildP'
			}

			pp.setColor(red)
			parent.setColor(black)
			t.rotateLeft(pp)
		}
	}

	// 4.3. if 'parent' have brother node
	pp.setColor(red)
	parent.setColor(black)
	if isLeftChildP {
		pp.left.setColor(black)
	} else {
		pp.right.setColor(black)
	}

	// 4.3.1.
	if pp == t.root {
		pp.setColor(black)
		return
	}

	// 4.3.2.
	if pp.parent.color == black {
		return
	}

	// 4.3.3.
	n := pp
	p = n.parent
	r := p.parent

	isLeftChildP = r.left == p
	isLeftChild = p.left == n

	if isLeftChildP {
		if !isLeftChild {
			t.rotateLeft(p)
			p = r.left // according to 'isLeftChildP'
		}

		r.setColor(red)
		p.setColor(black)
		t.rotateRight(r)
	} else {
		if isLeftChild {
			t.rotateRight(p)
			p = r.right // according to 'isLeftChildP'
		}

		r.setColor(red)
		p.setColor(black)
		t.rotateLeft(r)
	}
}

func (t *redBlackTreeImpl) Delete(key int) {
	//TODO implement me
	panic("implement me")
}

// rotateRight implement: draw structure before and after rotate, summarize all changes
// in function, just mark all changed nodes and do the changes
func (t *redBlackTreeImpl) rotateRight(node *redBlackTreeNode) {
	p := node.parent // maybe nil
	l := node.left
	reMigrate := node.left.right // maybe nil

	// maintain node.parent
	if p == nil { // 'node' is root node
		t.root = l
	} else {
		if p.left == node { // 'node' is left child of parent
			p.left = l
		} else {
			p.right = l
		}
	}

	// maintain node
	node.parent = l
	node.left = reMigrate // nil 're-migrate' is ok

	// maintain node.left
	l.parent = p // nil 'parent' is ok
	l.right = node

	// maintain re-migrate node, also node.left.right
	if reMigrate != nil {
		reMigrate.parent = node
	}
}

func (t *redBlackTreeImpl) rotateLeft(node *redBlackTreeNode) {
	p := node.parent // maybe nil
	r := node.right
	reMigrate := node.right.left // maybe nil

	// maintain node.parent
	if p == nil { // 'node' is root node
		t.root = r
	} else {
		if p.left == node { // 'node' is left child of parent
			p.left = r
		} else {
			p.right = r
		}
	}

	// maintain node
	node.parent = r
	node.right = reMigrate // nil 're-migrate' is ok

	// maintain node.right
	r.parent = p // nil 'parent' is ok
	r.left = node

	// maintain re-migrate node, also node.left.right
	if reMigrate != nil {
		reMigrate.parent = node
	}
}

func (n *redBlackTreeNode) setColor(c redBlackTreeNodeColor) {
	if n == nil {
		return
	}

	n.color = c
}
