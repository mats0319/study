package tree

type binarySearchTreeImpl struct {
	root *bstNode
}

type bstNode struct {
	key   int
	value int

	left  *bstNode
	right *bstNode
}

var _ binarySearchTree = (*binarySearchTreeImpl)(nil)

func newBST(data []int) *binarySearchTreeImpl {
	bstIns := &binarySearchTreeImpl{}

	for i := range data {
		bstIns.insert(data[i], data[i])
	}

	return bstIns
}

func (t *binarySearchTreeImpl) Find(key int) (value int, ok bool) {
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

func (t *binarySearchTreeImpl) insert(key int, value int) {
	if t.root == nil {
		t.root = &bstNode{
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
				p.left = &bstNode{
					key:   key,
					value: value,
				}

				break
			}

			p = p.left
		} else { // key > p.key
			if p.right == nil {
				p.right = &bstNode{
					key: key,
					value: value,
				}

				break
			}

			p = p.right
		}
	}
}
