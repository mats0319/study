package tree

import "fmt"

func printAVL(tree *avlTreeImpl) string {
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
