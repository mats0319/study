package mario

import (
	"fmt"
)

type listImpl struct {
	Value int
	Next *listImpl
}

var _ List = (*listImpl)(nil)

func NilList() List {
	return (*listImpl)(nil)
}

func NewList() List {
	return &listImpl{}
}

func MakeList(value ...int) List {
	res := &listImpl{}
	p := res
	for i := range value {
		p.Next = &listImpl{Value: value[i]}
		p = p.Next
	}

	return res.Next
}

func (l *listImpl) Print() string {
	p := l

	listStr := ""
	for ; p != nil; p = p.Next {
		listStr += fmt.Sprintf("%d -> ", p.Value)
	}

	if len(listStr) > 0 {
		listStr = listStr[:len(listStr)-3]
	}

	return fmt.Sprintf("List: %s.", listStr)
}

func (l *listImpl) Find(key int) int {
	p := l

	index := -1
	count := 0
	for ; p != nil; p = p.Next {
		if key == p.Value {
			index = count
			break
		}

		count++
	}

	return index
}

func (l *listImpl) Insert(key int) List {
	h := &listImpl{Next: l}
	p := h
	for ; p.Next != nil; p = p.Next {
	}
	p.Next = &listImpl{Value: key}

	return h.Next // h define for situation: l is nil
}

func (l *listImpl) Delete(key int) (res List, ok bool) {
	h := &listImpl{Next: l}
	pre, curr := h, l
	for ; curr != nil; pre, curr = pre.Next, curr.Next {
		if curr.Value == key {
			pre.Next = curr.Next
			ok = true
			break
		}
	}

	res = h.Next
	if res.IsNil() {
		res = NewList()
	}

	return
}

func (l *listImpl) IsEqual(li *listImpl) bool {
	a, b := l, li
	for ; a != nil && b != nil && a.Value == b.Value; a, b = a.Next, b.Next{
	}

	return a == nil && b == nil
}

func (l *listImpl) IsEmpty() bool {
	return l.IsNil() || (l.Value == 0 && l.Next == nil)
}

func (l *listImpl) IsNil() bool {
	return l == nil
}
