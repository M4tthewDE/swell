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
	devEncoderConfig := zap.NewDevelopmentEncoderConfig()
	devEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	devEncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-25s", caller.TrimmedPath()))
	}

	logLevel := getLogLevel()

	consoleEncoder := zapcore.NewConsoleEncoder(devEncoderConfig)

	jsonEncoderConfig := zap.NewProductionEncoderConfig()
	jsonEncoderConfig.TimeKey = "timestamp"
	jsonEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(jsonEncoderConfig)

	stdoutCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	logFile, err := os.Create("swell.log.json")
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %w", err)
	}

	fileCore := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(logFile),
		logLevel,
	)

	core := zapcore.NewTee(stdoutCore, fileCore)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

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
