package listp_test

import (
	listp "github.com/mats9693/study/data_structures/list/list_pointer"
	"testing"
)

func TestNode_Insert(t *testing.T) {
	list := listp.CreateList(1)
	list = list.Insert(2, 1)

	l := listp.MakeList(2, 1)
	if !listp.IsEqual(list, l) {
		t.Errorf("表插入失败\n")
	}
}

func TestNode_Find(t *testing.T) {
	list := listp.MakeList(1, 2, 3, 5)

	if list.Find(3) != 2 {
		t.Errorf("表查找失败\n")
	}
}

func TestNode_Delete(t *testing.T) {
	list := listp.MakeList(1, 2, 3, 5)

	if list.Delete(3) != 2 {
		t.Errorf("表删除失败\n")
	}
}
