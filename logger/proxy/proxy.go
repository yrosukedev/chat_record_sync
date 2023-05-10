package proxy

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelWarn  LogLevel = 3
	LogLevelError LogLevel = 4
)

type LoggerProxy struct {
	LogLevel LogLevel
	Logger   logger.Logger
}

func NewLoggerProxy(logLevel LogLevel, logger logger.Logger) *LoggerProxy {
	return &LoggerProxy{
		LogLevel: logLevel,
		Logger:   logger,
	}
}

func (p *LoggerProxy) Debug(ctx context.Context, fmt string, args ...any) {
	if p.Logger != nil && p.LogLevel <= LogLevelDebug {
		p.Logger.Debug(ctx, fmt, args...)
	}
}

func (p *LoggerProxy) Info(ctx context.Context, fmt string, args ...any) {
	if p.Logger != nil && p.LogLevel <= LogLevelInfo {
		p.Logger.Info(ctx, fmt, args...)
	}
}

func (p *LoggerProxy) Warn(ctx context.Context, fmt string, args ...any) {
	if p.Logger != nil && p.LogLevel <= LogLevelWarn {
		p.Logger.Warn(ctx, fmt, args...)
	}
}

func (p *LoggerProxy) Error(ctx context.Context, fmt string, args ...any) {
	if p.Logger != nil && p.LogLevel <= LogLevelError {
		p.Logger.Error(ctx, fmt, args...)
	}
}
