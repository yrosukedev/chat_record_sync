package http_app

import (
	"context"
	"fmt"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"github.com/yrosukedev/chat_record_sync/logger"
	"github.com/yrosukedev/chat_record_sync/utils"
	"net/http"
)

func (f *HTTPApp) createWecomEventHTTPHandler(ctx context.Context) http.Handler {
	innerHandler := createWecomEventHandler(ctx, f.logger, f.createChatSyncUseCase(ctx))
	handler, err := workwx.NewHTTPHandler(f.wecomConfig.EventPushToken, f.wecomConfig.EventPushEncodingAESKey, innerHandler)

	if err != nil {
		f.logger.Error(ctx, "[http app] fails to create workwx.HTTPHandler, err: %v", err)
		panic(fmt.Sprintf("fails to create workwx.HTTPHandler, err: %v", err))
	}

	return handler
}

type wecomEventHandler struct {
	ctx     context.Context
	logger  logger.Logger
	usecase use_case.UseCase
}

func createWecomEventHandler(ctx context.Context, logger logger.Logger, usecase use_case.UseCase) *wecomEventHandler {
	return &wecomEventHandler{
		ctx:     ctx,
		logger:  logger,
		usecase: usecase,
	}
}

// OnIncomingMessage 一条消息到来时的回调。
func (h *wecomEventHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	h.logger.Info(h.ctx, "[wecom event handler] receive incoming message: %#v", msg)

	h.logger.Info(h.ctx, "[wecom event handler] use case is about to start")

	errs := h.usecase.Run(h.ctx)

	if len(errs) != 0 {
		err := utils.CombineErrors(errs)
		h.logger.Error(h.ctx, "[wecom event handler] use case is finished with errors: %v", err)
		return err
	}

	h.logger.Info(h.ctx, "[wecom event handler] use case is finished successfully")

	return nil
}
