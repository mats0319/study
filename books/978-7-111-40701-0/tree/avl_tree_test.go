package tree

import (
	"fmt"
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewAVLTree(t *testing.T) {
	data := support.GenerateRandomIntSlice(1000, 1000)

	avlTreeIns := newAVLTree(data...)

	if !isAVLTree(avlTreeIns) {
		t.Logf("data: %v\ntree: \n%s", data, printAVLTree(avlTreeIns))
		t.Fail()
	}
}

func TestAvlTreeImpl_Find(t *testing.T) {
	data := support.GenerateRandomIntSlice(20, 100, 0)

	avlTreeIns := newAVLTree(data...)

	values := []int{0, 200}
	expected := []bool{true, false}

	if len(values) != len(expected) {
		t.Logf("unexpected amount")
		t.Fail()
	}

	for i := range values {
		_, ok := avlTreeIns.Find(values[i])
		if ok != expected[i] {
			t.Logf("test avl tree find failed, index: %d\n\twant: %t\n\tget: %t\n", i, expected[i], ok)
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
		avlTreeIns.Insert(200+i, 200+i)
	}

	if _, ok := avlTreeIns.Find(200); !ok {
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

	if _, ok := avlTreeIns.Find(200); ok {
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

	nodeList := []*avlTreeNode{tree.root}
	res := ""
	for len(nodeList) > 0 {
		nextNodeList := make([]*avlTreeNode, 0, len(nodeList)*2)

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
	return isValidAVLTreeNode(tree.root)
}

func isValidAVLTreeNode(node *avlTreeNode) bool {
	if node == nil {
		return true
	} else if node.isLeaf() {
		node.height = 0

		return true
	}

	isLeftValid := isValidAVLTreeNode(node.left)
	isRightValid := isValidAVLTreeNode(node.right)

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
