package zapp

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const loggerKeyString = "zap-ctx-logger"

func CtxWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey(loggerKeyString), logger)
}

func LoggerFromCtx(ctx context.Context) *zap.Logger {
	var value *zap.Logger
	ctxValue := ctx.Value(contextKey(loggerKeyString))
	if ctxValue == nil {
		return zap.L() // Revert to global logger
	}
	value, ok := ctxValue.(*zap.Logger)
	if !ok {
		return zap.L() // Revert to global logger
	}
	return value
}
