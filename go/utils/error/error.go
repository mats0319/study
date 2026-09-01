package merr

import (
	"fmt"
	"maps"
	"runtime"
	"strings"
)

// Error 结构体类型是导出的，这是考虑到外部可能需要定义该类型的变量（例如用作函数返回值）或类型断言，
// 不建议直接使用结构体创建变量，建议通过本包内提供的'Newxxx'函数创建变量
type Error struct {
	Code    int
	message string

	cause    error
	params   map[string]any
	stack    []uintptr
	stackStr string
}

var _ error = (*Error)(nil)

// NewErrorFactory 工厂模式，每次调用返回不同实例，不支持多goroutine共享
//
// 用法参考测试文件，定义一组函数类型变量，使用时创建新的实例并附加error、参数等信息
// 具体的定义规则应可快速定位出错位置、可快速知道基本信息，举个例子：
// - 'code:10203'，'1'表示用户系统，'02'表示登录模块，'03'表示这是登录模块的第三个错误
func NewErrorFactory(code int, detail string) func() *Error {
	return func() *Error { return newError(code, detail) }
}

//go:noinline
func newError(code int, message string) *Error {
	var stacks [16]uintptr
	n := runtime.Callers(3, stacks[:]) // 0:runtime.callers 1:newError 2:NewError.func1

	return &Error{
		Code:    code,
		message: message,
		stack:   append([]uintptr(nil), stacks[:n]...), // callers默认会让stacks逃逸到堆上
	}
}

// Error 只打印简单的信息，用于外部系统（例如从后端返回给前端）；同时用于实现error接口
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("error code: %d, message: %s\nerr: %v\n", e.Code, e.message, e.cause)
}

// String print all details, use in server log
func (e *Error) String() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("error code: %d, message: %s\nerr: %v\nparams: %#v\nstack trace: \n%s\n",
		e.Code, e.message, e.cause, e.params, e.stackTrace())
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	t, ok := target.(*Error)
	if !ok || t == nil {
		return false
	}

	return e.Code == t.Code
}

func (e *Error) WithCause(err error) *Error {
	if e == nil || err == nil {
		return e
	}

	e.cause = err
	return e
}

func (e *Error) WithParams(params map[string]any) *Error {
	if e == nil || params == nil {
		return e
	}

	if e.params == nil {
		e.params = make(map[string]any)
	}

	maps.Copy(e.params, params)

	return e
}

func (e *Error) WithParam(key string, value any) *Error {
	if e == nil {
		return nil
	}

	if e.params == nil {
		e.params = make(map[string]any)
	}

	e.params[key] = value

	return e
}

func (e *Error) stackTrace() string {
	if e.stackStr != "" || len(e.stack) == 0 { // 已经格式化过，或者没有数据
		return e.stackStr
	}

	var builder strings.Builder
	frames := runtime.CallersFrames(e.stack)

	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}

	e.stackStr = builder.String()

	return e.stackStr
}
