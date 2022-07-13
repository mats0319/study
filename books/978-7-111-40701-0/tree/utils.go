package tree

import (
	"fmt"
	"github.com/mats9693/utils/support"
	"sort"
	"testing"
)

func test100000Times(t *testing.T, testFunc func(*testing.T)) {
	for range [10_0000]struct{}{} {
		testFunc(t)
	}
}

func printBST(root IBSTNode) string {
	if root.IsEmpty() {
		return ""
	}

	nodeList := []IBSTNode{root}
	res := ""
	for {
		nextNodeList := make([]IBSTNode, 0, len(nodeList)*2)
		layerRes := ""

		for len(nodeList) > 0 {
			node := nodeList[0]
			nodeList = nodeList[1:]

			if node.IsEmpty() {
				layerRes += "null "
			} else {
				layerRes += fmt.Sprintf("%d ", node.Key())
				nextNodeList = append(nextNodeList, node.Left(), node.Right())
			}
		}

		layerRes += "\n"

		if len(nextNodeList) < 1 {
			break
		}

		res += layerRes
		nodeList = nextNodeList
	}

	return res
}

func isBST(node IBSTNode) bool {
	values := dfsBSTNode(node)
	backup := support.DeepCopyIntSlice(values)

	sort.Ints(values)

	return support.CompareOnIntSlice(backup, values)
}

func dfsBSTNode(node IBSTNode) []int {
	if node.IsEmpty() {
		return nil
	} else if node.Left() == nil && node.Right() == nil {
		return []int{node.Key()}
	}

	leftValues := dfsBSTNode(node.Left())

	rightValues := dfsBSTNode(node.Right())

	res := make([]int, 0, len(leftValues)+len(rightValues)+1)

	res = append(res, leftValues...)
	res = append(res, node.Key())
	res = append(res, rightValues...)

	return res
}
