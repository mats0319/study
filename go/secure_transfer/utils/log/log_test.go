package mlog

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestCustomLog(t *testing.T) {
	Initialize()
	defer Close()

	logger := DefaultLogger().WithGroup("group1").WithGroup("group2")

	Log(logger, slog.LevelDebug, "log msg",
		slog.Int("key1", 10),
		slog.String("key2", "value2"),
	)

	// 结论：双写测试成功，可以根据需要写入日志文件或写入默认输出
}

func TestLogLevel(t *testing.T) {
	Initialize()
	defer Close()

	Debug("debug level log", slog.Any("error", errors.New("debug error")))
	Info("info level log")
	Warn("warn level log")
	Error("error level log", slog.Any("error", errors.New("error")))
}

func TestLogSplitFile(t *testing.T) {
	Initialize()
	defer Close()

	lastSize := handler.Size
	currentSize := handler.Size
	logger := DefaultLogger().WithGroup("groupName")

	for lastSize <= currentSize { // exit loop when emit log split
		lastSize = currentSize

		Log(logger, 100, "this is a long long test log message",
			slog.String("key1", "value1"),
			slog.String("key2", "value2"),
		)

		currentSize = handler.Size
	}

	time.Sleep(time.Second * 3) // 阻塞，避免异步压缩goroutine因主程序退出而停止执行
}
