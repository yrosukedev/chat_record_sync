package http_app

import (
	"context"
	"net/http"
)

func (f *HTTPApp) createMultiplexer(ctx context.Context) http.Handler {
	mux := http.NewServeMux()

	f.registerHandler(ctx, mux, "/chat/sync", f.createChatSyncHTTPHandler(ctx))

	f.registerHandler(ctx, mux, "/wecom/event", f.createWecomEventHTTPHandler(ctx))

	f.registerHandlerFunc(ctx, mux, "/ping", f.pingPongHandler)

	f.registerHandlerFunc(ctx, mux, "/WW_verify_ug8Z0l9pl4sSR9a6.txt", f.fileHandler)

	return mux
}

func (f *HTTPApp) registerHandler(ctx context.Context, mux *http.ServeMux, pattern string, handler http.Handler) {
	mux.Handle(pattern, handler)
	f.logger.Info(ctx, "[http app] handler registered for path: %v", pattern)
}

func (f *HTTPApp) registerHandlerFunc(ctx context.Context, mux *http.ServeMux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, handler)
	f.logger.Info(ctx, "[http app] handler registered for path: %v", pattern)
}
