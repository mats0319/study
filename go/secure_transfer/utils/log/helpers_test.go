package mlog

import (
	"bytes"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"
)

// newCaptureHandler 创建独立的 handler 并把输出重定向到内存 buffer，便于断言日志格式。
func newCaptureHandler(level slog.Level) (*handler, *bytes.Buffer) {
	h, err := newHandler(W_Stdout, level)
	if err != nil {
		panic(err) // 测试环境 W_Stdout 必然合法
	}
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

// linePattern 匹配单行日志，结构为 [Time] [Level] [a/b.go:10] msg | k1=v1 k2=v2。
// 位置段不做 .go 限制（避免与实现细节耦合），具体文件由 TestLogSourcePosition 校验。
var linePattern = regexp.MustCompile(
	`^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\] \[[A-Z][^]]*\] \[[^]]+:\d+\] .+\n$`,
)
