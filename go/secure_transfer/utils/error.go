package utils

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Code    int
	Message string

	Cause  error
	Params map[string]any
	Stack  []uintptr
}

var _ error = (*Error)(nil)

func NewError(code int, message string) *Error {
	var stack [32]uintptr
	n := runtime.Callers(2, stack[:]) // skip 'runtime.caller' and 'NewError'

	return &Error{
		Code:    code,
		Message: message,
		Params:  make(map[string]any),
		Stack:   stack[:n],
	}
}

// Error print simple error message, without params
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("error code: %d", e.Code)
}

// String print all details, use in server log
func (e *Error) String() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("error code: %d, message: %s\nerr: %v\nparams: %#v\nstack trace: \n%s\n",
		e.Code, e.Message, e.Cause, e.Params, e.stackTrace())
}

func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

func (e *Error) WithParam(key string, value any) *Error {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}

	e.Params[key] = value
	return e
}

func (e *Error) stackTrace() string {
	var builder strings.Builder
	frames := runtime.CallersFrames(e.Stack)

	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}

	return builder.String()
}
