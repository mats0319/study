package mlog

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestCustomLog(t *testing.T) {
	Initialize(W_File)
	defer Close()

	h, e := NewHandler(W_Stdout)
	if e != nil {
		t.Fatal(e)
	}
	logger := slog.New(h).WithGroup("group1").WithGroup("group2")

	{
		Debug("log msg", slog.String("key", "value"))
		Log(logger, slog.LevelDebug, "log msg",
			slog.Int("key1", 10),
			slog.String("key2", "value2"),
		)
	}

	// 结论：双写测试成功，可以根据需要写入日志文件或默认输出，可以为文件/命令行写入不同内容
}

func TestLogLevel(t *testing.T) {
	Initialize(W_Stdout)
	defer Close()

	Debug("debug level log", slog.Any("error", errors.New("debug error")))
	Info("info level log")
	Warn("warn level log")
	Error("error level log", slog.Any("error", errors.New("error")))
}

func TestLogSplitFile(t *testing.T) {
	Initialize(W_File)
	defer Close()

	lastSize := handler.Size
	currentSize := handler.Size

	for lastSize <= currentSize { // exit loop when emit log split
		lastSize = currentSize

		Log(nil, 100, "this is a long long test log message",
			slog.String("key1", "value1"),
			slog.String("key2", "value2"),
		)

		currentSize = handler.Size
	}

	time.Sleep(time.Second * 3) // 阻塞，避免异步压缩goroutine因主程序退出而停止执行
}
