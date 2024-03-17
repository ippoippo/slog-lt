package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/ippoippo/slog-lt/zapdemo/zapp"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

const traceIdKeyString = "trace-id"

type contextKey string

// AddXRequestIdToCtx will (if client supplies X-Request-Id) insert into context
func AddXRequestIdToCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Request().Header.Get(echo.HeaderXRequestID)
			if rid != "" {
				ctx := context.WithValue(c.Request().Context(), contextKey(echo.HeaderXRequestID), rid)
				request := c.Request().Clone(ctx)
				c.SetRequest(request)
			}
			return next(c)
		}
	}
}

// AddTraceIdWithZLoggerToCtx is our own internal tracing ID to be inserted
// into the context. We also insert a reference to a logger
func AddTraceIdWithZLoggerToCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add TraceIdToCtx
			traceId, err := uuidV7String()
			if err != nil {
				zap.L().Error("error in generating trace id", zap.Error(err))
			}
			ctx := context.WithValue(c.Request().Context(), contextKey(traceIdKeyString), traceId)

			// Create Trace Ctx Logger, and add to Ctx
			ctxLogger := zap.L().With(zap.String(string(traceIdKeyString), traceId))
			ctx = zapp.CtxWithLogger(ctx, ctxLogger)

			request := c.Request().Clone(ctx)
			c.SetRequest(request)
			return next(c)
		}
	}
}

// RequestLogging logs appropriate info about the request
func RequestLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			zlogger := zapp.LoggerFromCtx(ctx)
			zlogger.Info(
				"REQUEST",
				zap.String("uri", req.RequestURI),
				zap.String("method", req.Method),
				zap.String(echo.HeaderXRequestID, req.Header.Get(echo.HeaderXRequestID)),
			)
			return next(c)
		}
	}
}

func uuidV7String() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
