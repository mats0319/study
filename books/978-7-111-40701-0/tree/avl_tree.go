package tree

import "sort"

type avlTreeImpl struct {
	root *avlNode
}

type avlNode struct {
	value int

	height int
	parent *avlNode
	left   *avlNode
	right  *avlNode
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
    // todo: impl
}

func (t *avlTreeImpl) Delete(v int) {
    // todo: impl
}

func (t *avlTreeImpl) rotateLeft(node *avlNode) {
	r := node.right

    // maintain relation with parent
	if node.parent == nil {
		t.root = r
        r.parent = nil
	} else {
		p := node.parent
		if p.left == node {
			p.left = r
		} else {
			p.right = r
		}

		r.parent = p
	}

    // if 'r' has left child, move it to 'node'.right
	if r.left != nil {
		l := r.left

		node.right = l
		l.parent = node
	}

    // maintain relation between 'node' and 'r'
	r.left = node
	node.parent = r
}

func (t *avlTreeImpl) rotateRight(node *avlNode) {
    l := node.left

    // maintain relation with parent
    if node.parent == nil {
        t.root = l
        l.parent = nil
    } else {
        p := node.parent
        if p.left == node {
            p.left = l
        } else {
            p.right = l
        }

        l.parent = p
    }

    // if 'l' has right child, move it to 'node'.left
    if l.right != nil {
        r := l.right

        node.left = r
        r.parent = node
    }

    // maintain relation between 'node' and 'l'
    l.right = node
    node.parent = l
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
