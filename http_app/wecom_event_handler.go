package http_app

import (
	"context"
	"fmt"
	"github.com/xen0n/go-workwx"
	"net/http"
)

func (f *HTTPApp) createWecomEventHTTPHandler(ctx context.Context) http.Handler {
	handler, err := workwx.NewHTTPHandler(f.wecomConfig.EventPushToken, f.wecomConfig.EventPushEncodingAESKey, dummyWecomMessageHandler{})

	if err != nil {
		f.logger.Error(ctx, "[http app] fails to create workwx.HTTPHandler, err: %v", err)
		panic(fmt.Sprintf("fails to create workwx.HTTPHandler, err: %v", err))
	}

	return handler
}

type dummyWecomMessageHandler struct{}

// OnIncomingMessage 一条消息到来时的回调。
func (dummyWecomMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	// You can do much more!
	fmt.Printf("incoming message: %s\n", msg)
	return nil
}
