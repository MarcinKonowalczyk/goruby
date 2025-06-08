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
	if ctx == nil {
		// no context, no logger
		return nil
	}
	logger, ok := ctx.Value(loggerKey{}).(*log.Logger)
	if !ok {
		return nil
	}
	return logger
}

func Logf(ctx context.Context, format string, args ...any) {
	logger := GetLogger(ctx)
	if logger == nil {
		// no logger
		return
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
