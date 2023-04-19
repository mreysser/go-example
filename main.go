package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mreysser/go-example/handler"
	"github.com/mreysser/go-lifecycle"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var totalHttpRequests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "total_http_requests",
	Help: "Number of HTTP requests received",
})

func init() {
	// TODO: import from env
	log.SetLevel(log.DebugLevel)

	prometheus.Register(totalHttpRequests)
}

func main() {
	token := lifecycle.GetDefaultLifecycleToken()

	e := echo.New()

	e.Pre(metricsMiddleware)
	e.Use(middleware.Logger())

	e.GET("/", handler.Hello)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go runServer(e)
	token.RegisterShutdownHandler(func(ctx context.Context) { e.Shutdown(ctx) })

	<-token.GetContext().Done()
	log.Warn("application exit")
}

func runServer(e *echo.Echo) {
	s := http.Server{
		Addr: ":8080",
	}

	err := e.StartServer(&s)
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("server failed to start: %w", err)
	}
}

func metricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		totalHttpRequests.Inc()
		return next(c)
	}
}
