package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/ragarwalll/mta-forge.git/pkg/cli"
)

type Handler struct {
	opts   HandlerOptions
	out    io.Writer
	attrs  []slog.Attr
	groups []string
	mu     sync.Mutex
}

type HandlerOptions struct {
	Level       slog.Leveler
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	AddSource   bool
}

func NewHandler(out io.Writer, opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	return &Handler{out: out, opts: *opts}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)

	return h2
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)

	return h2
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error { //nolint:gocritic // false positive
	h.mu.Lock()
	defer h.mu.Unlock()

	levelColor := getLevelColor(r.Level)
	timeColor := color.New(color.FgHiBlack)
	messageColor := color.New(color.FgHiWhite)

	// Format time
	timeStr := timeColor.Sprintf("[%s]", r.Time.Format("15:04:05.000"))

	// Format level
	levelStr := levelColor.Sprintf("%-6s", r.Level.String())

	// Format message
	messageStr := messageColor.Sprint(r.Message)

	// Format attributes
	attrs := h.formatAttrs(r)

	// Print log entry
	fmt.Fprintf(h.out, "%s %s %s %s\n", timeStr, levelStr, messageStr, attrs)

	return nil
}

func (h *Handler) formatAttrs(r slog.Record) string { //nolint:gocritic // false positive
	attrs := make(map[string]interface{})

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	// Add handler attributes
	for _, a := range h.attrs {
		attrs[a.Key] = a.Value.Any()
	}

	// Apply ReplaceAttr if set
	if h.opts.ReplaceAttr != nil {
		for k, v := range attrs {
			a := h.opts.ReplaceAttr(h.groups, slog.Any(k, v))
			if a.Key != "" {
				attrs[a.Key] = a.Value.Any()
			} else {
				delete(attrs, k)
			}
		}
	}

	if len(attrs) == 0 {
		return ""
	}

	b, _ := json.MarshalIndent(attrs, "", "  ")

	return color.New(color.FgHiBlack).Sprint(string(b))
}

func (h *Handler) clone() *Handler {
	return &Handler{
		out:    h.out,
		opts:   h.opts,
		attrs:  append([]slog.Attr{}, h.attrs...),
		groups: append([]string{}, h.groups...),
	}
}

func getLevelColor(level slog.Level) *color.Color {
	switch level {
	case slog.LevelDebug:
		return color.New(color.FgHiBlack)
	case slog.LevelInfo:
		return color.New(color.FgHiCyan)
	case slog.LevelWarn:
		return color.New(color.FgHiYellow)
	case slog.LevelError:
		return color.New(color.FgHiRed)
	default:
		return color.New(color.FgHiWhite)
	}
}

func InitLogger() {
	var logger *slog.Logger

	opts := &HandlerOptions{
		AddSource: cli.GetForgerArgs().ExpandSource,
	}

	if cli.GetForgerArgs().Verbose {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}

	if cli.GetForgerArgs().Local {
		logger = slog.New(NewHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		}))
	}

	slog.SetDefault(logger)
}
