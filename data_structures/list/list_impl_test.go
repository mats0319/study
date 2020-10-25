package mario

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestListImpl_Print(t *testing.T) {
	convey.Convey("test print nil List", t, func() {
		nilList := NilList()
		listStr := nilList.Print()

		convey.So(listStr, convey.ShouldEqual, "List: .")
	})

	convey.Convey("test print an empty List", t, func() {
		emptyList := NewList()
		listStr := emptyList.Print()

		convey.So(listStr, convey.ShouldEqual, "List: 0 .")
	})

	convey.Convey("test print a normal list", t, func() {
		list := MakeList(1, 2, 3)
		listStr := list.Print()

		convey.So(listStr, convey.ShouldEqual, "List: 1 -> 2 -> 3 .")
	})
}

func TestListImpl_Find(t *testing.T) {
	convey.Convey("test find in nil List", t, func() {
		nilList := NilList()
		res := nilList.Find(10)

		convey.So(res, convey.ShouldEqual, -1)
	})

	convey.Convey("test find in an empty List", t, func() {
		emptyList := NewList()
		res := emptyList.Find(10)

		convey.So(res, convey.ShouldEqual, -1)
	})

	convey.Convey("test find an item not exist", t, func() {
		list := MakeList(1, 2, 3)
		res := list.Find(10)

		convey.So(res, convey.ShouldEqual, -1)
	})

	convey.Convey("test find an item at start", t, func() {
		list := MakeList(1, 2, 3)
		res := list.Find(1)

		convey.So(res, convey.ShouldEqual, 0)
	})

	convey.Convey("test find an item at middle", t, func() {
		list := MakeList(1, 2, 3)
		res := list.Find(2)

		convey.So(res, convey.ShouldEqual, 1)
	})

	convey.Convey("test find an item at end", t, func() {
		list := MakeList(1, 2, 3)
		res := list.Find(3)

		convey.So(res, convey.ShouldEqual, 2)
	})
}

func TestListImpl_Insert(t *testing.T) {
	convey.Convey("test insert to nil List", t, func() {
		nilList := NilList()
		list := nilList.Insert(1)

		// nil List can not save values, reference to method definition in List interface
		convey.So(nilList.Print(), convey.ShouldEqual, "List: .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 1 .")
	})

	convey.Convey("test insert to an empty List", t, func() {
		emptyList := NewList()
		list := emptyList.Insert(1)

		convey.So(emptyList.Print(), convey.ShouldEqual, "List: 0 -> 1 .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 0 -> 1 .")
	})

	convey.Convey("test insert to a normal List", t, func() {
		emptyList := MakeList(1)
		list := emptyList.Insert(1)

		convey.So(emptyList.Print(), convey.ShouldEqual, "List: 1 -> 1 .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 1 -> 1 .")
	})
}

func TestListImpl_Delete(t *testing.T) {
	convey.Convey("test delete a nil List", t, func() {
		nilList := NilList()
		list, ok := nilList.Delete(10)

		convey.So(nilList.Print(), convey.ShouldEqual, "List: .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 0 .")
		convey.So(ok, convey.ShouldEqual, false)
	})

	convey.Convey("test delete an empty List", t, func() {
		emptyList := NewList()
		list, ok := emptyList.Delete(10)

		convey.So(emptyList.Print(), convey.ShouldEqual, "List: 0 .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 0 .")
		convey.So(ok, convey.ShouldEqual, false)
	})

	convey.Convey("test delete an empty List, delete zero value", t, func() {
		emptyList := NewList()
		list, ok := emptyList.Delete(0) // zero value

		convey.So(emptyList.Print(), convey.ShouldEqual, "List: 0 .")
		convey.So(list.Print(), convey.ShouldEqual, "List: 0 .")
		convey.So(ok, convey.ShouldEqual, true)
	})

	convey.Convey("test delete an item not exist", t, func() {
		list := MakeList(1, 2, 3)
		res, ok := list.Delete(10)

		convey.So(list.Print(), convey.ShouldEqual, "List: 1 -> 2 -> 3 .")
		convey.So(res.Print(), convey.ShouldEqual, "List: 1 -> 2 -> 3 .")
		convey.So(ok, convey.ShouldEqual, false)
	})

	convey.Convey("test delete at start", t, func() {
		list := MakeList(1, 2, 3)
		res, ok := list.Delete(1)

		// left this bug in code, reference to definition in List interface
		convey.So(list.Print(), convey.ShouldEqual, "List: 1 -> 2 -> 3 .")
		convey.So(res.Print(), convey.ShouldEqual, "List: 2 -> 3 .")
		convey.So(ok, convey.ShouldEqual, true)
	})

	convey.Convey("test delete at middle", t, func() {
		list := MakeList(1, 2, 3)
		res, ok := list.Delete(2)

		convey.So(list.Print(), convey.ShouldEqual, "List: 1 -> 3 .")
		convey.So(res.Print(), convey.ShouldEqual, "List: 1 -> 3 .")
		convey.So(ok, convey.ShouldEqual, true)
	})

	convey.Convey("test delete at start", t, func() {
		list := MakeList(1, 2, 3)
		res, ok := list.Delete(3)

		convey.So(list.Print(), convey.ShouldEqual, "List: 1 -> 2 .")
		convey.So(res.Print(), convey.ShouldEqual, "List: 1 -> 2 .")
		convey.So(ok, convey.ShouldEqual, true)
	})
}

// enum situation of one List:
// nil / empty / normal
// after permutation: 9 cases, add one: normal - normal, not equal, total 10 cases.
func TestListImpl_IsEqual(t *testing.T) {
	convey.Convey("test ListA(nil) is equal to ListB(nil)", t, func() {
		var la *listImpl = (*listImpl)(nil)
		var lb *listImpl = (*listImpl)(nil)
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, true)
	})

	convey.Convey("test ListA(nil) is equal to ListB(empty)", t, func() {
		var la *listImpl = (*listImpl)(nil)
		var lb *listImpl = &listImpl{}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(nil) is equal to ListB(normal)", t, func() {
		var la *listImpl = (*listImpl)(nil)
		var lb *listImpl = &listImpl{value: 1}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(empty) is equal to ListB(nil)", t, func() {
		var la *listImpl = &listImpl{}
		var lb *listImpl = (*listImpl)(nil)
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(empty) is equal to ListB(empty)", t, func() {
		var la *listImpl = &listImpl{}
		var lb *listImpl = &listImpl{}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, true)
	})

	convey.Convey("test ListA(empty) is equal to ListB(normal)", t, func() {
		var la *listImpl = &listImpl{}
		var lb *listImpl = &listImpl{value: 1}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(normal) is equal to ListB(nil)", t, func() {
		var la *listImpl = &listImpl{value: 1}
		var lb *listImpl = (*listImpl)(nil)
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(normal) is equal to ListB(empty)", t, func() {
		var la *listImpl = &listImpl{value: 1}
		var lb *listImpl = &listImpl{}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})

	convey.Convey("test ListA(normal) is equal to ListB(normal), equal", t, func() {
		var la *listImpl = &listImpl{value: 1}
		var lb *listImpl = &listImpl{value: 1}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, true)
	})

	convey.Convey("test ListA(normal) is equal to ListB(normal), not equal", t, func() {
		var la *listImpl = &listImpl{value: 1}
		var lb *listImpl = &listImpl{value: 2}
		res := la.IsEqual(lb)

		convey.So(res, convey.ShouldEqual, false)
	})
}
