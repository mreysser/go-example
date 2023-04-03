package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mreysser/go-example/handler"
	"github.com/mreysser/go-example/lifecycle"
)

func init() {
	// TODO: import from env
	log.SetLevel(log.DEBUG)
}

func main() {
	token := lifecycle.InitializeLifecycle(context.Background(), []syscall.Signal{syscall.SIGTERM})

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", handler.Hello)

	go runServer(e)
	token.RegisterShutdownHandler(func() { e.Shutdown(token.Ctx) })

	<-token.Ctx.Done()
	log.Warn("application exit")
}

func runServer(e *echo.Echo) {
	s := http.Server{
		Addr: ":8080",
	}

	log.Debug("starting server")

	err := e.StartServer(&s)
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("server failed to start: %w", err)
	}
}
