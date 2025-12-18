package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

type Field struct {
	Key   string
	Value any
}

type ZapLogger struct {
	logger *zap.Logger
}

func New(level string) *ZapLogger {
	config := zap.NewProductionConfig()
	// Parse level, default to info
	l, err := zapcore.ParseLevel(level)
	if err == nil {
		config.Level = zap.NewAtomicLevelAt(l)
	}

	logger, _ := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	return &ZapLogger{logger: logger}
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) Error(msg string, err error, fields ...Field) {
	f := l.toZapFields(fields)
	f = append(f, zap.Error(err))
	l.logger.Error(msg, f...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{logger: l.logger.With(l.toZapFields(fields)...)}
}

// WithContext extracts Trace ID from context if available
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	// Example: Extract "trace_id" from context
	// val := ctx.Value("trace_id")
	// if val != nil {
	// 	return l.With(Field{Key: "trace_id", Value: val})
	// }
	return l
}

func (l *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, len(fields))
	for i, f := range fields {
		zf[i] = zap.Any(f.Key, f.Value)
	}
	return zf
}
