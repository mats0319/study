package tree

type binarySearchTree interface {
	Find(key int) (value int, ok bool)
}

type balancedBST interface {
	Find(key int) (value int, ok bool)
	Insert(key int, value int)
	Delete(key int)
}

// IBSTNode for test
type IBSTNode interface {
	IsEmpty() bool
	Key() int
	Left() IBSTNode
	Right() IBSTNode
}
