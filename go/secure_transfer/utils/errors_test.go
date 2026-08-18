package utils

import (
	"errors"
	"fmt"
	"testing"
)

var ErrForTest = newError(-1, "test error string")

func TestLogStyle(t *testing.T) {
	e := ErrForTest().WithCause(errors.New("a new error")).
		WithParam("first param", "first value").
		WithParam("second param", 10000)

	t.Log(e.String())
}

func TestMultiUseOnOneInstance(t *testing.T) {
	e := ErrForTest().WithParam("key1", "value1")
	e = ErrForTest().WithParam("key2", "value2")

	t.Log(e.String())

	// 如果Err直接是*Error类型的实例，那么在多处使用该实例时，会继承历史数据
	// （即：假设数据库未启动，dbError一直触发，并且每次都会使用不同参数调用withParam()，
	// 那么一段时间以后，Err的map可能就会带满所有程序中写到的key；
	// 另外，多次触发同一位置的Err没有这个问题，因为error和相同的key会覆盖历史数据）
	//
	// 我们可以看到，实际场景中确实存在同一个Err多处使用的情况，所以我们将Err统一编辑成函数，
	// 这样在每次调用Err时，都会产生新的实例，不会存在继承历史数据的情况。
}

func TestNilError(t *testing.T) {
	fmt.Println(f1() == nil) // false
	fmt.Println(f2() == nil)

	// 结论：需要注意，任何自定义类型实例赋值给接口类型，接口 != nil，
	// 如果接口类型想要正常判空，则必须在正确流程中明确返回'nil'，而非值为nil的自定义类型或者其他
}

func f1() error {
	var err *Error = nil // '= nil' is not necessary
	return err
}

func f2() error {
	return nil
}
