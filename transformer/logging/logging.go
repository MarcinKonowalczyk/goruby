package logging

import (
	"context"
	"log"
)

// attach a logger to a context

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func GetLogger(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*log.Logger)
	if !ok {
		return nil // or return a default logger
	}
	return logger
}

func Logf(ctx context.Context, format string, args ...interface{}) {
	logger := GetLogger(ctx)
	if logger == nil {
		return // or handle the case where no logger is set
	}
	logger.Printf(format, args...)
}

func Log(ctx context.Context, message string) {
	logger := GetLogger(ctx)
	if logger == nil {
		return // or handle the case where no logger is set
	}
	logger.Print(message)
}
