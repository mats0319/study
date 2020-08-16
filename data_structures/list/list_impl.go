package list

import "fmt"

type ListImpl struct {
	Value int
	Next *ListImpl
}

var _ List = (*ListImpl)(nil)

func MakeEmpty() List {
	return &ListImpl{}
}

func (l *ListImpl) PrintList() {
	p := l

	listStr := ""
	for ; p.Next != nil; p = p.Next {
		listStr += fmt.Sprintf("%d -> ", p.Value)
	}
	listStr += fmt.Sprintf("%d .", p.Value)

	fmt.Println("List:", listStr)
}

func (l *ListImpl) Find(key int) int {
	p := l

	index := -1
	for ; p != nil; p = p.Next {
		index++
		if key == p.Value {
			break
		}
	}

	return index
}

func (l *ListImpl) Insert(key int) List {
	// todo: test is necessary?
	//if l.IsNil() {
	//	return &ListImpl{Value: key}
	//}

	p := l
	for ; p != nil; p = p.Next {
	}
	p = &ListImpl{Value: key}

	return l
}

// todo: review code, optimize
func (l *ListImpl) Delete(key int) (res List, ok bool) {
	if l.IsNil() { // nil
		return MakeEmpty(), false
	}

	if l.Next == nil { // one node
		if l.Value == key {
			return MakeEmpty(), true
		} else if l.Value != key {
			return l, false
		}
	}


	pre, curr := &ListImpl{Next: l}, l
	for ; curr != nil; pre, curr = pre.Next, curr.Next {
		if curr.Value == key {
			pre.Next = curr.Next
			ok = true
			break
		}
	}

	return l, ok
}

func (l *ListImpl) IsEmpty() bool {
	return l.IsNil() || (l.Value == 0 && l.Next == nil)
}

func (l *ListImpl) IsNil() bool {
	return l == nil
}
