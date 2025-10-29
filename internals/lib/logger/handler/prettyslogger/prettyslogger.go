package prettyslogger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

func (opts *PrettyHandlerOptions) NewPrettyLogger(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	// TODO: write doc about any method of Value struct
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h == nil || &h.opts.SlogOpts == nil || h.opts.SlogOpts.Level == nil {
		return level >= slog.LevelInfo
	}
	defer func() {
		if r := recover(); r != nil {
			h.l.Fatalf("Recovered from panic: %v", r)
		}
	}()

	return level >= h.opts.SlogOpts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
		attrs:   h.attrs,
		opts:    h.opts,
	}
}
