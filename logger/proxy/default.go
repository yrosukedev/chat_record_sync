package proxy

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/logger"
	"log"
	"os"
)

type defaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func NewDefaultLogger() logger.Logger {
	logger := defaultLogger{
		debugLogger: log.New(os.Stdout, "[Debug]", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "[Info]", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[Warn]", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stdout, "[Error]", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return logger
}

func (l defaultLogger) Debug(ctx context.Context, fmt string, args ...any) {
	l.debugLogger.Printf(fmt, args...)
}

func (l defaultLogger) Info(ctx context.Context, fmt string, args ...any) {
	l.infoLogger.Printf(fmt, args...)
}

func (l defaultLogger) Warn(ctx context.Context, fmt string, args ...any) {
	l.warnLogger.Printf(fmt, args...)
}

func (l defaultLogger) Error(ctx context.Context, fmt string, args ...any) {
	l.errorLogger.Printf(fmt, args...)
}
