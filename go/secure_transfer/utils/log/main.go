package mlog

import (
	"log"
	"log/slog"
)

var handler *Handler

func DefaultLogger() *slog.Logger {
	h, err := defaultHandler()
	if err != nil {
		panic(err)
	}

	return slog.New(h)
}

func Initialize() {
	var err error
	handler, err = newHandler("log.log", 1)
	if err != nil {
		log.Fatalln("open log file failed, error:", err)
	}

	slog.SetDefault(slog.New(handler))

	Info("> Config init.") // 不是这里才初始化的，但是只有这里（日志）初始化之后才能使用自定义结构打印这句话
	Info("> Log init.")
}

func Close() {
	handler.close()
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

	switch level {
	case slog.LevelInfo:
		logger.Info(msg, fields...)
	case slog.LevelWarn:
		logger.Warn(msg, fields...)
	case slog.LevelError:
		logger.Error(msg, fields...)
	default:
		logger.Debug(msg, fields...)
	}
}
