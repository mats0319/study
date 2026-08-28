package mlog

import (
	"context"
	"log/slog"
)

type writeFlag int

const (
	W_Stdout writeFlag = 1 << iota
	W_File
)

var h = &handler{}

func Initialize(wf writeFlag) {
	h = newHandler(wf, slog.LevelDebug)

	slog.SetDefault(slog.New(h))

	Info("> Log init.")
}

func Close() {
	h.close()
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
// 已调整为与默认打印函数使用相同的调用层级数（日志中的代码位置一项，使用'mlog.Info()'/'mlog.Log()'均显示为该函数位置）
func Log(logger *slog.Logger, level slog.Level, msg string, fields ...any) {
	if logger == nil {
		logger = slog.Default()
	}

	logger.Log(context.TODO(), level, msg, fields...)
}
