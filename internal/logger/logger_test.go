package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/nabind47/sensor_api/internal/logger"
)

func Test_ContextWithLogger(t *testing.T) {

	testCases := []struct {
		name   string
		ctx    context.Context
		logger *slog.Logger
		exists bool
	}{

		{
			name: "returns context without logger",
			ctx:  context.Background(),
		},
		{
			name:   "returns context if it is",
			ctx:    context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			exists: true,
		}, {
			name:   "inject logger",
			ctx:    context.Background(),
			logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := logger.ContextWithLogger(tc.ctx, tc.logger)

			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)
			if tc.exists != ok {
				t.Errorf("expected :%v got: %v", tc.exists, ok)
			}
		})
	}

}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "logger exits",
			ctx:      context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			expected: true,
		}, {
			name:     "new logger returned",
			ctx:      context.Background(),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.FromContext(tc.ctx)

			if tc.expected && logger == nil {
				t.Errorf("expected: %v, got: %v", tc.expected, logger)
			}
		})
	}
}
