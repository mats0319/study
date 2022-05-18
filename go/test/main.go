package main

import (
	"github.com/mats9693/utils/support"
	"log"
)

func main() {
	data := support.GenerateRandomIntSlice(20, 100, 300)

	hashMapIns := newHashMap(data)

	ok := hashMapIns.find(300) // true
	if !ok {
		log.Println("can't find exist item")
		return
	}
	ok = hashMapIns.find(301) // false
	if ok {
		log.Println("find not exist item")
		return
	}

	hashMapIns.modify(300, 400)
	ok = hashMapIns.find(300) // false
	if ok { // 'find' method has tested above
		log.Printf("find not exist item del by 'modify'")
		return
	}
	ok = hashMapIns.find(400) // true
	if !ok { // 'find' method has tested above
		log.Printf("can't find exist item insert by 'modify'")
		return
	}

	hashMapIns.delete(300)
	ok = hashMapIns.find(300) // false
	if ok {
		log.Println("delete failed")
		return
	}

	log.Println("pass test")
}

type hashTable struct {
	data []*listNode
}

type listNode struct {
	exist bool

	child []*listNode
}

func newHashMap(data []int) *hashTable {
	ins := &hashTable{
		data: make([]*listNode, 10),
	}

	for i := range data {
		ins.insert(uint(data[i]))
	}

	return ins
}

func (i *hashTable) insert(key uint) {
	p := i.data
	for key >= 10 {
		index := key % 10
		key /= 10

		if p[index] == nil {
			p[index] = &listNode{
				child: make([]*listNode, 10),
			}
		}

		p = p[index].child
	}

	if p[key] != nil { // 'key' exist
		return
	}

	p[key] = &listNode{
		exist: true,
		child: make([]*listNode, 10),
	}
}

func (i *hashTable) find(key uint) bool {
	node := i.getExistNode(key)

	return node != nil && node.exist
}

func (i *hashTable) modify(oldKey uint, newKey uint) {
	i.delete(oldKey)
	i.insert(newKey)

	return
}

func (i *hashTable) delete(key uint) {
	node := i.getExistNode(key)

	if node != nil {
		node.exist = false
	}
}

func (i *hashTable) getExistNode(key uint) *listNode {
	p := i.data
	isExist := true
	for key >= 10 {
		index := key % 10
		key /= 10

		if p[index] == nil { // not exist
			isExist = false
			break
		}

		p = p[index].child
	}

	if !isExist || p[key] == nil || !p[key].exist { // not exist
		return nil
	}

	return p[key]
}
