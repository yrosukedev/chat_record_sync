package http_app

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/logger"
	"net/http"
)

type LogHandler struct {
	ctx     context.Context
	handler http.Handler
	logger  logger.Logger
}

func NewLogHandler(ctx context.Context, handler http.Handler, logger logger.Logger) *LogHandler {
	return &LogHandler{
		ctx:     ctx,
		handler: handler,
		logger:  logger,
	}
}

func (l *LogHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	l.logger.Info(l.ctx, "[http handler] receive request, url: %v, method: %v, headers: %+v", request.URL, request.Method, request.Header)
	l.handler.ServeHTTP(writer, request)
}
