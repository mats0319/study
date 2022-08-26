package tree

import (
	"fmt"
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewRedBlackTree(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
		data := support.GenerateRandomIntSlice(1000, 1000)

		redBlackTreeIns := newRedBlackTree(data...)

		if !isRedBlackTree(redBlackTreeIns) {
			t.Logf("data: %v\ntree: \n%s", data, printRedBlackTree(redBlackTreeIns))
			t.Fail()
		}
	})
}

func TestRedBlackTreeImpl_Find(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
		data := support.GenerateRandomIntSlice(20, 100, 0)

		redBlackTreeIns := newRedBlackTree(data...)

		values := []int{0, 200}
		expected := []bool{true, false}

		if len(values) != len(expected) {
			t.Logf("unexpected amount")
			t.Fail()
		}

		for i := range values {
			_, ok := redBlackTreeIns.Find(values[i])
			if ok != expected[i] {
				t.Logf("test red-black tree find failed, index: %d\n\twant: %t\n\tget: %t\n", i, expected[i], ok)
				t.Fail()
			}
		}

		if !isRedBlackTree(redBlackTreeIns) {
			t.Logf("data: %v\ntree: \n%s", data, printRedBlackTree(redBlackTreeIns))
			t.Fail()
		}
	})
}

func TestRedBlackTreeImpl_Insert(t *testing.T) {
	TestNewAVLTree(t)
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
