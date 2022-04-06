package tree

type binarySearchTree struct {
	root *treeNode
}

type treeNode struct {
	value int

	left  *treeNode
	right *treeNode
}

func NewBST(data []int) *binarySearchTree {
	return &binarySearchTree{
		root: buildBST(data),
	}
}

func(t *binarySearchTree) Find(v int) bool {
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

func buildBST(data []int) *treeNode {
	if len(data) < 1 {
		return nil
	} else if len(data) == 1 {
		return &treeNode{value: data[0]}
	}

	length := len(data)

	middle := length/2 + length%2

    root := &treeNode{value: data[middle]}

    root.left = buildBST(data[:middle])

    root.right = buildBST(data[middle+1:])

    return root
}
