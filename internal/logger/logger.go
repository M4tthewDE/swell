package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKeyType string

var (
	contextKey = contextKeyType("logger")
)

func NewLogger() (*zap.SugaredLogger, error) {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func OnContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx != nil {
		if v := ctx.Value(contextKey); v != nil {
			return v.(*zap.SugaredLogger)
		}
	}

	panic("No logger found in context")
	return nil
}
