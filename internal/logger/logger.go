package logger

import (
	"context"
	"fmt"
	"os"

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
	zapConfig.Level = getLogLevel()

	zapConfig.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-20s", caller.TrimmedPath()))
	}

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
}

func getLogLevel() zap.AtomicLevel {
	swellLog, ok := os.LookupEnv("SWELL_LOG")
	if !ok {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	switch swellLog {
	case "DEBUG":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "WARN":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}

}
