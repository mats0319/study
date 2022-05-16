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

func isValidRedBlackTreeNode2(node *redBlackTreeNode) bool {
	// layer-order traversal, calc black node amount on road, until leaf node
	// 'node' is valid only all amounts last step are equal
	if node.color == red {
		return false
	}

	type redBlackTreeNodeWrapper struct {
		node            *redBlackTreeNode
		blackNodeAmount int
	}

	rootBlackNodeAmount := -1 // <0: unsigned

	isValid := true
	nodeList := []*redBlackTreeNodeWrapper{{node: node, blackNodeAmount: 1}}
	for len(nodeList) > 0 {
		n := nodeList[0]
		nodeList = nodeList[1:]

		if n.node.left == nil && n.node.right == nil { // leaf node
			if rootBlackNodeAmount < 0 {
				rootBlackNodeAmount = n.blackNodeAmount
			} else if rootBlackNodeAmount != n.blackNodeAmount {
				isValid = false
				break
			}

			continue
		}

		if n.node.left != nil {
			blackNodeAmount := n.blackNodeAmount
			if n.node.left.color == black {
				blackNodeAmount++
			} else if n.node.color == red {
				isValid = false
				break
			}

			nodeList = append(nodeList, &redBlackTreeNodeWrapper{
				node:            n.node.left,
				blackNodeAmount: blackNodeAmount,
			})
		}
		if n.node.right != nil {
			blackNodeAmount := n.blackNodeAmount
			if n.node.right.color == black {
				blackNodeAmount++
			} else if n.node.color == red {
				isValid = false
				break
			}

			nodeList = append(nodeList, &redBlackTreeNodeWrapper{
				node:            n.node.right,
				blackNodeAmount: blackNodeAmount,
			})
		}
	}

	return isValid
}

func printRedBlackTree2(tree *redBlackTreeImpl) string {
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
		fmt.Println(data[i], ": ")
		fmt.Println(printRedBlackTree2(ins))
		ins.Insert(data[i], data[i])
		if !isValidRedBlackTreeNode2(ins.root) {
			fmt.Println()
			fmt.Println(printRedBlackTree2(ins))
			fmt.Println()
		}
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
	r := parent.parent
	isLeftChildP := r.left == parent
	hasBrotherNode := false
	if isLeftChildP {
		hasBrotherNode = r.right != nil
	} else {
		hasBrotherNode = r.left != nil
	}

	// 4.2. if 'parent' do not have brother node, tree now degenerate to a linked list
	if !hasBrotherNode {
		if isLeftChildP {
			if !isLeftChild { // LR situation
				t.rotateLeft(parent)
				parent = r.left // according to 'isLeftChildP'
			}

			r.setColor(red)
			parent.setColor(black)
			t.rotateRight(r)
		} else {
			if isLeftChild {
				t.rotateRight(parent)
				parent = r.right // according to 'isLeftChildP'
			}

			r.setColor(red)
			parent.setColor(black)
			t.rotateLeft(r)
		}

		return
	}

	// 4.3. if 'parent' have brother node
	r.setColor(red)
	r.left.setColor(black)
	r.right.setColor(black)

	// 4.3.2.
	if r == t.root {
		r.setColor(black)
		return
	}

	// 4.3.3.
	if r.parent.color == black {
		return
	}

	// 4.3.4.2.
	n := r
	p = n.parent
	r = p.parent

	isLeftChildP = r.left == p
	isLeftChild = p.left == n

	if (isLeftChildP && r.right.color == black) || (!isLeftChildP && r.left.color == black) {
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

		return
	}

	// 4.3.4.3
	for {
		p.setColor(red)
		p.left.setColor(black)
		p.right.setColor(black)

		n = p
		p = r
		r = p.parent

		if p == t.root {
			p.setColor(black)
			break
		}

		if r.color == black {
			break
		}

		isLeftChildP = r.left == p
		isLeftChild = p.left == n

		if (isLeftChildP && r.right.color == black) || (!isLeftChildP && r.left.color == black) {
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

			break
		}
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
