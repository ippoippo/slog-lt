package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

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

// AddTraceIdToCtx is our own internal tracing ID to be inserted into the context
func AddTraceIdToCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceId, err := uuidV7String()
			if err != nil {
				slog.ErrorContext(
					c.Request().Context(),
					"error in generating trace id",
					slog.String("error", err.Error()),
				)
			}
			ctx := context.WithValue(c.Request().Context(), contextKey(traceIdKeyString), traceId)
			request := c.Request().Clone(ctx)
			c.SetRequest(request)
			return next(c)
		}
	}
}

// RequestLogging logs appropriate info about the request
func RequestLogging(slogger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			slogger.InfoContext(
				req.Context(),
				"REQUEST",
				slog.String("uri", req.RequestURI),
				slog.String("method", req.Method),
				slog.String(echo.HeaderXRequestID, req.Header.Get(echo.HeaderXRequestID)),
			)
			return next(c)
		}
	}
}

func stringFromCtx(ctx context.Context, key string) string {
	value := ""
	ctxValue := ctx.Value(contextKey(key))
	if ctxValue != nil {
		value = ctxValue.(string)
	}
	return value
}

func uuidV7String() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
