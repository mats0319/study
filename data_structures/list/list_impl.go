package mario

import (
	"fmt"
)

type listImpl struct {
	value NodeValue
	next  *listImpl
}

var _ List = (*listImpl)(nil)

func NilList() List {
	return (*listImpl)(nil)
}

func NewList() List {
	return &listImpl{}
}

func MakeList(value ...NodeValue) List {
	res := &listImpl{}
	p := res
	for i := range value {
		p.next = &listImpl{value: value[i]}
		p = p.next
	}

	return res.next
}

func (l *listImpl) Print() string {
	p := l

	listStr := ""
	for ; p != nil; p = p.next {
		listStr += fmt.Sprintf("%v -> ", p.value)
	}

	if len(listStr) > 0 {
		listStr = listStr[:len(listStr)-3]
	}

	return fmt.Sprintf("List: %s.", listStr)
}

func (l *listImpl) Find(key NodeValue) int {
	p := l

	index := -1
	count := 0
	for ; p != nil; p = p.next {
		if key == p.value {
			index = count
			break
		}

		count++
	}

	return index
}

func (l *listImpl) Insert(key NodeValue) List {
	h := &listImpl{next: l}
	p := h
	for ; p.next != nil; p = p.next {
	}
	p.next = &listImpl{value: key}

	return h.next // h define for situation: l is nil
}

func (l *listImpl) Delete(key NodeValue) (res List, ok bool) {
	h := &listImpl{next: l}
	pre, curr := h, l
	for ; curr != nil; pre, curr = pre.next, curr.next {
		if curr.value == key {
			pre.next = curr.next
			ok = true
			break
		}
	}

	res = h.next
	if res.IsNil() {
		res = NewList()
	}

	return
}

func (l *listImpl) IsEqual(li *listImpl) bool {
	a, b := l, li
	for ; a != nil && b != nil && a.value == b.value; a, b = a.next, b.next {
	}

	return a == nil && b == nil
}

func (l *listImpl) IsEmpty() bool {
	return l.IsNil() || (l.value == ZeroValue && l.next == nil)
}

func (l *listImpl) IsNil() bool {
	return l == nil
}

func (l *listImpl) Next() List {
	if l.IsNil() {
		return nil
	}

	return l.next
}
