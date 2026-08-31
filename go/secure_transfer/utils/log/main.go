package mlog

import (
	"context"
	"log/slog"
)

type writerFlag int

const (
	W_Stdout writerFlag = 1 << iota
	W_File
)

var mLogger *slog.Logger

func Initialize(wf writerFlag) {
	l, err := NewLogger(wf, slog.LevelDebug)
	if err != nil {
		panic(err)
	}

	mLogger = l
	slog.SetDefault(mLogger)

	Info("> Log init.")
}

func Close() {
	if mLogger != nil {
		CloseLogger(mLogger)
	}
}

func NewLogger(wf writerFlag, level slog.Level) (*slog.Logger, error) {
	h, err := newHandler(wf, level)
	if err != nil {
		return nil, err
	}

	return slog.New(h), nil
}

func CloseLogger(l *slog.Logger) {
	if l == nil {
		return
	}

	h, ok := l.Handler().(*handler)
	if ok {
		h.close()
	}
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Log 适用于需要通过'WithAttrs'/'WithGroup'定制输出内容、生成新的logger实例的场景
// 例如：http请求、db操作......
//
// 其实可以把slog default logger / slog level复制到这里，在使用过程中就不需要使用slog
// 日志中的代码位置由 codePosition 跳过本包与 slog 的栈帧自动定位，
// 经 mlog.Info()/mlog.Log() 或直接 slog 调用均指向真实调用点
func Log(logger *slog.Logger, level slog.Level, msg string, fields ...any) {
	if logger == nil {
		logger = slog.Default() // Initialize 后即 mLogger；未初始化时用标准库默认，不 panic
	}

	logger.Log(context.TODO(), level, msg, fields...)
}
