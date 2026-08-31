package merr

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var (
	errTestA = NewErrorFactory(10001, "error A")
	errTestB = NewErrorFactory(10002, "error B")
)

func Example() {
	e := errTestA().WithCause(errors.New("err msg")).
		WithParam("first param", "first value").
		WithParam("second param", 10)

	fmt.Println(e.Error())

	// Output:
	// error code: 10001, message: error A
	// err: err msg
}

// ---- 字段与构造 ----

func TestFactoryReturnsFreshInstances(t *testing.T) {
	e1 := errTestA().WithParam("k", "v1")
	e2 := errTestA().WithParam("k", "v2")

	if e1 == e2 {
		t.Fatal("factory should return a new instance per call")
	}

	// 新实例不应继承前一个实例的参数
	if strings.Contains(e1.String(), "v2") || strings.Contains(e2.String(), "v1") {
		t.Error("instances should not share params")
	}
}

// ---- errors.Is 相关 ----

func TestIsSameCode(t *testing.T) {
	// 工厂每次返回新实例，但错误码相同应视为同一类错误
	if !errors.Is(errTestA(), errTestA()) {
		t.Error("errors.Is should be true for two instances with the same code")
	}
}

func TestIsDifferentCode(t *testing.T) {
	if errors.Is(errTestA(), errTestB()) {
		t.Error("errors.Is should be false for different codes")
	}
}

func TestIsViaCauseChain(t *testing.T) {
	outer := errTestA().WithCause(errTestB())

	if !errors.Is(outer, errTestB()) {
		t.Error("errors.Is should match the wrapped cause through Unwrap")
	}
	if !errors.Is(outer, errTestA()) {
		t.Error("errors.Is should match the outer error itself")
	}
}

func TestIsNonErrTarget(t *testing.T) {
	if errors.Is(errTestA(), errors.New("plain error")) {
		t.Error("errors.Is should be false for a non-err target")
	}
}

func TestIsNilTarget(t *testing.T) {
	if errors.Is(errTestA(), nil) {
		t.Error("errors.Is should be false when err is non-nil and target is nil")
	}
	if !errors.Is(nil, nil) {
		t.Error("errors.Is(nil, nil) should be true")
	}
}

func TestIsDirectCall(t *testing.T) {
	e := errTestA()

	if !e.Is(errTestA()) {
		t.Error("direct Is call should match same code")
	}
	if e.Is(errTestB()) {
		t.Error("direct Is call should not match different code")
	}
	if e.Is(nil) || e.Is(errors.New("plain")) {
		t.Error("direct Is call should be false for nil or non-merr target")
	}
}

// ---- Unwrap 相关 ----

func TestUnwrap(t *testing.T) {
	cause := errors.New("root cause")
	e := errTestA().WithCause(cause)

	if e.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
	if errTestA().Unwrap() != nil {
		t.Error("Unwrap should return nil when there is no cause")
	}
	// errors.Is 应能沿 Unwrap 链到达普通 error
	if !errors.Is(e, cause) {
		t.Error("errors.Is should traverse the Unwrap chain to reach a plain cause")
	}
}

// ---- 参数 ----

func TestWithParam(t *testing.T) {
	e := errTestA().WithParam("k", "v")

	if !strings.Contains(e.String(), "k") || !strings.Contains(e.String(), "v") {
		t.Errorf("String() should contain the param, got: %s", e.String())
	}

	// Error() 用于对外输出，不应包含 params
	if strings.Contains(e.Error(), "k") {
		t.Error("Error() should not contain params")
	}
}

func TestWithParamsCopiesInput(t *testing.T) {
	m := map[string]any{"k1": "v1"}
	e := errTestA().WithParams(m)

	// 调用方后续修改原 map，不应影响错误实例
	m["k1"] = "CHANGED"
	m["k2"] = "ADDED"

	s := e.String()
	if strings.Contains(s, "CHANGED") || strings.Contains(s, "ADDED") {
		t.Errorf("WithParams should copy the input map, got: %s", s)
	}
}

func TestWithParamsMergesExisting(t *testing.T) {
	e := errTestA().WithParam("k1", "v1").WithParams(map[string]any{"k2": "v2"})

	s := e.String()
	if !strings.Contains(s, "k1") || !strings.Contains(s, "k2") {
		t.Errorf("WithParams should merge into existing params, got: %s", s)
	}
}

func TestWithParamsNilMap(t *testing.T) {
	e := errTestA().WithParams(nil)

	// 不应 panic
	if !strings.Contains(e.String(), "params:") {
		t.Error("WithParams(nil) should be a no-op without panic")
	}
}

// ---- 输出格式 ----

func TestErrorFormat(t *testing.T) {
	e := errTestA().WithCause(errors.New("boom"))
	want := "error code: 10001, message: error A\nerr: boom\n"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestStringContainsStackAndParams(t *testing.T) {
	e := errTestA().WithParam("k", "v")
	s := e.String()

	for _, want := range []string{
		"error code: 10001",
		"message: error A",
		"params:",
		"stack trace:",
		"TestStringContainsStackAndParams",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("String() should contain %q, got: %s", want, s)
		}
	}

	// 栈文本已缓存，两次输出应一致
	if e.String() != s {
		t.Error("String() should be stable across calls (stack cached)")
	}
}

func TestEmptyStackString(t *testing.T) {
	// 直接构造（绕过工厂），stack 为空
	e := &Error{Code: 1, message: "manual"}

	if got := e.stackTrace(); got != "" {
		t.Errorf("stackTrace() with empty stack = %q, want empty", got)
	}
	if !strings.Contains(e.String(), "stack trace:") {
		t.Errorf("String() should not panic with empty stack, got: %s", e.String())
	}
}

// ---- nil receiver ----

func TestNilReceiverReads(t *testing.T) {
	var e *Error

	if e.Error() != "" {
		t.Error("nil *Error.Error() should be empty")
	}
	if e.String() != "" {
		t.Error("nil *Error.String() should be empty")
	}
	if e.Unwrap() != nil {
		t.Error("nil *Error.Unwrap() should be nil")
	}

	// Is 不应 panic
	if e.Is(errTestA()) {
		t.Error("nil *Error should not match a non-nil target")
	}
	if !e.Is(nil) {
		t.Error("nil *Error should match nil target")
	}
}

func TestNilReceiverMutators(t *testing.T) {
	var e *Error

	if got := e.WithCause(errors.New("x")); got != nil {
		t.Error("nil *Error.WithCause should return nil")
	}
	if got := e.WithParam("k", "v"); got != nil {
		t.Error("nil *Error.WithParam should return nil")
	}
	if got := e.WithParams(map[string]any{"k": "v"}); got != nil {
		t.Error("nil *Error.WithParams should return nil")
	}
}
