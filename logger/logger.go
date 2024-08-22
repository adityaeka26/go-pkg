package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger() *Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	writer := zapcore.AddSync(os.Stdout)
	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.InfoLevel))

	return &Logger{
		log: logger,
	}
}

func (l *Logger) GetLog() *zap.Logger {
	return l.log
}

func (l *Logger) Close(ctx context.Context) error {
	l.log.Sync()
	return nil
}
