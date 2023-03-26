package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, fmt string, args ...any)
	Info(ctx context.Context, fmt string, args ...any)
	Warn(ctx context.Context, fmt string, args ...any)
	Error(ctx context.Context, fmt string, args ...any)
}
