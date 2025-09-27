package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type CtxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}

	if ctxLog, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok && ctxLog == logger {
		return ctx
	}

	return context.WithValue(ctx, CtxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

// MIDDLEWARES
func AddLoggerMiddleware(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextWithLogger(r.Context(), logger)

		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	}
}

func LogRequestMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := FromContext(r.Context())

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		logger.Info("incoming request",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger.Info("request completed",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("duration", duration.String()),
			slog.Int("status_code", http.StatusOK),
		)
	}
}
