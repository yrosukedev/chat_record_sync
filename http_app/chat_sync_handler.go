package http_app

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/buffer_reader"
	"github.com/yrosukedev/chat_record_sync/chat_record_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/config"
	"github.com/yrosukedev/chat_record_sync/http_controller"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"github.com/yrosukedev/chat_record_sync/pagination_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/retry_writer"
	"github.com/yrosukedev/chat_record_sync/use_case"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
	"github.com/yrosukedev/chat_record_sync/wecom_chat_adapter"
	"net/http"
)

func (f *HTTPApp) createChatSyncHTTPHandler(ctx context.Context) http.Handler {
	return http_controller.NewChatSyncHTTPController(ctx, f.createChatSyncUseCase(ctx), f.logger)
}

func (f *HTTPApp) createChatSyncUseCase(ctx context.Context) use_case.UseCase {
	useCase := use_case.NewSyncChatRecordUseCase(
		buffer_reader.NewChatRecordBufferedReaderAdapter(
			paginated_reader.NewChatRecordPaginatedReader(
				wecom_chat.NewPaginatedBufferedReaderAdapter(
					wecom_chat_adapter.NewWeComChatRecordServiceAdapter(f.wecomClient, "", "", config.WeComChatRecordSDKTimeout),
					nil,
					wecom_chat.NewWeComMessageTransformerFactory(map[string]wecom_chat.ChatRecordTransformer{
						wecom_chat.WeComMessageTypeText: wecom_chat.NewWeComTextMessageTransformer(),
					},
						wecom_chat.NewWeComDefaultMessageTransformer())),
				pagination_bitable_storage.NewPaginationStorageAdapter(ctx, f.larkClient, config.PaginationStorageBitableAppToken, config.PaginationStorageBitableTableId),
				config.PaginatedReaderPageSize)),
		retry_writer.NewRetryWriterAdapter(
			chat_record_bitable_storage.NewChatRecordStorageAdapter(ctx, f.larkClient, config.ChatStorageBitableAppToken, config.ChatStorageBitableTableId, f.logger)),
	)

	return useCase
}
