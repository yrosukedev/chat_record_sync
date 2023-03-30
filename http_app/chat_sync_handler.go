package http_app

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/chat_sync/chat_record_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/chat_sync/http_controller"
	"github.com/yrosukedev/chat_record_sync/chat_sync/paginated_reader"
	"github.com/yrosukedev/chat_record_sync/chat_sync/pagination_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/buffer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/retry_writer"
	use_case2 "github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom_chat"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom_chat_adapter"
	"github.com/yrosukedev/chat_record_sync/config"
	"net/http"
)

func (f *HTTPApp) createChatSyncHTTPHandler(ctx context.Context) http.Handler {
	return http_controller.NewChatSyncHTTPController(ctx, f.createChatSyncUseCase(ctx), f.logger)
}

func (f *HTTPApp) createChatSyncUseCase(ctx context.Context) use_case2.UseCase {
	useCase := use_case2.NewChatSyncUseCase(
		buffer.NewReader(
			paginated_reader.NewChatRecordPaginatedReader(
				wecom_chat2.NewPaginatedBufferedReaderAdapter(
					wecom_chat_adapter.NewWeComChatRecordServiceAdapter(ctx, f.wecomClient, "", "", config.WeComChatRecordSDKTimeout, f.logger),
					nil,
					wecom_chat2.NewWeComMessageTransformerFactory(map[string]wecom_chat2.ChatRecordTransformer{
						wecom_chat2.WeComMessageTypeText: wecom_chat2.NewWeComTextMessageTransformer(ctx, f.logger),
					},
						wecom_chat2.NewWeComDefaultMessageTransformer(ctx, f.logger))),
				pagination_bitable_storage.NewPaginationStorageAdapter(ctx, f.larkClient, config.PaginationStorageBitableAppToken, config.PaginationStorageBitableTableId, f.logger),
				config.PaginatedReaderPageSize)),
		retry_writer.NewRetryWriterAdapter(
			chat_record_bitable_storage.NewChatRecordStorageAdapter(ctx, f.larkClient, config.ChatStorageBitableAppToken, config.ChatStorageBitableTableId, f.logger)),
	)

	return useCase
}
