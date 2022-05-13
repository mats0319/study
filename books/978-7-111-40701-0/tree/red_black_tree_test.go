package tree

import (
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewRedBlackTree(t *testing.T) {
	data := support.GenerateRandomIntSlice(1000, 1000)

	redBlackTreeIns := newRedBlackTree(data...)

	if !isRedBlackTree(redBlackTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printRedBlackTree(redBlackTreeIns))
		t.Fail()
	}
}

func TestRedBlackTreeImpl_Find(t *testing.T) {
	// todo: impl
}

func TestRedBlackTreeImpl_Insert(t *testing.T) {
	// todo: impl
}

func TestRedBlackTreeImpl_Delete(t *testing.T) {
	// todo: impl
}

func isRedBlackTree(tree *redBlackTreeImpl) bool {
	return isValidRedBlackTreeNode(tree.root)
}

func isValidRedBlackTreeNode(node *redBlackTreeNode) bool {
	// layer-order traversal, calc black node amount on road, until leaf node
	// 'node' is valid only all amounts last step are equal

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
		} else {
			if n.node.left != nil {
				nodeList = append(nodeList, &redBlackTreeNodeWrapper{
					node: n.node.left,
					blackNodeAmount: n.blackNodeAmount+1,
				})
			}
			if n.node.right != nil {
				nodeList = append(nodeList, &redBlackTreeNodeWrapper{
					node: n.node.right,
					blackNodeAmount: n.blackNodeAmount+1,
				})
			}
		}
	}

	return isValid
}

//func printRedBlackTree(tree *redBlackTreeImpl) string {
//	if tree == nil {
//		return ""
//	}
//
//	nodeList := []*redBlackTreeNode{tree.root}
//	res := ""
//	for len(nodeList) > 0 {
//		nextLayer := make([]*redBlackTreeNode, 0, 2*len(nodeList))
//
//		for len(nodeList) > 0 {
//			node := nodeList[0]
//			nodeList = nodeList[1:]
//
//			if node == nil {
//				res += "null "
//			} else {
//				if node.color == black {
//					res += fmt.Sprintf("%d ", node.key)
//				} else { // node.color == red
//					res += fmt.Sprintf("%d_R ", node.key)
//				}
//
//				nextLayer = append(nextLayer, node.left, node.right)
//			}
//		}
//
//		res += "\n"
//		nodeList = nextLayer
//	}
//
//	return res
//}
