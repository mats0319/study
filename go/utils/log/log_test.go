package mlog

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// ExampleDefault / ExampleLogger 作为文档示例。
// 注意：日志行含毫秒时间戳，而 example 的 // Output: 是逐字节精确比较
// （不支持正则/通配），无法校验带时间的行，因此不加 Output 注释——
// 无 Output 注释的 example 只编译不执行，行为校验由 Test* 负责。
func ExampleDefault() {
	Initialize(W_Stdout)
	defer Close()

	Info("default use", slog.Int("key", 10))
}

func ExampleLogger() {
	// 不同logger可以将相同/不同的内容写进相同/不同的输出
	logger, err := NewLogger(W_Stdout, slog.LevelDebug)
	if err != nil {
		panic(err)
	}
	defer CloseLogger(logger)

	Log(logger, slog.LevelDebug, "use custom logger")
}

// TestHandlerLineFormat 校验整行日志结构：
// [Time] [Level] [pos] msg | g.k1=v1 g.k2=v2
// 一条日志覆盖：时间/级别格式、源码位置、消息、handler 级 attrs 与 record attrs
// （含 Any）的组前缀、单行输出。
func TestHandlerLineFormat(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h.WithGroup("g1").WithGroup("g2").
		WithAttrs([]slog.Attr{slog.String("hkey", "hval")}))

	l.Info("hello",
		slog.Int("k1", 10),
		slog.String("k2", "value2"),
		slog.Any("err", errors.New("oops")),
	)

	out := buf.String()
	wantLine := regexp.MustCompile(
		`^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\] \[INFO\] \[[^]]+:\d+\] ` +
			`hello \| g1\.g2\.hkey=hval g1\.g2\.k1=10 g1\.g2\.k2=value2 g1\.g2\.err=oops\n$`,
	)
	if !wantLine.MatchString(out) {
		t.Errorf("日志行格式不符合预期:\n%s", out)
	}
	if strings.Count(out, "\n") != 1 {
		t.Errorf("一行日志应只有一个换行符，实际输出: %q", out)
	}
	// 直接 slog 调用时位置也应指向调用点（log_test.go）
	wantContains(t, out, "log_test.go:")
}

func TestLevelFiltering(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelWarn)
	l := slog.New(h)

	l.Debug("debug should be filtered")
	l.Info("info should be filtered")
	l.Warn("warn kept")
	l.Error("error kept")

	out := buf.String()
	wantContains(t, out, "[WARN]")
	wantContains(t, out, "] warn kept")
	wantContains(t, out, "[ERROR]")
	wantContains(t, out, "] error kept")
	wantNotContains(t, out, "should be filtered")
}

func TestLogLevelPreserved(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h)

	Log(l, slog.LevelError, "err msg", slog.String("k", "v"))
	wantContains(t, buf.String(), "[ERROR]")
	wantContains(t, buf.String(), "] err msg | k=v")

	// 自定义级别应原样保留
	buf.Reset()
	customLevel := slog.Level(12)
	Log(l, customLevel, "custom level")
	wantContains(t, buf.String(), "["+customLevel.String()+"]")
	wantContains(t, buf.String(), "] custom level")
}

func TestLogSourcePosition(t *testing.T) {
	// 位置段应指向调用点所在文件 log_test.go，无论经哪种路径调用
	re := regexp.MustCompile(`\[[^]]*log_test\.go:\d+\]`)

	// 1) 直接 slog 调用
	h1, buf := newCaptureHandler(slog.LevelDebug)
	slog.New(h1).Info("direct position check")
	if !re.MatchString(buf.String()) {
		t.Errorf("直接 slog 调用的位置应指向 log_test.go，实际输出:\n%s", buf.String())
	}

	// 2) 经 mlog.Log 包装
	buf.Reset()
	Log(slog.New(h1), slog.LevelInfo, "log position check")
	if !re.MatchString(buf.String()) {
		t.Errorf("Log() 的位置应指向 log_test.go，实际输出:\n%s", buf.String())
	}

	// 3) 经 mlog.Info 包装（走包级 mLogger；Initialize 会新建 mLogger，需重定向其输出）
	buf.Reset()
	Initialize(W_Stdout)
	defer Close()
	if h, ok := mLogger.Handler().(*handler); ok {
		h.handlerWriter.writer = buf
	}

	Info("info position check")
	if !re.MatchString(buf.String()) {
		t.Errorf("Info() 的位置应指向 log_test.go，实际输出:\n%s", buf.String())
	}
}

func TestInitializeClose(t *testing.T) {
	// 未初始化时 Close 不应 panic
	Close()

	Initialize(W_Stdout)
	defer Close()

	if mLogger == nil {
		t.Fatal("Initialize 后 mLogger 不应为 nil")
	}
	if _, ok := mLogger.Handler().(*handler); !ok {
		t.Error("Initialize 后 mLogger 的 handler 应为 mlog.handler")
	}
	// Close 可重复调用且不应 panic（W_Stdout 模式 file 为 nil）
	Close()
	Close()
}

