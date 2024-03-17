package middleware

import (
	"context"
	"log/slog"
)

// ContextHandler is our base context handler, it will handle all requests
type ContextHandler struct {
	slog.Handler
}

// Enabled - pass through to underlying handler
func (ch ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ch.Handler.Enabled(ctx, level)
}

// Handle includes our ctx logging before defering to underlying handler
func (ch ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(ch.addContextInfoToLog(ctx)...)
	return ch.Handler.Handle(ctx, r)
}

// WithAttrs overriding default implementation
func (ch ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{ch.Handler.WithAttrs(attrs)}
}

// WithGroup overriding default implementation
func (ch ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{ch.Handler.WithGroup(name)}
}

func (ch ContextHandler) addContextInfoToLog(ctx context.Context) []slog.Attr {
	var as []slog.Attr

	group := slog.Group("ctx-information",
		slog.String(traceIdKeyString, stringFromCtx(ctx, traceIdKeyString)),
	)
	as = append(as, group)
	return as
}
