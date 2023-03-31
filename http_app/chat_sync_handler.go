package http_app

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/chat_sync/bitable_storage/chat_record"
	pagination2 "github.com/yrosukedev/chat_record_sync/chat_sync/bitable_storage/pagination"
	"github.com/yrosukedev/chat_record_sync/chat_sync/http_controller"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/buffer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	use_case2 "github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/chat_record_service"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/writer/retry"
	"github.com/yrosukedev/chat_record_sync/config"
	"net/http"
)

func (f *HTTPApp) createChatSyncHTTPHandler(ctx context.Context) http.Handler {
	return http_controller.NewChatSyncHTTPController(ctx, f.createChatSyncUseCase(ctx), f.logger)
}

func (f *HTTPApp) createChatSyncUseCase(ctx context.Context) use_case2.UseCase {
	useCase := use_case2.NewChatSyncUseCase(
		buffer.NewReader(
			pagination.NewBatchReaderAdapter(
				wecom_chat2.NewPaginatedReaderAdapter(
					chat_record_service.NewAdapter(ctx, f.wecomClient, "", "", config.WeComChatRecordSDKTimeout, f.logger),
					nil,
					transformer.NewWeComMessageTransformerFactory(map[string]wecom_chat2.ChatRecordTransformer{
						wecom_chat2.MessageTypeText: transformer.NewWeComTextMessageTransformer(ctx, f.logger),
					},
						transformer.NewWeComDefaultMessageTransformer(ctx, f.logger))),
				pagination2.NewStorageAdapter(ctx, f.larkClient, config.PaginationStorageBitableAppToken, config.PaginationStorageBitableTableId, f.logger),
				config.PaginatedReaderPageSize)),
		retry.NewWriterAdapter(
			chat_record.NewStorageAdapter(ctx, f.larkClient, config.ChatStorageBitableAppToken, config.ChatStorageBitableTableId, f.logger)),
	)

	return useCase
}
