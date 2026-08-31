package mlog

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
)

var mlogPkgPrefix = ""

func init() {
	pc := make([]uintptr, 1)
	runtime.Callers(1, pc) // 0=runtime.Callers, 1=init

	fn := runtime.FuncForPC(pc[0]).Name() // 'xxx/log.init.0'，'xxx'是基于go module name的相对路径
	if fn != "" {
		lastIndex := strings.LastIndex(fn, ".")
		if lastIndex >= 0 {
			index := strings.LastIndex(fn[:lastIndex], ".")
			if index >= 0 {
				mlogPkgPrefix = fn[:index+1] // include '.'
				return
			}
		}
	}

	mlogPkgPrefix = "mlog."
}

func codePosition(buf *bytes.Buffer) {
	var pcs [32]uintptr
	n := runtime.Callers(1, pcs[:])
	if n == 0 {
		buf.WriteString("unknown:0")
		return
	}

	fs := runtime.CallersFrames(pcs[:n])
	for {
		f, more := fs.Next()

		if !isChainFrame(f.Function) {
			buf.WriteString(trimPath(f.File))
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(f.Line))
			return
		}

		if !more {
			break
		}
	}

	buf.WriteString("unknown:0") // 栈上只有链路内部帧（正常情况不会发生）
}

func isChainFrame(function string) bool {
	switch {
	case strings.HasPrefix(function, "log/slog."):
		return true // slog
	case !strings.HasPrefix(function, mlogPkgPrefix):
		return false // mlog
	}

	rest := function[len(mlogPkgPrefix):]
	return !(strings.HasPrefix(rest, "Test") ||
		strings.HasPrefix(rest, "Benchmark") ||
		strings.HasPrefix(rest, "Example") ||
		strings.HasPrefix(rest, "Fuzz"))
}

// 将文件绝对路径裁减到最后两级
func trimPath(fileName string) string {
	lastIndex := strings.LastIndex(fileName, "/")
	if lastIndex >= 0 {
		index := strings.LastIndex(fileName[:lastIndex], "/")
		if index >= 0 {
			return fileName[index+1:]
		}
	}

	return fileName
}