func TestWriteFlag(t *testing.T) {
	t.Run("stdout only", func(t *testing.T) {
		h, err := newHandler(W_Stdout, slog.LevelDebug)
		if err != nil {
			t.Fatal(err)
		}
		if h.handlerWriter.writer != os.Stdout {
			t.Errorf("W_Stdout 时 writer 应为 os.Stdout, got %T", h.handlerWriter.writer)
		}
		if h.handlerWriter.file != nil {
			t.Error("W_Stdout 时不应打开日志文件")
		}
	})

	t.Run("file only", func(t *testing.T) {
		chdirTemp(t)
		h, err := newHandler(W_File, slog.LevelDebug)
		if err != nil {
			t.Fatal(err)
		}
		defer h.close()

		if h.handlerWriter.file == nil {
			t.Fatal("W_File 时应打开日志文件")
		}
		if h.handlerWriter.writer != h.handlerWriter.file.f {
			t.Errorf("W_File 时 writer 应为共享文件句柄, got %T", h.handlerWriter.writer)
		}
		if err := h.handlerWriter.Write([]byte("abc\n")); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("log.log")
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "abc\n" {
			t.Errorf("文件内容应为 %q, 实际 %q", "abc\n", string(data))
		}
	})

	t.Run("stdout and file", func(t *testing.T) {
		chdirTemp(t)
		h, err := newHandler(W_Stdout|W_File, slog.LevelDebug)
		if err != nil {
			t.Fatal(err)
		}
		defer h.close()

		if h.handlerWriter.file == nil {
			t.Fatal("多路输出时应打开日志文件")
		}
		// io.MultiWriter 返回未导出类型，用 %T 校验组合写入
		if typ := fmt.Sprintf("%T", h.handlerWriter.writer); !strings.Contains(typ, "multiWriter") {
			t.Errorf("W_Stdout|W_File 时 writer 应为 MultiWriter, got %s", typ)
		}
	})

	t.Run("shared file handle", func(t *testing.T) {
		// 指向同一文件的多个 writer 复用同一个 *os.File（引用计数），
		// 最后一个 close 才真正关闭文件并从池中移除。
		chdirTemp(t)

		w1, err := newHandlerWriter(W_File)
		if err != nil {
			t.Fatal(err)
		}
		w2, err := newHandlerWriter(W_File)
		if err != nil {
			t.Fatal(err)
		}

		if w1.file == nil || w2.file == nil {
			t.Fatal("W_File 应持有共享文件引用")
		}
		if w1.file != w2.file {
			t.Fatal("指向同一文件的 writer 应共享同一个 fileRef")
		}
		ref := w1.file
		if ref.refs != 2 {
			t.Errorf("打开两个 writer 后引用计数应为 2, 实际 %d", ref.refs)
		}

		// 两个 writer 各自写入，内容都落到同一文件
		if err := w1.Write([]byte("a\n")); err != nil {
			t.Fatal(err)
		}
		if err := w2.Write([]byte("b\n")); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("log.log")
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "a\nb\n" {
			t.Errorf("共享文件内容应为 %q, 实际 %q", "a\nb\n", string(data))
		}

		// 释放一个引用：文件仍打开
		w1.close()
		if ref.refs != 1 {
			t.Errorf("释放一个引用后计数应为 1, 实际 %d", ref.refs)
		}
		if err := w2.Write([]byte("c\n")); err != nil {
			t.Fatal(err)
		}

		// 最后一个引用释放：文件真正关闭并从池中移除
		w2.close()
		if ref.refs != 0 {
			t.Errorf("全部释放后计数应为 0, 实际 %d", ref.refs)
		}
		filePoolMu.Lock()
		_, inPool := filePool[ref.path]
		filePoolMu.Unlock()
		if inPool {
			t.Error("引用归零后文件应从池中移除")
		}
	})
}

// TestNewLogger 验证公开 API：多个 logger 可独立指定输出与级别，互不影响。
// 合并自原 TestCustomLogger：经公共 API 验证组前缀与"不同输出"能力。
func TestNewLogger(t *testing.T) {
	chdirTemp(t)

	l1, err := NewLogger(W_File, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseLogger(l1)
	l1 = l1.WithGroup("g1")

	l2, err := NewLogger(W_Stdout, slog.LevelWarn)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseLogger(l2)

	l1.Info("file logger only", slog.Int("key", 10))
	l2.Debug("filtered debug") // l2 级别为 Warn，Debug 应被丢弃
	l2.Info("stdout only")     // l2 输出到 stdout，不进文件

	data, err := os.ReadFile("log.log")
	if err != nil {
		t.Fatal(err)
	}
	wantContains(t, string(data), "] file logger only | g1.key=10")
	wantNotContains(t, string(data), "filtered debug")
	wantNotContains(t, string(data), "stdout only")
}

func TestConcurrentWrites(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h)

	const goroutines, perGoroutine = 20, 50
	var wg sync.WaitGroup
	for g := range goroutines {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := range perGoroutine {
				l.Info("concurrent", slog.Int("g", g), slog.Int("i", i))
			}
		}(g)
	}
	wg.Wait()

	// mutex 应保证每行完整、不丢行（配合 -race 验证无数据竞争）
	out := buf.String()
	if got := strings.Count(out, "\n"); got != goroutines*perGoroutine {
		t.Errorf("并发写入行数应为 %d, 实际 %d", goroutines*perGoroutine, got)
	}
	for line := range strings.SplitSeq(strings.TrimSuffix(out, "\n"), "\n") {
		if !linePattern.MatchString(line + "\n") {
			t.Errorf("存在格式损坏的行: %q", line)
		}
		// goroutine 内调用，位置也应指向调用点（log_test.go）
		if !strings.Contains(line, "log_test.go:") {
			t.Errorf("goroutine 调用的位置应指向 log_test.go，实际: %q", line)
		}
	}
}
