package tree

import (
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewAVLTree(t *testing.T) {
	test100000Times(t, func(t *testing.T) {

		data := support.GenerateRandomIntSlice(1000, 1000)

		avlTreeIns := newAVLTree(data...)

		if !isAVLTree(avlTreeIns.root) {
			t.Logf("data: %v\ntree: \n%s", data, printBST(avlTreeIns.root))
			t.Fail()
		}
	})
}

func TestAvlTreeImpl_Find(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
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

		if !isAVLTree(avlTreeIns.root) {
			t.Logf("data: %v\ntree: \n%s", data, printBST(avlTreeIns.root))
			t.Fail()
		}
	})
}

func TestAvlTreeImpl_Insert(t *testing.T) {
	TestNewAVLTree(t)
}

func TestAvlTreeImpl_Delete(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
		specialValue := 10

		data := support.GenerateRandomIntSlice(20, 100, specialValue)

		avlTreeIns := newAVLTree(data...)

		treeBeforeDel := printBST(avlTreeIns.root)

		avlTreeIns.Delete(specialValue)

		if _, ok := avlTreeIns.Find(specialValue); ok {
			t.Logf("test avl tree delete failed")
			t.Fail()
		}

		if !isAVLTree(avlTreeIns.root) {
			t.Logf("is bst: %t\ndata: %v\ntree before del: \n%s\ntree after del: \n%s", isBST(avlTreeIns.root), dfsBSTNode(avlTreeIns.root), treeBeforeDel, printBST(avlTreeIns.root))
			t.Fail()
		}
	})
}

func isAVLTree(node *avlTreeNode) bool {
	return isBST(node) && isValidAVLTreeNode(node)
}

func isValidAVLTreeNode(node *avlTreeNode) bool {
	if node == nil {
		return true
	} else if node.left == nil && node.right == nil {
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
