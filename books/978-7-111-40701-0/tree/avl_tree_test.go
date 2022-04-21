package tree

import (
	"fmt"
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewAVLTree(t *testing.T) {
	data := support.GenerateRandomIntSlice(20, 100)

	avlTreeIns := newAVLTree(data...)

	if !isAVLTree(avlTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printAVLTree(avlTreeIns))
		t.Fail()
	}
}

func TestAvlTreeImpl_Find(t *testing.T) {
	data := support.GenerateRandomIntSlice(20, 100, 0)

	avlTreeIns := newAVLTree(data...)

	res := make([]bool, 0, 2)
	res = append(res, avlTreeIns.Find(0))
	res = append(res, avlTreeIns.Find(200))

	expected := []bool{true, false}

	if len(res) != len(expected) {
		t.Logf("unexpected res amount")
		t.Fail()
	}

	for i := range res {
		if res[i] != expected[i] {
			t.Logf("test avl tree find failed, index: %d\n\twant: %t\n\tget: %t\n", i, expected[i], res[i])
			t.Fail()
		}
	}

	if !isAVLTree(avlTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printAVLTree(avlTreeIns))
		t.Fail()
	}
}

func TestAvlTreeImpl_Insert(t *testing.T) {
	data := support.GenerateRandomIntSlice(20, 100)

	avlTreeIns := newAVLTree(data...)

	for i := 0; i < 10; i++ {
		avlTreeIns.Insert(200+i)
	}

	if !avlTreeIns.Find(200) {
		t.Logf("test avl tree insert failed")
		t.Fail()
	}

	if !isAVLTree(avlTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printAVLTree(avlTreeIns))
		t.Fail()
	}
}

func TestAvlTreeImpl_Delete(t *testing.T) {
	specialValues := []int{-10, 0, 10, 200}

	data := support.GenerateRandomIntSlice(20, 100, specialValues...)

	avlTreeIns := newAVLTree(data...)

	for i := range specialValues {
		avlTreeIns.Delete(specialValues[i])
	}

	if avlTreeIns.Find(200) {
		t.Logf("test avl tree delete failed")
		t.Fail()
	}

	if !isAVLTree(avlTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printAVLTree(avlTreeIns))
		t.Fail()
	}
}

func printAVLTree(tree *avlTreeImpl) string {
	if tree == nil {
		return ""
	}

	nodeList := []*avlNode{tree.root}
	res := ""
	for len(nodeList) > 0 {
		nextNodeList := make([]*avlNode, 0, len(nodeList)*2)

		for len(nodeList) > 0 {
			node := nodeList[0]
			nodeList = nodeList[1:]

			if node == nil {
				res += "null "
			} else {
				res += fmt.Sprintf("%d ", node.value)
				nextNodeList = append(nextNodeList, node.left, node.right)
			}
		}

		res += "\n"
		nodeList = nextNodeList
	}

	return res
}

func isAVLTree(tree *avlTreeImpl) bool {
	return isValidNode(tree.root)
}

func isValidNode(node *avlNode) bool {
	if node == nil {
		return true
	} else if node.isLeaf() {
		node.height = 0

		return true
	}

	isLeftValid := isValidNode(node.left)
	isRightValid := isValidNode(node.right)

	if !isLeftValid || !isRightValid {
		return false
	}

	leftHeight := -1
	if node.left != nil {
		leftHeight = node.left.height
	}

	rightHeight := -1
	if node.right != nil {
		rightHeight = node.right.height
	}

	node.height = big(leftHeight, rightHeight) + 1

	balanceFactor := leftHeight - rightHeight

	return -1 <= balanceFactor && balanceFactor <= 1
}
