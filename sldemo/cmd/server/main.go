package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ippoippo/slog-lt/sldemo/account"
	mware "github.com/ippoippo/slog-lt/sldemo/middleware"
)

func main() {
	slogger := setupSlogger()

	// Setup Echo
	e := setupEchoMiddlewareAndRoutes(slogger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			slogErrorWithOSExit("failed to e.Start", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slogErrorWithOSExit("failed to e.Shutdown", err)
	}
	slog.Info("graceful shutdown complete")
}

func setupSlogger() *slog.Logger {
	var slogger *slog.Logger
	levelVar := strings.ToUpper(os.Getenv("SLOG_LEVEL"))
	if levelVar == "" {
		levelVar = slog.LevelInfo.String()
	}

	lvl := slog.LevelInfo
	err := lvl.UnmarshalText([]byte(levelVar))
	if err != nil {
		func() {
			if slogger != nil {
				slogger.Error("unable to setup log level, defaulting to INFO")
			}
		}()
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}
	baseHandler := slog.NewTextHandler(os.Stdout, opts)
	ctxHandler := mware.ContextHandler{Handler: baseHandler}
	slogger = slog.New(ctxHandler)
	return slogger
}

func setupEchoMiddlewareAndRoutes(slogger *slog.Logger) *echo.Echo {
	e := newEcho() // echo.New(), but with banner and port stdout supressed

	slog.SetDefault(slogger) // Overrides the default slog.* AND! log.* functions to use the handlers

	slog.Info("slog:configuring routes and middlware")

	e.Use(middleware.Recover())
	e.Use(mware.AddXRequestIdToCtx())
	e.Use(mware.AddTraceIdToCtx())
	e.Use(mware.RequestLogging(slogger))

	accountsGroup := e.Group("/accounts")
	aHandler := account.NewHandler()
	accountsGroup.GET("", aHandler.GetAllAccounts)
	accountsGroup.POST("", aHandler.CreateAccount)
	accountsGroup.GET("/:id", aHandler.GetAccount)
	accountsGroup.DELETE("/:id", aHandler.DeleteAccount)

	return e
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func slogErrorWithOSExit(msg string, err error) {
	slog.Error(msg, slog.String("err", err.Error()))
	os.Exit(1)
}
