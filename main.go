package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mreysser/go-example/handler"
)

func init() {
	log.SetLevel(log.DEBUG)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", handler.Hello)

	go runServer(e)

	<-ctx.Done()
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
