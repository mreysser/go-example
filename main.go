package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mreysser/go-example/handler"
	"github.com/mreysser/go-example/logger"
	"github.com/mreysser/go-lifecycle"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// TODO: this is already a bit of a duplicate of the existing prometheus metrics, but works as a
// starting example.
var totalHttpRequests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "total_http_requests",
	Help: "Number of HTTP requests received",
})

var logr *log.Logger

func init() {
	logr = logger.GetLoggerFromContextOrDefault(context.Background())
	// TODO: import from env
	logr.SetLevel(log.DebugLevel)

	// TODO: more metrics
	prometheus.Register(totalHttpRequests)
}

func main() {
	token := lifecycle.GetDefaultLifecycleToken()

	e := echo.New()

	// This middleware handles per-request metrics
	e.Pre(metricsMiddleware)

	// This middleware loads a log entry with UUID into the echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_ = logger.GetEntryFromEchoContext(c)
			return next(c)
		}
	})

	// This middleware logs when a request is completed by a handler.
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: false,
		LogLatency:  true,
		Skipper:     skipInfraEndpoints,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			entry := logger.GetEntryFromEchoContext(c)
			if v.Error == nil {
				entry.WithFields(log.Fields{
					"URI":       v.URI,
					"method":    v.Method,
					"latencyMs": v.Latency.Milliseconds(),
					"status":    v.Status,
				}).Info("request complete")
			} else {
				entry.WithFields(log.Fields{
					"URI":       v.URI,
					"method":    v.Method,
					"latencyMs": v.Latency.Milliseconds(),
					"status":    v.Status,
					"error":     v.Error,
				}).Error("request error")
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.GET("/", handler.Hello)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/live", live)
	e.GET("/ready", ready)

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
		if !skipInfraEndpoints(c) {
			totalHttpRequests.Inc()
		}
		return next(c)
	}
}

func skipInfraEndpoints(c echo.Context) bool {
	switch c.Request().URL.Path {
	case "/metrics":
		fallthrough
	case "/live":
		fallthrough
	case "/ready":
		return true
	}
	return false
}

func live(c echo.Context) error {
	return c.String(http.StatusOK, "Alive")
}

func ready(c echo.Context) error {
	// TODO: actual readiness should depend on the other resources this service relies on
	return c.String(http.StatusOK, "Ready")
}
