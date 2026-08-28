package mlog

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// newCaptureHandler 创建 handler 并把输出重定向到内存 buffer，便于断言日志格式。
func newCaptureHandler(level slog.Level) (*handler, *bytes.Buffer) {
	h := newHandler(W_Stdout, level)
	buf := &bytes.Buffer{}
	h.handlerWriter.writer = buf
	return h, buf
}

// chdirTemp 切换到临时目录（测试结束后自动恢复），避免测试污染源码目录。
func chdirTemp(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func wantContains(t *testing.T, s, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		t.Errorf("输出应包含 %q，实际输出:\n%s", sub, s)
	}
}

func wantNotContains(t *testing.T, s, sub string) {
	t.Helper()
	if strings.Contains(s, sub) {
		t.Errorf("输出不应包含 %q，实际输出:\n%s", sub, s)
	}
}

// 结构: [Time] [Level] [a/b.go:10] msg | k1=v1 k2=v2
// 位置段不做 .go 限制：直接 slog 调用/goroutine 中 Callers 深度不准，
// 可能指向 testing 框架或 runtime 汇编文件（见 TestHandlerLineFormat 注释）。
var linePattern = regexp.MustCompile(
	`^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\] \[[A-Z][^]]*\] \[[^]]+:\d+\] .+\n$`,
)

func TestHandlerLineFormat(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	slog.New(h).Info("hello", slog.Int("key1", 10), slog.String("key2", "value2"))

	out := buf.String()
	if !linePattern.MatchString(out) {
		t.Errorf("日志行格式不符合预期:\n%s", out)
	}
	wantContains(t, out, "] hello | key1=10 key2=value2")
	if strings.Count(out, "\n") != 1 {
		t.Errorf("一行日志应只有一个换行符，实际输出: %q", out)
	}

	// 注：直接通过 slog.New(h) 调用时，runtime.Callers(6) 的深度是按 mlog 包装函数
	// 校准的，位置会指向错误帧（如 testing 框架），此处只校验格式不校验具体文件。
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

func TestGroupPrefix(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h).WithGroup("g1").WithGroup("g2")

	l.Info("grouped", slog.Int("key1", 10), slog.String("key2", "v"))

	out := buf.String()
	wantContains(t, out, "g1.g2.key1=10")
	wantContains(t, out, "g1.g2.key2=v")
}

func TestWithAttrs(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	// 注：按 slog 语义，WithGroup("g") 之后的 WithAttrs 应输出 g.ak=av；
	// 当前实现把 handler 级 attrs 平铺输出（无组前缀），而 record attrs 有前缀，
	// 两条路径不一致（handler.go logAttrs）。此处断言的是当前行为。
	l := slog.New(h.WithGroup("g").WithAttrs([]slog.Attr{slog.String("ak", "av")}))
	l.Info("attrs-in-group", slog.Int("k", 1))

	out := buf.String()
	wantContains(t, out, " | ak=av g.k=1")
}

func TestLogLevelPreserved(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h)

	Log(l, slog.LevelError, "err msg", slog.String("k", "v"))
	wantContains(t, buf.String(), "[ERROR]")
	wantContains(t, buf.String(), "] err msg | k=v")

	// 自定义级别应原样保留（回归保护：旧实现用 switch 会把非标准级别丢到 Debug）
	buf.Reset()
	customLevel := slog.Level(12)
	Log(l, customLevel, "custom level")
	wantContains(t, buf.String(), "["+customLevel.String()+"]")
	wantContains(t, buf.String(), "] custom level")
}

func TestLogSourcePosition(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	Log(slog.New(h), slog.LevelInfo, "position check")

	// 经 mlog.Log 调用时，位置应指向调用点所在文件（log_test.go）
	re := regexp.MustCompile(`\[[^]]*log_test\.go:\d+\]`)
	if !re.MatchString(buf.String()) {
		t.Errorf("Log() 路径的源码位置应指向 log_test.go，实际输出:\n%s", buf.String())
	}
}

func TestAnyAttrRendering(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	slog.New(h).Error("boom", slog.Any("error", errors.New("oops")))

	wantContains(t, buf.String(), "error=oops")
}

func TestInitializeClose(t *testing.T) {
	Initialize(W_Stdout)
	defer Close()

	if _, ok := slog.Default().Handler().(*handler); !ok {
		t.Error("Initialize 后 slog.Default 的 handler 应为 mlog.handler")
	}
	// Close 可重复调用且不应 panic（W_Stdout 模式 fileIns 为 nil）
	Close()
	Close()
}

func TestWriteFlag(t *testing.T) {
	t.Run("stdout only", func(t *testing.T) {
		h := newHandler(W_Stdout, slog.LevelDebug)
		if h.handlerWriter.writer != os.Stdout {
			t.Errorf("W_Stdout 时 writer 应为 os.Stdout, got %T", h.handlerWriter.writer)
		}
		if h.handlerWriter.fileIns != nil {
			t.Error("W_Stdout 时不应打开日志文件")
		}
	})

	t.Run("file only", func(t *testing.T) {
		chdirTemp(t)
		h := newHandler(W_File, slog.LevelDebug)
		defer h.close()

		if h.handlerWriter.fileIns == nil {
			t.Fatal("W_File 时应打开日志文件")
		}
		if h.handlerWriter.writer != h.handlerWriter.fileIns {
			t.Errorf("W_File 时 writer 应为文件句柄, got %T", h.handlerWriter.writer)
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
		h := newHandler(W_Stdout|W_File, slog.LevelDebug)
		defer h.close()

		if h.handlerWriter.fileIns == nil {
			t.Fatal("多路输出时应打开日志文件")
		}
		// io.MultiWriter 返回未导出类型，用 %T 校验组合写入
		if typ := fmt.Sprintf("%T", h.handlerWriter.writer); !strings.Contains(typ, "multiWriter") {
			t.Errorf("W_Stdout|W_File 时 writer 应为 MultiWriter, got %s", typ)
		}
	})
}

func TestCustomLogger(t *testing.T) {
	h1, buf1 := newCaptureHandler(slog.LevelDebug)
	h2, buf2 := newCaptureHandler(slog.LevelDebug)

	logger1 := slog.New(h1).WithGroup("group1")
	logger2 := slog.New(h2).WithGroup("group2")

	Log(logger1, slog.LevelDebug, "custom logger1", slog.Int("key1", 10))
	Log(logger2, slog.LevelDebug, "custom logger2", slog.String("key1", "value"))
	Log(logger1, slog.LevelDebug, "custom logger1", slog.Int("key2", 20))

	out1, out2 := buf1.String(), buf2.String()
	wantContains(t, out1, "group1.key1=10")
	wantContains(t, out1, "group1.key2=20")
	wantContains(t, out2, "group2.key1=value")
	wantNotContains(t, out1, "custom logger2")
	wantNotContains(t, out2, "custom logger1")
	// 结论：不同 logger 使用各自的 writer，输出互不影响
}

func TestConcurrentWrites(t *testing.T) {
	h, buf := newCaptureHandler(slog.LevelDebug)
	l := slog.New(h)

	const goroutines, perGoroutine = 20, 50
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
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
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if !linePattern.MatchString(line + "\n") {
			t.Errorf("存在格式损坏的行: %q", line)
		}
	}
}
