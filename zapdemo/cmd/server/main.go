package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/ippoippo/slog-lt/zapdemo/account"
	mware "github.com/ippoippo/slog-lt/zapdemo/middleware"
)

func main() {
	var zlogger = zap.Must(zap.NewDevelopment())
	defer zlogger.Sync()
	zap.ReplaceGlobals(zlogger)

	// Setup Echo
	e := setupEchoMiddlewareAndRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":1324"); err != nil && err != http.ErrServerClosed {
			zlogger.Fatal("failed to e.Start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		zlogger.Fatal("failed to e.Shutdown", zap.Error(err))
	}
	zlogger.Info("graceful shutdown complete")
}

func setupEchoMiddlewareAndRoutes() *echo.Echo {
	e := newEcho() // echo.New(), but with banner and port stdout supressed

	// Example of Global
	zap.L().Info("zlog:configuring routes and middlware")

	e.Use(middleware.Recover())
	e.Use(mware.AddXRequestIdToCtx())
	e.Use(mware.AddTraceIdWithZLoggerToCtx())
	e.Use(mware.RequestLogging())

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
