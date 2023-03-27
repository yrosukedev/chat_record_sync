package http_app

import (
	"context"
	"fmt"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/logger"
	"net/http"
)

func (f *HTTPApp) createWecomEventHTTPHandler(ctx context.Context) http.Handler {
	innerHandler := &dummyWecomMessageHandler{
		ctx:    ctx,
		logger: f.logger,
	}
	handler, err := workwx.NewHTTPHandler(f.wecomConfig.EventPushToken, f.wecomConfig.EventPushEncodingAESKey, innerHandler)

	if err != nil {
		f.logger.Error(ctx, "[http app] fails to create workwx.HTTPHandler, err: %v", err)
		panic(fmt.Sprintf("fails to create workwx.HTTPHandler, err: %v", err))
	}

	return handler
}

type dummyWecomMessageHandler struct {
	ctx    context.Context
	logger logger.Logger
}

// OnIncomingMessage 一条消息到来时的回调。
func (h *dummyWecomMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	h.logger.Info(h.ctx, "[wecom event handler] receive incoming message: %#v", msg)
	return nil
}
