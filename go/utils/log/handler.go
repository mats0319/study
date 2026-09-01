package mlog

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
)

type handler struct {
	*handlerWriter

	level  slog.Level
	attrs  []slog.Attr
	groups []string
}

var _ slog.Handler = (*handler)(nil)

var bufferPool = sync.Pool{New: func() any { return new(bytes.Buffer) }} // 避免反复创建新的buffer

func newHandler(wf writerFlag, level slog.Level) (*handler, error) {
	hw, err := newHandlerWriter(wf)
	if err != nil {
		return nil, err
	}

	return &handler{
		handlerWriter: hw,
		level:         level,
		attrs:         []slog.Attr{},
		groups:        []string{},
	}, nil
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// structure: `[Time] [Level] [xxx/xx.go:10] log message | g.k1=v1 g.k2=v2`
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

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	newInstance := &handler{
		handlerWriter: h.handlerWriter,
		level:         h.level,
		attrs:         make([]slog.Attr, len(h.attrs)+len(attrs)),
		groups:        make([]string, len(h.groups)),
	}

	copy(newInstance.attrs, h.attrs)
	copy(newInstance.attrs[len(h.attrs):], attrs)
	copy(newInstance.groups, h.groups)

	return newInstance
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newInstance := &handler{
		handlerWriter: h.handlerWriter,
		level:         h.level,
		attrs:         make([]slog.Attr, len(h.attrs)),
		groups:        make([]string, len(h.groups)+1),
	}

	copy(newInstance.attrs, h.attrs)
	copy(newInstance.groups, h.groups)
	newInstance.groups[len(h.groups)] = name

	return newInstance
}

func (h *handler) logAttrs(buf *bytes.Buffer, r slog.Record) {
	if len(h.attrs) == 0 && r.NumAttrs() == 0 {
		return
	}

	buf.WriteString(" |")

	for _, attr := range h.attrs {
		writeAttr(buf, attr, h.groups)
	}

	r.Attrs(func(attr slog.Attr) bool {
		writeAttr(buf, attr, h.groups)

		return true
	})
}

func writeAttr(buf *bytes.Buffer, attr slog.Attr, groups []string) {
	buf.WriteByte(' ')
	for _, v := range groups {
		buf.WriteString(v)
		buf.WriteByte('.')
	}
	buf.WriteString(attr.Key)
	buf.WriteByte('=')
	buf.WriteString(attr.Value.String())
}
