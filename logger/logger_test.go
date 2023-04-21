package logger_test

import (
	"context"
	"testing"

	logr "github.com/mreysser/go-example/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	// Set up test case with context containing logger
	logger := logrus.New()
	ctx := logr.AddLoggerToContext(logger, context.Background())
	result1 := logr.GetLoggerFromContextOrDefault(ctx)
	assert.NotNil(t, result1)

	// Call function and check if it returns expected logger
	result2 := logr.GetLoggerFromContextOrDefault(ctx)
	if result2 != result1 {
		t.Errorf("Expected logger %v, but got %v", result1, result2)
	}

	// Test case where context does not contain logger
	ctx2 := context.Background()
	result3 := logr.GetLoggerFromContextOrDefault(ctx2)
	assert.NotNil(t, result3)
	assert.NotEqual(t, result1, result3)
	if result3.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected debug level, but got %v", result3.GetLevel())
	}
	_, ok := result3.Formatter.(*logrus.JSONFormatter)
	if !ok {
		t.Errorf("Expected JSON formatter, but got %T", result3.Formatter)
	}
}
