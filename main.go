package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	lifecycle "github.com/mreysser/go-example/lifecycle"
	log "github.com/sirupsen/logrus"
)

func main() {
	token := lifecycle.GetDefaultLifecycleToken()

	e := echo.New()
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello, world!") })

	go func() {
		if err := e.StartServer(&http.Server{Addr: ":8080"}); err != nil && err != http.ErrServerClosed {
			log.Errorf("server failed to start: %w", err)
			token.TerminateLifecycle()
		}
	}()

	token.RegisterShutdownHandler(func(ctx context.Context) { e.Shutdown(ctx) })
	<-token.GetContext().Done()
	log.Warn("application exit")
}
