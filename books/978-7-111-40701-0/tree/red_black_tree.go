package tree

type redBlackTreeImpl struct {
    root *redBlackNode
}

type redBlackNode struct {
	key   int
	value int

	parent *redBlackNode
	left   *redBlackNode
	right  *redBlackNode
}

var _ balancedBST = (*redBlackTreeImpl)(nil)

func (r redBlackTreeImpl) Find(key int) (value int, ok bool) {
    //TODO implement me
    panic("implement me")
}

func (r redBlackTreeImpl) Insert(key int, value int) {
    //TODO implement me
    panic("implement me")
}

func (r redBlackTreeImpl) Delete(key int) {
    //TODO implement me
    panic("implement me")
}
