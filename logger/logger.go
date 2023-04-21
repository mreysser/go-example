package logger

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ContextLoggerKey struct{}

var CtxLoggerKey = &ContextLoggerKey{}

type RequestIdKey struct{}

var ReqIdKey = &RequestIdKey{}

// Adds a logger to a context. Returns the new context with the logger attached.
func AddLoggerToContext(log *logrus.Logger, ctx context.Context) context.Context {
	return context.WithValue(context.Background(), CtxLoggerKey, log)
}

// Retrieves the logger from the given context. If there is none, a default logger is returned, so
// that this function is guaranteed to never return nil.
func GetLoggerFromContextOrDefault(ctx context.Context) *logrus.Logger {
	log, found := ctx.Value(CtxLoggerKey).(*logrus.Logger)
	if !found || log == nil {
		log = logrus.New()
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	}
	return log
}

// Gets a logrus entry from the Echo handler context. If none has been added yet, a new entry is
// created along with a UUID so that the logging for this unique request can be tracked.
func GetEntryFromEchoContext(c echo.Context) *logrus.Entry {
	var entry *logrus.Entry
	value := c.Get("log_entry")
	if value == nil {
		entry = logrus.NewEntry(GetLoggerFromContextOrDefault(context.Background())).WithField("req_id", uuid.New().String())
		c.Set("log_entry", entry)
	} else {
		entry = value.(*logrus.Entry)
	}
	return entry
}

// Gets a logrus entry from the given context. If one
func GetEntryFromContextOrDefault(ctx context.Context) *logrus.Entry {
	entry, found := ctx.Value(ReqIdKey).(*logrus.Entry)
	if !found || entry == nil {
		entry = logrus.NewEntry(GetLoggerFromContextOrDefault(context.Background())).WithField("req_id", "0")
	}
	return entry
}

func AddEntryToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, ReqIdKey, entry)
}
