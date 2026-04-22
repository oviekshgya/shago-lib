package txt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmittmann/tint"
)

const HeaderCorrelationID = "X-Correlation-ID"

type ctxKey string

const correlationIDKey ctxKey = "correlation_id"

func NewCorrelationID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "fallback-correlation-id"
	}
	return hex.EncodeToString(buf)
}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	if id == "" {
		id = NewCorrelationID()
	}
	return context.WithValue(ctx, correlationIDKey, id)
}

func CorrelationID(ctx context.Context) string {
	v, _ := ctx.Value(correlationIDKey).(string)
	if v == "" {
		return "unknown"
	}
	return v
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get(HeaderCorrelationID)
		if corrID == "" {
			corrID = NewCorrelationID()
		}
		w.Header().Set(HeaderCorrelationID, corrID)
		next.ServeHTTP(w, r.WithContext(WithCorrelationID(r.Context(), corrID)))
	})
}

func New(service string, level string) *slog.Logger {
	lvl := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}

	logFormat := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if logFormat == "" {
		logFormat = "text"
	}

	var handler slog.Handler
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      lvl,
			AddSource:  true,
			TimeFormat: "2006-01-02 15:04:05",
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if attr.Key == slog.SourceKey {
					if source, ok := attr.Value.Any().(*slog.Source); ok {
						attr.Value = slog.StringValue(formatSource(source))
					}
				}
				return attr
			},
		})
	}

	return slog.New(handler).With("service", service)
}

func colorizeLevel(level slog.Level) string {
	label := strings.ToUpper(level.String())

	switch {
	case level <= slog.LevelDebug:
		return "\x1b[36m" + label + "\x1b[0m"
	case level < slog.LevelWarn:
		return "\x1b[32m" + label + "\x1b[0m"
	case level < slog.LevelError:
		return "\x1b[33m" + label + "\x1b[0m"
	default:
		return "\x1b[31m" + label + "\x1b[0m"
	}
}

func formatSource(source *slog.Source) string {
	if source == nil {
		return ""
	}

	file := filepath.Base(source.File)
	function := source.Function
	if function == "" {
		function = "unknown"
	} else {
		function = shortFunctionName(function)
	}

	return fmt.Sprintf("%s:%d %s", file, source.Line, function)
}

func shortFunctionName(function string) string {
	lastSlash := strings.LastIndex(function, "/")
	if lastSlash >= 0 && lastSlash+1 < len(function) {
		function = function[lastSlash+1:]
	}

	lastDot := strings.LastIndex(function, ".")
	if lastDot >= 0 && lastDot+1 < len(function) {
		return function[lastDot+1:]
	}

	return function
}

func ExampleTlog() {
	logger := New("mainan", "debug")

	ctx := context.Background()
	ctx = WithCorrelationID(ctx, NewCorrelationID())

	logger.Debug("debug biasa", "step", 1, "feature", "logger")
	logger.Info("info biasa", "data apa", map[string]interface{}{"data": "data"})
	logger.Warn("warn biasa", "module", "payment", "retry", true)
	logger.Error("error biasa", "err", "koneksi database timeout")

	logger.DebugContext(ctx, "debug context", "correlation_id", CorrelationID(ctx), "payload", "cek debug")
	logger.InfoContext(ctx, "info context", "correlation_id", CorrelationID(ctx), "status", "running")
	logger.WarnContext(ctx, "warn context", "correlation_id", CorrelationID(ctx), "attempt", 3)
	logger.ErrorContext(ctx, "error context", "correlation_id", CorrelationID(ctx), "reason", "failed process")

	err := fmt.Errorf("koneksi database timeout")

	logger.ErrorContext(ctx,
		"error context",
		"correlation_id", CorrelationID(ctx),
		"err", err,
	)

	logger.InfoContext(ctx, "no pending logDone to process", "correlation_id", CorrelationID(ctx))
	logger.DebugContext(ctx, "logDone Data", "Data", "logDone")
}
