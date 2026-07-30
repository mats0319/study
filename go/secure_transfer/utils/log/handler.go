package mlog

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Handler struct {
	*HandlerWriter

	Level  slog.Level
	Attrs  []slog.Attr
	Groups []string
}

var _ slog.Handler = (*Handler)(nil)

var bufferPool = sync.Pool{New: func() any { return new(bytes.Buffer) }} // 减少GC抖动

func defaultHandler() (*Handler, error) {
	h := &Handler{
		HandlerWriter: &HandlerWriter{WriterFlag: w_File | w_Stdout},
		Level:         slog.LevelDebug,
		Attrs:         []slog.Attr{},
		Groups:        []string{},
	}

	err := h.HandlerWriter.New("log.log", 1<<20)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func newHandler(fileName string, maxSize int64) (*Handler, error) {
	h := &Handler{
		HandlerWriter: &HandlerWriter{WriterFlag: w_File},
		Level:         slog.LevelDebug,
		Attrs:         []slog.Attr{},
		Groups:        []string{},
	}

	maxSize = maxSize << 20 // unit: MB
	err := h.HandlerWriter.New(fileName, maxSize)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.Level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// structure: `[Time] [Level] [a/b.go:10] log message | k1=v1 g.k2=v2`
	buf.WriteByte('[')
	buf.WriteString(r.Time.Format("2006-01-02 15:04:05.000"))
	buf.WriteString("] [")
	buf.WriteString(r.Level.String())
	buf.WriteString("] [")
	codePosition(buf)
	buf.WriteString("] ")
	buf.WriteString(r.Message)
	h.logAttrs(buf, r)
	buf.WriteByte('\n')

	return h.Write(buf.Bytes())
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	newInstance := &Handler{
		HandlerWriter: h.HandlerWriter,
		Level:         h.Level,
		Attrs:         make([]slog.Attr, len(h.Attrs)+len(attrs)),
		Groups:        make([]string, len(h.Groups)),
	}

	copy(newInstance.Attrs, h.Attrs)
	copy(newInstance.Attrs[len(h.Attrs):], attrs)
	copy(newInstance.Groups, h.Groups)

	return newInstance
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newInstance := &Handler{
		HandlerWriter: h.HandlerWriter,
		Level:         h.Level,
		Attrs:         make([]slog.Attr, len(h.Attrs)),
		Groups:        make([]string, len(h.Groups)+1),
	}

	copy(newInstance.Attrs, h.Attrs)
	copy(newInstance.Groups, h.Groups)
	newInstance.Groups[len(h.Groups)] = name

	return newInstance
}

func codePosition(buf *bytes.Buffer) {
	pc := make([]uintptr, 1)
	runtime.Callers(6, pc)

	fs := runtime.CallersFrames(pc)
	f, _ := fs.Next()

	fileName := f.File
	lastIndex := strings.LastIndex(fileName, "/")
	if lastIndex >= 0 {
		index := strings.LastIndex(fileName[:lastIndex], "/")
		if index >= 0 {
			fileName = fileName[index+1:]
		}
	}

	buf.WriteString(fileName)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(f.Line))
}

func (h *Handler) logAttrs(buf *bytes.Buffer, r slog.Record) {
	if len(h.Attrs) == 0 && r.NumAttrs() == 0 {
		return
	}

	buf.WriteString(" |")

	for _, attr := range h.Attrs {
		buf.WriteByte(' ')
		buf.WriteString(attr.Key)
		buf.WriteByte('=')
		buf.WriteString(attr.Value.String())
	}

	r.Attrs(func(attr slog.Attr) bool {
		buf.WriteByte(' ')
		for _, v := range h.Groups {
			buf.WriteString(v)
			buf.WriteByte('.')
		}
		buf.WriteString(attr.Key)
		buf.WriteByte('=')
		buf.WriteString(attr.Value.String())

		return true
	})
}
