package tree

import (
	"github.com/mats9693/utils/support"
	"testing"
)

func TestNewBST(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
		data := support.GenerateRandomIntSlice(1000, 1000)

		bstIns := newBST(data)

		if !isBST(bstIns.root) {
			t.Logf("data: %v\ntree: \n%s", data, printBST(bstIns.root))
			t.Fail()
		}
	})
}

func TestBinarySearchTree_Find(t *testing.T) {
	test100000Times(t, func(t *testing.T) {
		data := support.GenerateRandomIntSlice(20, 100, 0, 10, 101)

		bstIns := newBST(data)

		values := []int{0, 10, 101, 200}
		expected := []bool{true, true, true, false}

		if len(values) != len(expected) {
			t.Logf("unmatched amount of values and expected results")
			t.Fail()
		}

		for i := range values {
			_, ok := bstIns.Find(values[i])
			if ok != expected[i] {
				t.Logf("index: %d, find: %d\ntree: \n%s\n\twant: %t\n\tget: %t\n",
					i, values[i], printBST(bstIns.root), expected[i], ok)
				t.Fail()
			}
		}

		if !isBST(bstIns.root) {
			t.Logf("data: %v\ntree: \n%s", data, printBST(bstIns.root))
			t.Fail()
		}
	})
}
