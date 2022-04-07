package tree

type binarySearchTree interface {
	Find(int) bool
}

type avlTree interface {
	Find(int) bool
	Insert(int)
	Delete(int)
}
